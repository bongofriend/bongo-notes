package data

import (
	"log"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type NotebooksRepository interface {
	FetchByUserId(userId uuid.UUID) ([]models.Notebook, error)
	CreateNotebook(userId uuid.UUID, title string, description string) error
	HasNotebook(userId uuid.UUID, notebookId uuid.UUID) (bool, error)
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
func (n notebooksRepositoryImpl) CreateNotebook(userId uuid.UUID, title string, description string) error {
	tx, err := n.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	if _, err := tx.Exec("INSERT INTO notebooks(creater_id, id, title, description) VALUES ($1, $2, $3, $4)", userId, uuid.New(), title, description); err != nil {
		return err
	}
	return nil
}

type notebooksEntity struct {
	Id          int32  `db:"rowid"`
	UUIDId      string `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

// FetchByUserId implements NotebooksRepository.
func (n notebooksRepositoryImpl) FetchByUserId(userId uuid.UUID) ([]models.Notebook, error) {
	var notebooksEntities []notebooksEntity
	if err := n.db.Select(&notebooksEntities, "SELECT rowid, id, title, description FROM notebooks WHERE creater_id = $1", userId.String()); err != nil {
		return nil, err
	}
	notebooks := make([]models.Notebook, 0, len(notebooksEntities))
	for _, e := range notebooksEntities {
		notebookModel, err := n.entityToModel(e)
		if err != nil {
			log.Println(err)
			continue
		}
		notebooks = append(notebooks, notebookModel)
	}
	return notebooks, nil
}

// HasNotebook implements NotebooksRepository.
func (n notebooksRepositoryImpl) HasNotebook(userId uuid.UUID, notebookId uuid.UUID) (bool, error) {
	var count int32
	if err := n.db.Get(&count, "SELECT COUNT(*) FROM notebooks WHERE id = $1 and creater_id = $2", notebookId.String(), userId.String()); err != nil {
		return false, err
	}
	return count == 1, nil
}

func (r notebooksRepositoryImpl) entityToModel(n notebooksEntity) (models.Notebook, error) {
	notebookId, err := uuid.Parse(n.UUIDId)
	if err != nil {
		return models.Notebook{}, err
	}
	return models.Notebook{
		Id:          notebookId,
		Description: n.Description,
		Title:       n.Title,
	}, nil
}
