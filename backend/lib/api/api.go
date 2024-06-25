package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/data"
	"github.com/bongofriend/bongo-notes/backend/lib/api/handlers"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
)

func InitApi(appContext context.Context, errCh chan struct{}, doneCh chan struct{}, c config.Config) {
	repoContainer := data.NewRepositoryContainer(c)
	servicesContainer := services.NewServicesContainer(c, repoContainer)
	serviceDoneCh := make(chan struct{})
	repoDoneCh := make(chan struct{})
	muxDoneCh := make(chan struct{})

	apiMux := handlers.NewApiMux(c, servicesContainer.AuthService())
	handlers := []handlers.ApiHandler{
		handlers.NewSwaggerHandler(c),
		handlers.NewAuthHandler(servicesContainer),
		handlers.NewNotebooksHandler(servicesContainer),
		handlers.NewNotesHandler(servicesContainer),
	}

	for _, h := range handlers {
		h.Register(apiMux)
	}

	middlewares := CreateMiddlewareStack(Logger)

	server := &http.Server{
		Handler: middlewares(apiMux),
		Addr:    fmt.Sprintf(":%d", c.Port),
	}

	log.Printf("Server ready on port: %d\n", c.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
			errCh <- struct{}{}
		}
	}()

	<-appContext.Done()
	log.Println("Shuttting down services")
	go servicesContainer.Shutdown(serviceDoneCh)
	<-serviceDoneCh

	log.Println("Shutting down database")
	go repoContainer.Shutdown(repoDoneCh)
	<-repoDoneCh

	log.Println("Shutting down API")
	go func() {
		if err := server.Shutdown(appContext); err != nil {
			log.Println(err)
			errCh <- struct{}{}
		}
		muxDoneCh <- struct{}{}
	}()
	<-muxDoneCh
	close(repoDoneCh)
	close(serviceDoneCh)
	close(muxDoneCh)
	doneCh <- struct{}{}
}
