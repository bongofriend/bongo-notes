package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

type SwaggerHandler struct {
	config config.Config
}

func NewSwaggerHandler(c config.Config) SwaggerHandler {
	return SwaggerHandler{
		config: c,
	}
}

func (s SwaggerHandler) Register(a *apiMux) {
	if !s.config.IncludeSwagger() {
		log.Println("Skipping swagger UI")
		return
	}

	a.Handle("/swagger/docs/swagger.json", http.StripPrefix("/swagger/docs/", http.FileServer(http.Dir("docs"))))
	a.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/docs/swagger.json", s.config.Port)),
	))
}
