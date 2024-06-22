package models

import "github.com/google/uuid"

//TODO Specific user model to pass to endpoint handlers

type User struct {
	Id           uuid.UUID
	Username     string
	PasswordHash string
}
