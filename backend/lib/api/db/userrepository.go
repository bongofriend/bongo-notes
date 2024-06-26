package db

import (
	"time"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindUserById(id uuid.UUID) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

type userEntity struct {
	Id           int32     `db:"rowid"`
	UUIDId       string    `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password"`
	UpdatedAt    time.Time `db:"updated_at"`
	CreatedAt    time.Time `db:"created_at"`
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return userRepositoryImpl{
		db: db,
	}
}

func (u userRepositoryImpl) FindUserById(id uuid.UUID) (models.User, error) {
	var userEntity userEntity
	if err := u.db.Get(&userEntity, "SELECT rowid, id ,username, password FROM users WHERE id = $1", id); err != nil {
		return models.User{}, err
	}
	return entityToModel(userEntity)
}

func (u userRepositoryImpl) GetUserByUsername(username string) (models.User, error) {
	var userEntity userEntity
	if err := u.db.Get(&userEntity, "SELECT rowid, id, username, password FROM users WHERE username = $1 LIMIT 1", username); err != nil {
		return models.User{}, err
	}
	return entityToModel(userEntity)
}

func entityToModel(e userEntity) (models.User, error) {
	userId, err := uuid.Parse(e.UUIDId)
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		Id:           userId,
		Username:     e.Username,
		PasswordHash: e.PasswordHash,
	}, nil
}
