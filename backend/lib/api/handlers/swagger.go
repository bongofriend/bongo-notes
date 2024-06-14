package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	httputils "github.com/bongofriend/bongo-notes/backend/lib/api/utils"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
	httpSwagger "github.com/swaggo/http-swagger"
)

type swaggerHandler struct {
	config      config.Config
	projectPath string
}

func NewSwaggerHandler(c config.Config) ApiHandler {
	_, b, _, _ := runtime.Caller(0)
	projectPath := filepath.Dir(b)

	return swaggerHandler{
		config:      c,
		projectPath: projectPath,
	}
}

func (s swaggerHandler) serveSwaggerDefinition(w http.ResponseWriter, r *http.Request) {
	swaggerDefinitionsPath := filepath.Join(s.projectPath, "..", "..", "..", "docs", "swagger.json")
	if _, err := os.Stat(swaggerDefinitionsPath); err != nil {
		if os.IsNotExist(err) {
			httputils.NotFoundError(w)
			return
		}
		log.Print(err)
		httputils.InternalServerError(w)
		return
	}
	data, err := os.ReadFile(swaggerDefinitionsPath)
	if err != nil {
		log.Println(err)
		httputils.InternalServerError(w)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if _, err = w.Write(data); err != nil {
		log.Println(err)
		httputils.InternalServerError(w)
		return
	}

}

func (s swaggerHandler) Register(a *ApiMux) {
	if !s.config.IncludeSwagger {
		log.Println("Skipping swagger UI")
		return
	}

	a.HandleFunc("/swagger/swagger.json", s.serveSwaggerDefinition)
	a.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/swagger.json", s.config.Port)),
	))
}
