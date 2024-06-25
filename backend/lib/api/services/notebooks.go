package services

import (
	"errors"
	"strings"

	"github.com/bongofriend/bongo-notes/backend/lib/api/db"
	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
)

type NotebookService interface {
	FetchNotebooks(user models.User) ([]models.Notebook, error)
	CreateNotebook(user models.User, title string, descroption string) error
}

type notebooksServiceImpl struct {
	notebooksRepo db.NotebooksRepository
}

func NewNotebooksService(notebooksRepo db.NotebooksRepository) NotebookService {
	return notebooksServiceImpl{
		notebooksRepo: notebooksRepo,
	}
}

func (n notebooksServiceImpl) CreateNotebook(user models.User, title string, description string) error {
	cleanTitle := strings.Trim(title, "")
	cleanDesc := strings.Trim(description, "")
	if len(cleanTitle) == 0 || len(cleanDesc) == 0 {
		return errors.New("notebook title or description was empty")
	}
	return n.notebooksRepo.CreateNotebook(user.Id, title, description)
}

func (n notebooksServiceImpl) FetchNotebooks(user models.User) ([]models.Notebook, error) {
	return n.notebooksRepo.FetchByUserId(user.Id)
}
