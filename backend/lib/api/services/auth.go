package services

import (
	"log"
	"net/http"
	"strings"

	"github.com/bongofriend/bongo-notes/backend/lib/api/data"
	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
)

type AuthService interface {
	Authenticate(*http.Request) (models.User, bool)
}

type authServiceImpl struct {
	userRepo data.UserRepository
}

func NewAuthService(u data.UserRepository) AuthService {
	return authServiceImpl{
		userRepo: u,
	}
}

func (a authServiceImpl) Authenticate(r *http.Request) (models.User, bool) {
	token, ok := extractAuthToken(r)
	if !ok {
		return models.User{}, false
	}
	userId, err := decodeToken(token)
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
	bearerToken := r.Header.Get("Authentication")
	if !strings.HasPrefix(bearerToken, "Bearer: ") {
		return "", false
	}
	bearerToken = strings.Replace(bearerToken, "Bearer: ", "", 1)
	return bearerToken, true
}

// TODO
func decodeToken(token string) (int32, error) {
	return -1, nil
}
