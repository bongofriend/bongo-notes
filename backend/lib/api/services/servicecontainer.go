package services

import (
	"github.com/bongofriend/bongo-notes/backend/lib/api/data"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
)

type servicesContainerImpl struct {
	authService      AuthService
	notebooksService NotebookService
	notesService     NotesService
}

type ServicesContainer interface {
	AuthService() AuthService
	NotebooksService() NotebookService
	NotesService() NotesService
	Shutdown(chan struct{})
}

func (s servicesContainerImpl) Shutdown(doneCh chan struct{}) {
	doneCh <- struct{}{}
}

func (s servicesContainerImpl) AuthService() AuthService {
	return s.authService
}

func (s servicesContainerImpl) NotebooksService() NotebookService {
	return s.notebooksService
}

func (s servicesContainerImpl) NotesService() NotesService {
	return s.notesService
}

// TODO Better way to serialize errors to responses
func NewServicesContainer(c config.Config, r data.RepositoryContainer) servicesContainerImpl {
	return servicesContainerImpl{
		authService:      NewAuthService(c, r.UserRepository()),
		notebooksService: NewNotebooksService(r.NotebooksRepository()),
		notesService:     NewNotesService(c, r),
	}
}
