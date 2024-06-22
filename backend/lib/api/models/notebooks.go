package models

import "github.com/google/uuid"

type Notebook struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"descripton"`
}
