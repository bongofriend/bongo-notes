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
	panic("unimplemented")
}

func NewDiffingRepository(db *sqlx.DB) DiffingRepository {
	return diffingRepositoryImpl{
		db: db,
	}
}
