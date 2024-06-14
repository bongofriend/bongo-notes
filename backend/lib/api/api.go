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

func InitApi(appContext context.Context, doneCh chan struct{}, c config.Config) {
	context, cancel := context.WithCancel(appContext)
	repoContainer := data.NewRepositoryContainer(c)
	servicesContainer := services.NewServicesContainer(c, repoContainer)

	apiMux := handlers.NewApiMux(c, servicesContainer.AuthService())
	handlers := []handlers.ApiHandler{
		handlers.NewSwaggerHandler(c),
		handlers.NewAuthHandler(servicesContainer),
	}

	for _, h := range handlers {
		h.Register(apiMux)
	}

	middlewares := CreateMiddlewareStack(LoggingMiddleware())

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
