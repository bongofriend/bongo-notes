package data

import (
	"log"

	"github.com/bongofriend/bongo-notes/backend/lib/config"
	"github.com/jmoiron/sqlx"
)

type repositoryContainerImpl struct {
	userRepository UserRepository
}

type RepositoryContainer interface {
	UserRepository() UserRepository
}

func (r repositoryContainerImpl) UserRepository() UserRepository {
	return r.userRepository
}

func NewRepositoryContainer(c config.Config) RepositoryContainer {
	db, err := sqlx.Connect(c.Db.Driver, c.Db.Path)
	if err != nil {
		log.Fatal(err)
	}
	return repositoryContainerImpl{
		userRepository: NewUserRepository(db),
	}
}
