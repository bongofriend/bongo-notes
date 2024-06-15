package models

//TODO Specific user model to pass to endpoint handlers

type User struct {
	Id           int
	Username     string
	PasswordHash string
}
