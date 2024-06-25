package services

import (
	"context"

	"github.com/bongofriend/bongo-notes/backend/lib/api/db"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
)

type servicesContainerImpl struct {
	authService      AuthService
	notebooksService NotebookService
	notesService     NotesService
	diffingService   DiffingService
}

type ServicesContainer interface {
	AuthService() AuthService
	NotebooksService() NotebookService
	NotesService() NotesService
	DiffingService() DiffingService
	Shutdown(chan struct{})
	Init(appContext context.Context)
}

// Init implements ServicesContainer.
func (s servicesContainerImpl) Init(appContext context.Context) {
	go s.diffingService.Start(appContext)
}

func (s servicesContainerImpl) Shutdown(doneCh chan struct{}) {
	<-s.diffingService.done()
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

func (s servicesContainerImpl) DiffingService() DiffingService {
	return s.diffingService
}

func NewServicesContainer(c config.Config, r db.RepositoryContainer) ServicesContainer {
	return servicesContainerImpl{
		authService:      NewAuthService(c, r.UserRepository()),
		notebooksService: NewNotebooksService(r.NotebooksRepository()),
		notesService:     NewNotesService(c, r),
		diffingService:   NewDiffingService(c, r.DiffingRespository()),
	}
}
