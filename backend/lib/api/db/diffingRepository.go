package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DiffingRepository interface {
	AddDiff(noteId uuid.UUID, diffId uuid.UUID) error
}

type diffingRepositoryImpl struct {
	db *sqlx.DB
}

// AddDiff implements DiffingRepository.
func (d diffingRepositoryImpl) AddDiff(noteId uuid.UUID, diffId uuid.UUID) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	if _, err := tx.Exec(`INSERT INTO note_diffs(id, note_id) VALUES($1, $2)`, diffId, noteId); err != nil {
		return err
	}
	return nil
}

func NewDiffingRepository(db *sqlx.DB) DiffingRepository {
	return diffingRepositoryImpl{
		db: db,
	}
}
