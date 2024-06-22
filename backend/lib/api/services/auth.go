package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bongofriend/bongo-notes/backend/lib/api/data"
	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Authenticate(*http.Request) (models.User, bool)
	GenerateToken(username string, password string) (string, error)
}

type authServiceImpl struct {
	c        config.Config
	userRepo data.UserRepository
}

func NewAuthService(c config.Config, u data.UserRepository) AuthService {
	return authServiceImpl{
		c:        c,
		userRepo: u,
	}
}

func (a authServiceImpl) Authenticate(r *http.Request) (models.User, bool) {
	token, ok := extractAuthToken(r)
	if !ok {
		return models.User{}, false
	}
	userId, err := a.decodeToken(token)
	if err != nil {
		log.Println(err)
		return models.User{}, false
	}
	user, err := a.userRepo.FindUserById(userId)
	if err != nil {
		log.Println(err)
		return models.User{}, false
	}
	return user, true
}

func extractAuthToken(r *http.Request) (string, bool) {
	bearerToken := r.Header.Get("Authorization")
	if !strings.HasPrefix(bearerToken, "Bearer: ") {
		return "", false
	}
	bearerToken = strings.Replace(bearerToken, "Bearer: ", "", 1)
	return bearerToken, true
}

func (a authServiceImpl) decodeToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sigining method: %s", t.Header["alg"])
		}
		return []byte(a.c.JwtSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("JWT could not be parsed")
	}
	userId, ok := claims["userId"]
	if !ok {
		return uuid.Nil, errors.New("UserId not found in JWT claims")
	}
	id, ok := userId.(string)
	if !ok {
		return uuid.Nil, errors.New("could not parse userId in JWT claims")
	}
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New("could not parse userId in JWT claims")
	}
	return userUUID, nil
}

func (u authServiceImpl) GenerateToken(username string, password string) (string, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
	})
	tokenString, err := token.SignedString([]byte(u.c.JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, err
}
