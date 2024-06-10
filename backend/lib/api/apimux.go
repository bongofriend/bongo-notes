package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/config"
)

type apiMux struct {
	http.ServeMux
	config config.Config
}

type apiHandler interface {
	Register(*apiMux)
}

//TODO Add additional services

func newApiMux(c config.Config) *apiMux {
	return &apiMux{
		config: c,
	}
}

func InitApi(appContext context.Context, doneCh chan struct{}, c config.Config) {
	context, cancel := context.WithCancel(appContext)
	apiMux := newApiMux(c)
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
