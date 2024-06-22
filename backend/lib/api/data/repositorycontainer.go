package data

import (
	"log"

	"github.com/bongofriend/bongo-notes/backend/lib/config"
	"github.com/jmoiron/sqlx"
)

type repositoryContainerImpl struct {
	db                  *sqlx.DB
	userRepository      UserRepository
	notebooksRepository NotebooksRepository
	notesRepository     NotesRepository
}

// Shutdown implements RepositoryContainer.
func (r repositoryContainerImpl) Shutdown(doneCh chan struct{}) {
	r.db.Close()
	doneCh <- struct{}{}
}

type RepositoryContainer interface {
	UserRepository() UserRepository
	NotebooksRepository() NotebooksRepository
	NotesRepository() NotesRepository
	Shutdown(chan struct{})
}

func (r repositoryContainerImpl) UserRepository() UserRepository {
	return r.userRepository
}

func (r repositoryContainerImpl) NotebooksRepository() NotebooksRepository {
	return r.notebooksRepository
}

func (r repositoryContainerImpl) NotesRepository() NotesRepository {
	return r.notesRepository
}

func NewRepositoryContainer(c config.Config) RepositoryContainer {
	db, err := sqlx.Connect(c.Db.Driver, c.Db.Path)
	if err != nil {
		log.Fatal(err)
	}
	return repositoryContainerImpl{
		db:                  db,
		userRepository:      NewUserRepository(db),
		notebooksRepository: NewNotebooksRepository(db),
		notesRepository:     NewNotesRepository(db),
	}
}
