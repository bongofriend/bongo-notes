package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bongofriend/bongo-notes/backend/lib/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

type SwaggerHandler struct {
	config      config.Config
	projectPath string
}

func NewSwaggerHandler(c config.Config) SwaggerHandler {
	_, b, _, _ := runtime.Caller(0)
	projectPath := filepath.Dir(b)

	return SwaggerHandler{
		config:      c,
		projectPath: projectPath,
	}
}

func (s SwaggerHandler) serveSwaggerDefinition(w http.ResponseWriter, r *http.Request) {
	swaggerDefinitionsPath := filepath.Join(s.projectPath, "..", "..", "docs", "swagger.json")
	if _, err := os.Stat(swaggerDefinitionsPath); err != nil {
		if os.IsNotExist(err) {
			notFoundError(w)
			return
		}
		log.Print(err)
		internalServerError(w)
		return
	}
	data, err := os.ReadFile(swaggerDefinitionsPath)
	if err != nil {
		log.Println(err)
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err = w.Write(data); err != nil {
		log.Println(err)
		internalServerError(w)
		return
	}

}

func (s SwaggerHandler) Register(a *apiMux) {
	if !s.config.IncludeSwagger {
		log.Println("Skipping swagger UI")
		return
	}

	a.HandleFunc("/swagger/swagger.json", s.serveSwaggerDefinition)
	a.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/swagger.json", s.config.Port)),
	))
}
