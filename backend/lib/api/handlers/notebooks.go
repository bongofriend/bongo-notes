package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
	httputils "github.com/bongofriend/bongo-notes/backend/lib/api/utils"
)

type notebooksHandler struct {
	notebooksService services.NotebookService
}

// Register implements ApiHandler.
func (n notebooksHandler) Register(mux *ApiMux) {
	mux.AuthenticatedHandlerFunc("GET /notebooks", n.GetNotebooks)
	mux.AuthenticatedHandlerFunc("POST /notebooks", n.CreateNewNotebook)
}

func NewNotebooksHandler(services services.ServicesContainer) ApiHandler {
	return notebooksHandler{
		notebooksService: services.NotebooksService(),
	}
}

type getNotebooksResponse struct {
	Notebooks []models.Notebook `json:"notebooks"`
}

// GetNotebooks godoc
//
//	@Summary	Get notebooks created by user
//	@Tags		notebooks
//	@Router		/notebooks [get]
//	@Success	200	{object}	handlers.getNotebooksResponse
//	@Security	BearerAuth
func (n notebooksHandler) GetNotebooks(user models.User, w http.ResponseWriter, r *http.Request) {
	notebooks, err := n.notebooksService.FetchNotebooks(user)
	if err != nil {
		log.Println(err)
		httputils.BadRequestError(w)
		return
	}
	rsp := getNotebooksResponse{
		Notebooks: notebooks,
	}
	httputils.WriteJson(w, http.StatusOK, rsp)
}

type createNewnNotebookRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// CreateNewNotebooks godoc
//
//	@Summary	Create new notebook
//	@Tags		notebooks
//	@Router		/notebooks [post]
//	@Param		noteboolDetails	body	handlers.createNewnNotebookRequest	true	"Parameters for creating a new new notebook"
//	@Success	200
//	@Failure	401
//	@Security	BearerAuth
func (n notebooksHandler) CreateNewNotebook(user models.User, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params createNewnNotebookRequest
	if err := decoder.Decode(&params); err != nil {
		log.Println(err)
		httputils.BadRequestError(w)
		return
	}
	if err := n.notebooksService.CreateNotebook(user, params.Title, params.Description); err != nil {
		log.Println(err)
		httputils.InternalServerError(w)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}
