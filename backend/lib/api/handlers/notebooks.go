package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
)

type notebooksHandler struct {
	notebooksService services.NotebookService
}

// Register implements ApiHandler.
func (n notebooksHandler) Register(mux *ApiMux) {
	mux.AuthenticatedServiceResponseHandlerFunc("GET /notebooks", n.GetNotebooks)
	mux.AuthenticatedServiceResponseHandlerFunc("POST /notebooks", n.CreateNewNotebook)
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
//	@Failure	400
//	 @Failure 401
//	@Security	BearerAuth
func (n notebooksHandler) GetNotebooks(user models.User, r *http.Request) ServiceResponse {
	notebooks, err := n.notebooksService.FetchNotebooks(user)
	if err != nil {
		return BadRequest(err)
	}
	rsp := getNotebooksResponse{
		Notebooks: notebooks,
	}
	return Success(http.StatusOK, rsp)
}

type createNewnNotebookRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// CreateNewNotebooks godoc
//
//		@Summary	Create new notebook
//		@Tags		notebooks
//		@Router		/notebooks [post]
//		@Param		noteboolDetails	body	handlers.createNewnNotebookRequest	true	"Parameters for creating a new new notebook"
//		@Success	200
//		@Failure	400
//	 @Failure 401
//		@Security	BearerAuth
func (n notebooksHandler) CreateNewNotebook(user models.User, r *http.Request) ServiceResponse {
	decoder := json.NewDecoder(r.Body)
	var params createNewnNotebookRequest
	if err := decoder.Decode(&params); err != nil {
		return BadRequest(err)
	}
	if err := n.notebooksService.CreateNotebook(user, params.Title, params.Description); err != nil {
		return InternalServerError(err)
	} else {
		return Ok()
	}

}
