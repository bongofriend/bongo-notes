package db

import (
	"log"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type NotesRepository interface {
	AddNote(notebookId uuid.UUID, noteId uuid.UUID, title string, path string) error
	GetNotesForNotebook(notebookId uuid.UUID) ([]models.Note, error)
}

type notesRepositoryImpl struct {
	db *sqlx.DB
}

type noteEntity struct {
	Id    int    `db:"rowid"`
	UUID  string `db:"id"`
	Title string `db:"title"`
}

// GetNotesForNotebook implements NotesRepository.
func (n notesRepositoryImpl) GetNotesForNotebook(notebookId uuid.UUID) ([]models.Note, error) {
	var noteEntities []noteEntity
	if err := n.db.Select(&noteEntities, "SELECT rowid, id, title FROM notes WHERE notebook_id = $1", notebookId); err != nil {
		return nil, err
	}
	notes := make([]models.Note, 0, len(noteEntities))
	for _, e := range noteEntities {
		noteModel, err := n.entityToModel(e)
		if err != nil {
			log.Println(err)
			continue
		}
		notes = append(notes, noteModel)
	}
	return notes, nil
}

// AddNote implements NotesRepository.
func (n notesRepositoryImpl) AddNote(notebookId uuid.UUID, noteId uuid.UUID, title string, path string) error {
	tx, err := n.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	if _, err := tx.Exec("INSERT INTO notes(id, notebook_id, title, path) VALUES ($1, $2, $3, $4)", noteId, notebookId, title, path); err != nil {
		return err
	}
	return err
}

func NewNotesRepository(db *sqlx.DB) NotesRepository {
	return notesRepositoryImpl{
		db: db,
	}
}

func (n notesRepositoryImpl) entityToModel(e noteEntity) (models.Note, error) {
	id, err := uuid.Parse(e.UUID)
	if err != nil {
		return models.Note{}, err
	}
	return models.Note{
		Id:    id,
		Title: e.Title,
	}, err
}
