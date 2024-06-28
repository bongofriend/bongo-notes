package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bongofriend/bongo-notes/backend/lib/api/db"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
	"github.com/google/uuid"
)

const (
	errorExitCode int = 2
)

type DiffingService interface {
	QueueFile(pathToNewContent string, noteId uuid.UUID)
	Start(context context.Context)
	done() <-chan struct{}
}

type diffingJob struct {
	noteId         uuid.UUID
	newContentPath string
}

type diffingServiceImpl struct {
	jobCh       chan diffingJob
	config      config.Config
	diffingRepo db.DiffingRepository
	doneCh      chan struct{}
}

// Done implements DiffingService.
func (d diffingServiceImpl) done() <-chan struct{} {
	return d.doneCh
}

// Start implements DiffinggService.
func (d diffingServiceImpl) Start(context context.Context) {
	for {
		defer func() {
			d.doneCh <- struct{}{}
			close(d.jobCh)
			close(d.doneCh)
		}()
		select {
		case <-context.Done():
			return
		case job := <-d.jobCh:
			if err := d.processJob(job); err != nil {
				log.Println(err)
			}
		}
	}
}

func (d diffingServiceImpl) processJob(job diffingJob) error {
	diffId, err := d.generatedDiff(job.newContentPath, job.noteId)
	if err != nil {
		return err
	}
	if err := d.diffingRepo.AddDiff(job.noteId, diffId); err != nil {
		return err
	}
	if err := d.updateNoteContent(job.noteId, job.newContentPath); err != nil {
		return err
	}
	return nil
}

func (d diffingServiceImpl) generatedDiff(newContentPath string, noteId uuid.UUID) (uuid.UUID, error) {
	if _, err := os.Stat(newContentPath); err != nil {
		return uuid.Nil, fmt.Errorf("no new content found in %s for note %s", newContentPath, noteId.String())
	}
	notesPath := filepath.Join(d.config.NotesFolderPath, noteId.String())
	if _, err := os.Stat(notesPath); err != nil {
		return uuid.Nil, fmt.Errorf("no note in path %s for note with ID %s", notesPath, noteId.String())

	}
	currentNotePath := filepath.Join(notesPath, "recent")
	if _, err := os.Stat(currentNotePath); err != nil {
		return uuid.Nil, fmt.Errorf("current state for note %s could not be found in %s", noteId.String(), notesPath)
	}
	cmd := exec.Command("diff", newContentPath, currentNotePath)
	output, _ := cmd.Output()
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode == errorExitCode {
		diffError := errors.New(string(output))
		return uuid.Nil, fmt.Errorf("could not diff: %w", diffError)
	}

	diffPath := filepath.Join(notesPath, "diffs")
	if err := os.MkdirAll(diffPath, 0755); err != nil {
		return uuid.Nil, fmt.Errorf("could not diff: %w", err)
	}
	diffId := uuid.New()
	diffFilePath := filepath.Join(diffPath, fmt.Sprintf("%s.diff", diffId))
	diffFile, err := os.Create(diffFilePath)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not diff: %w", err)
	}
	defer diffFile.Close()
	buf := make([]byte, 1024*1024)
	if _, err := io.CopyBuffer(diffFile, bytes.NewReader(output), buf); err != nil {
		return uuid.Nil, err
	}
	return diffId, nil
}

func (d diffingServiceImpl) updateNoteContent(noteId uuid.UUID, newContentPath string) error {
	if _, err := os.Stat(newContentPath); os.IsNotExist(err) {
		return fmt.Errorf("no content at %s", newContentPath)
	}
	recentNotePath := filepath.Join(d.config.NotesFolderPath, noteId.String(), "recent")
	if _, err := os.Stat(recentNotePath); os.IsNotExist(err) {
		return fmt.Errorf("note not found at %s", recentNotePath)
	}
	if err := os.Remove(recentNotePath); err != nil {
		return err
	}
	if err := os.Rename(newContentPath, recentNotePath); err != nil {
		return err
	}
	return nil
}

// QueueFile implements DiffinggService.
func (d diffingServiceImpl) QueueFile(pathToNewContent string, noteId uuid.UUID) {
	job := diffingJob{
		newContentPath: pathToNewContent,
		noteId:         noteId,
	}
	go func() {
		d.jobCh <- job
	}()
}

func NewDiffingService(config config.Config, diffingRepo db.DiffingRepository) DiffingService {
	return diffingServiceImpl{
		config:      config,
		jobCh:       make(chan diffingJob, 10),
		diffingRepo: diffingRepo,
		doneCh:      make(chan struct{}),
	}
}
