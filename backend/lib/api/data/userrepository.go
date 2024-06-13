package data

import "github.com/bongofriend/bongo-notes/backend/lib/api/models"

type UserRepository interface {
	FindUserById(id int32) (models.User, error)
}

type userRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return userRepositoryImpl{}
}

// TODO
func (u userRepositoryImpl) FindUserById(id int32) (models.User, error) {
	return models.User{}, nil
}
