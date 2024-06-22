package models

import "github.com/google/uuid"

type Note struct {
	Id    uuid.UUID `json:"id"`
	Title string    `json:"Title"`
}
