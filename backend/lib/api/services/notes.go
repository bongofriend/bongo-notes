package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/bongofriend/bongo-notes/backend/lib/api/data"
	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
)

type NotesService interface {
	AddNoteToNotebook(user models.User, notebookId int32, noteTitle string, header *multipart.FileHeader) error
}

type notesServiceImpl struct {
	config       config.Config
	notebookRepo data.NotebooksRepository
	notesRepo    data.NotesRepository
}

// AddNoteToNotebook implements NotesService.
func (n notesServiceImpl) AddNoteToNotebook(user models.User, notebookId int32, noteTitle string, header *multipart.FileHeader) error {
	content, err := header.Open()
	if err != nil {
		return err
	}
	defer content.Close()
	if isPlainText, err := isPlainTextFile(content); err != nil || !isPlainText {
		return errors.New("file could not be validated")
	}
	hasNotebook, err := n.notebookRepo.HasNotebook(int32(user.Id), notebookId)
	if err != nil {
		return err
	}
	if !hasNotebook {
		return fmt.Errorf("user %d has not ownership of notebook %d", user.Id, notebookId)
	}
	filePath, err := n.writeNoteToDisk(content)
	if err != nil {
		return err
	}
	return n.notesRepo.AddNote(notebookId, noteTitle, filePath)
}

func isPlainTextFile(content io.Reader) (bool, error) {
	panic("unimplemented")
}

func (n notesServiceImpl) writeNoteToDisk(fileContent io.Reader) (string, error) {
	panic("not implemented")
}

func NewNotesService(c config.Config, r data.RepositoryContainer) NotesService {
	return notesServiceImpl{
		config:       c,
		notebookRepo: r.NotebooksRepository(),
		notesRepo:    r.NotesRepository(),
	}
}
