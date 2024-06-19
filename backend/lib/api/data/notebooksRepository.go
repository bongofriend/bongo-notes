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
	tx, err := n.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	if _, err := tx.Exec("INSERT INTO notebooks(creater_id, title, description) VALUES ($1, $2, $3)", userId, title, description); err != nil {
		return err
	}
	return nil
}

type notebooksEntity struct {
	Id          int32  `db:"rowid"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

// FetchByUserId implements NotebooksRepository.
func (n notebooksRepositoryImpl) FetchByUserId(userId int32) ([]models.Notebook, error) {
	var notebooksEntities []notebooksEntity
	if err := n.db.Select(&notebooksEntities, "SELECT rowid, title, description FROM notebooks WHERE creater_id = $1", userId); err != nil {
		return nil, err
	}
	notebooks := make([]models.Notebook, 0, len(notebooksEntities))
	for _, n := range notebooksEntities {
		notebookModel := models.Notebook{
			Id:          n.Id,
			Description: n.Description,
			Title:       n.Title,
		}
		notebooks = append(notebooks, notebookModel)
	}
	return notebooks, nil
}
