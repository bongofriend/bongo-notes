package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/data"
	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
)

type apiMux struct {
	http.ServeMux
	config      config.Config
	authService services.AuthService
}

type apiHandler interface {
	Register(*apiMux)
}

//TODO Add additional services

func newApiMux(c config.Config, a services.AuthService) *apiMux {
	return &apiMux{
		config:      c,
		authService: a,
	}
}

type AuthenticatedHttpHandlerFunc func(user models.User, w http.ResponseWriter, r *http.Request)

func (a *apiMux) authenticatedHandlerFunc(pattern string, h AuthenticatedHttpHandlerFunc) {
	a.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		user, ok := a.authService.Authenticate(r)
		if !ok {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}
		h(user, w, r)
	})
}

func InitApi(appContext context.Context, doneCh chan struct{}, c config.Config) {
	context, cancel := context.WithCancel(appContext)

	userRepo := data.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	apiMux := newApiMux(c, authService)
	handlers := []apiHandler{
		newTestHandler(),
		NewSwaggerHandler(c),
	}

	for _, h := range handlers {
		h.Register(apiMux)
	}

	middlewares := createMiddlewareStack(loggingMiddleware())

	server := &http.Server{
		Handler: middlewares(apiMux),
		Addr:    fmt.Sprintf(":%d", c.Port),
	}

	log.Printf("Server ready on port: %d\n", c.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
			cancel()
		}
	}()

	<-context.Done()
	log.Println("Shutting down API")
	if err := server.Shutdown(appContext); err != nil {
		log.Fatal(err)
	}
	doneCh <- struct{}{}

}
