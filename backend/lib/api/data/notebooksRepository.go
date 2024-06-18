package data

import (
	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/jmoiron/sqlx"
)

type NotebooksRepository interface {
	FetchByUserId(userId int32) ([]models.Notebook, error)
	CreateNotebook(userId int32, title string, description string) error
}

type notebooksRepositoryImpl struct {
	db *sqlx.DB
}

func NewNotebooksRepository(db *sqlx.DB) NotebooksRepository {
	return notebooksRepositoryImpl{
		db: db,
	}
}

// CreateNotebook implements NotebooksRepository.
func (n notebooksRepositoryImpl) CreateNotebook(userId int32, title string, description string) error {
	panic("unimplemented")
}

// FetchByUserId implements NotebooksRepository.
func (n notebooksRepositoryImpl) FetchByUserId(userId int32) ([]models.Notebook, error) {
	panic("unimplemented")
}
