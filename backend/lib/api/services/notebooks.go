package services

import (
	"github.com/bongofriend/bongo-notes/backend/lib/api/data"
	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
)

type NotebookService interface {
	FetchNotebooks(user models.User) ([]models.Notebook, error)
	CreateNotebook(user models.User, title string, descroption string) error
}

type notebooksServiceImpl struct {
	notebooksRepo data.NotebooksRepository
}

func NewNotebooksService(notebooksRepo data.NotebooksRepository) NotebookService {
	return notebooksServiceImpl{
		notebooksRepo: notebooksRepo,
	}
}

func (n notebooksServiceImpl) CreateNotebook(user models.User, title string, descroption string) error {
	panic("unimplemented")
}

func (n notebooksServiceImpl) FetchNotebooks(user models.User) ([]models.Notebook, error) {
	panic("unimplemented")
}
