package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bongofriend/bongo-notes/backend/lib/api/data"
	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
	"github.com/google/uuid"
)

type NotesService interface {
	AddNoteToNotebook(user models.User, notebookId uuid.UUID, noteTitle string, content string) error
	FetchNotes(user models.User, notebookId uuid.UUID) ([]models.Note, error)
}

type notesServiceImpl struct {
	config       config.Config
	notebookRepo data.NotebooksRepository
	notesRepo    data.NotesRepository
}

// FetchNotes implements NotesService.
func (n notesServiceImpl) FetchNotes(user models.User, notebookId uuid.UUID) ([]models.Note, error) {
	hasNotebook, err := n.notebookRepo.HasNotebook(user.Id, notebookId)
	if err != nil {
		return nil, err
	}
	if !hasNotebook {
		return nil, fmt.Errorf("user %s has not ownership of notebook %s", user.Id, notebookId)
	}
	return n.notesRepo.GetNotesForNotebook(notebookId)
}

// AddNoteToNotebook implements NotesService.
func (n notesServiceImpl) AddNoteToNotebook(user models.User, notebookId uuid.UUID, noteTitle string, content string) error {
	if isPlainText, err := isValidNote(content); err != nil || !isPlainText {
		return errors.New("file could not be validated")
	}
	hasNotebook, err := n.notebookRepo.HasNotebook(user.Id, notebookId)
	if err != nil {
		return err
	}
	if !hasNotebook {
		return fmt.Errorf("user %d has not ownership of notebook %d", user.Id, notebookId)
	}
	noteId := uuid.New()
	filePath, err := n.writeNoteToDisk(noteId, content)
	if err != nil {
		return err
	}
	return n.notesRepo.AddNote(notebookId, noteId, noteTitle, filePath)
}

func isValidNote(content string) (bool, error) {
	log.Println(content)
	return true, nil
}

func (n notesServiceImpl) writeNoteToDisk(noteId uuid.UUID, fileContent string) (string, error) {
	notePath := filepath.Join(n.config.NotesFolderPath, noteId.String())
	if err := os.MkdirAll(notePath, 0777); err != nil {
		return "", err
	}
	noteFilePath := filepath.Join(notePath, "recent")
	file, err := os.Create(noteFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	reader := strings.NewReader(fileContent)
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}
	return notePath, nil
}

func NewNotesService(c config.Config, r data.RepositoryContainer) NotesService {
	return notesServiceImpl{
		config:       c,
		notebookRepo: r.NotebooksRepository(),
		notesRepo:    r.NotesRepository(),
	}
}
