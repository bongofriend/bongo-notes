package models

type Notebook struct {
	Id          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"descripton"`
}
