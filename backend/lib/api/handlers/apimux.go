package handlers

import (
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
	httputils "github.com/bongofriend/bongo-notes/backend/lib/api/utils"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
)

type ApiMux struct {
	http.ServeMux
	config      config.Config
	authService services.AuthService
}

//TODO Add additional services

func NewApiMux(c config.Config, a services.AuthService) *ApiMux {
	return &ApiMux{
		config:      c,
		authService: a,
	}
}

type AuthenticatedHttpHandlerFunc func(user models.User, w http.ResponseWriter, r *http.Request)

func (a *ApiMux) authenticatedHandlerFunc(pattern string, h AuthenticatedHttpHandlerFunc) {
	a.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		user, ok := a.authService.Authenticate(r)
		if !ok {
			httputils.NotAuthenticatedError(w)
			return
		}
		h(user, w, r)
	})
}

type ApiHandler interface {
	Register(*ApiMux)
}
