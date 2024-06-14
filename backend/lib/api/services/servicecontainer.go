package services

import (
	"github.com/bongofriend/bongo-notes/backend/lib/api/data"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
)

type servicesContainerImpl struct {
	authService AuthService
}

type ServicesContainer interface {
	AuthService() AuthService
}

func (s servicesContainerImpl) AuthService() AuthService {
	return s.authService
}

func NewServicesContainer(c config.Config, r data.RepositoryContainer) servicesContainerImpl {
	return servicesContainerImpl{
		authService: NewAuthService(c, r.UserRepository()),
	}
}
