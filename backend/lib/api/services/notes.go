package services

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bongofriend/bongo-notes/backend/lib/api/db"
	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
	"github.com/google/uuid"
)

type NotesService interface {
	AddNoteToNotebook(user models.User, notebookId uuid.UUID, noteTitle string, content string) error
	FetchNotes(user models.User, notebookId uuid.UUID) ([]models.Note, error)
	UpdateNote(user models.User, notebookId uuid.UUID, noteId uuid.UUID, notebookIdnewContent string) error
}

type notesServiceImpl struct {
	config         config.Config
	notebookRepo   db.NotebooksRepository
	notesRepo      db.NotesRepository
	diffingService DiffingService
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
	filePath, err := n.writeNewNoteToDisk(noteId, content)
	if err != nil {
		return err
	}
	return n.notesRepo.AddNote(notebookId, noteId, noteTitle, filePath)
}

// TODO
func isValidNote(content string) (bool, error) {
	log.Println(content)
	return true, nil
}

func (n notesServiceImpl) writeNewNoteToDisk(noteId uuid.UUID, fileContent string) (string, error) {
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
	adjustLineBreak(&fileContent)
	reader := strings.NewReader(fileContent)
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}
	return notePath, nil
}

// UpdateNote implements NotesService.
func (n notesServiceImpl) UpdateNote(user models.User, notebookdId uuid.UUID, noteId uuid.UUID, newContent string) error {
	isPartOfNotebook, err := n.notesRepo.IsNotePartOfNotebook(user.Id, notebookdId, noteId)
	if err != nil {
		return err
	}
	if !isPartOfNotebook {
		return fmt.Errorf("wrong note access: User-Id %s Notebook-Id %s Note-Id %s", user.Id, notebookdId, noteId)
	}
	ok, err := isValidNote(newContent)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("invalid content for note %s", noteId)
	}
	newContentPath, err := n.writeTempNoteToDisk(newContent)
	if err != nil {
		return err
	}
	n.diffingService.QueueFile(newContentPath, noteId)
	return nil
}

func (n notesServiceImpl) writeTempNoteToDisk(content string) (string, error) {
	notesTempPath := filepath.Join(n.config.NotesFolderPath, "temp")
	if err := os.MkdirAll(notesTempPath, 0755); err != nil {
		return "", err
	}
	contentHash, err := getHashByContent(strings.NewReader(content))
	if err != nil {
		return "", err
	}
	newNoteContentPath := filepath.Join(notesTempPath, contentHash)
	file, err := os.Create(newNoteContentPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	adjustLineBreak(&content)
	buf := make([]byte, 1024*1024)
	if _, err := io.CopyBuffer(file, strings.NewReader(content), buf); err != nil {
		return "", err
	}
	return newNoteContentPath, nil
}

func getHashByContent(r io.Reader) (string, error) {
	buff := make([]byte, 1024*1024)
	hash := sha1.New()
	if _, err := io.CopyBuffer(hash, r, buff); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func adjustLineBreak(s *string) {
	if strings.HasSuffix(*s, "\n") {
		return
	}
	*s = *s + "\n"
}

func NewNotesService(c config.Config, diffService DiffingService, notesRepo db.NotesRepository, notebookRepo db.NotebooksRepository) NotesService {
	return notesServiceImpl{
		config:         c,
		notebookRepo:   notebookRepo,
		notesRepo:      notesRepo,
		diffingService: diffService,
	}
}
