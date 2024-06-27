package services

import (
	"bytes"
	"context"
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
	if err := cmd.Run(); err != nil {
		return uuid.Nil, err
	}
	diffPath := filepath.Join(notesPath, "diffs")
	if err := os.MkdirAll(diffPath, 0755); err != nil {
		return uuid.Nil, err
	}
	diffId := uuid.New()
	diffFilePath := filepath.Join(diffPath, fmt.Sprintf("%s.diff", diffId))
	diffFile, err := os.Create(diffFilePath)
	if err != nil {
		return uuid.Nil, err
	}
	defer diffFile.Close()
	cmdOutput, err := cmd.Output()
	if err != nil {
		return uuid.Nil, err
	}
	buf := make([]byte, 1024*1024)
	if _, err := io.CopyBuffer(diffFile, bytes.NewReader(cmdOutput), buf); err != nil {
		return uuid.Nil, err
	}
	return diffId, nil

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
