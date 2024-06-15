package data

import (
	"time"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindUserById(id int32) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

type userEntity struct {
	Id           int32     `db:"rowid"`
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

func (u userRepositoryImpl) FindUserById(id int32) (models.User, error) {
	var userEntity userEntity
	if err := u.db.Get(&userEntity, "SELECT rowid, username, password FROM users WHERE rowid = $1", id); err != nil {
		return models.User{}, err
	}
	return models.User{
		Id:           int(userEntity.Id),
		Username:     userEntity.Username,
		PasswordHash: userEntity.PasswordHash,
	}, nil
}

func (u userRepositoryImpl) GetUserByUsername(username string) (models.User, error) {
	var userEntity userEntity
	if err := u.db.Get(&userEntity, "SELECT rowid, username, password FROM users WHERE username = $1 LIMIT 1", username); err != nil {
		return models.User{}, err
	}
	return models.User{
		Id:           int(userEntity.Id),
		Username:     userEntity.Username,
		PasswordHash: userEntity.PasswordHash,
	}, nil
}
