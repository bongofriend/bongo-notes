package data

import "github.com/jmoiron/sqlx"

type NotesRepository interface {
	AddNote(notebookId int32, title string, path string) error
}

type notesRepositoryImpl struct {
	db *sqlx.DB
}

// AddNote implements NotesRepository.
func (n notesRepositoryImpl) AddNote(notebookId int32, title string, path string) error {
	panic("unimplemented")
}

func NewNotesRepository(db *sqlx.DB) NotesRepository {
	return notesRepositoryImpl{
		db: db,
	}
}
