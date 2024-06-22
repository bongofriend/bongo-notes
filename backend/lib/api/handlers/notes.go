package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
	httputils "github.com/bongofriend/bongo-notes/backend/lib/api/utils"
	"github.com/google/uuid"
)

type notesHandler struct {
	notesService services.NotesService
}

// Register implements ApiHandler.
func (n notesHandler) Register(m *ApiMux) {
	m.AuthenticatedHandlerFunc("POST /notes/{notebookId}", n.CreateNewNote)
	m.AuthenticatedHandlerFunc("GET /notes/{notebookId}", n.GetNotesForNotebook)
}

func NewNotesHandler(s services.ServicesContainer) ApiHandler {
	return notesHandler{
		notesService: s.NotesService(),
	}
}

type createNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func postCreateNewNoteParams(r *http.Request) (createNoteRequest, error) {
	var params createNoteRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Println(err)
		return createNoteRequest{}, err
	}
	return params, nil
}

// CreateNewNote godoc
//
//	@Summary	Create a new note for a notebook
//	@Tags		notes
//	@Router		/notes/{notebookId} [post]
//	@Param		noteParams	body	handlers.createNoteRequest	true	"Parameters for creating a new note"
//	@Param		notebookId	path	string						true	"Notebook Id for new note"
//	@Success	200
//	@Security	BearerAuth
func (n notesHandler) CreateNewNote(user models.User, w http.ResponseWriter, r *http.Request) {
	notebookIdPath := r.PathValue("notebookId")
	if notebookIdPath == "" {
		httputils.NotFoundError(w)
		return
	}
	notebookId, err := uuid.Parse(notebookIdPath)
	if err != nil {
		httputils.NotFoundError(w)
		return
	}
	params, err := postCreateNewNoteParams(r)
	if err != nil {
		log.Println(err)
		httputils.BadRequestError(w)
		return
	}
	if err := n.notesService.AddNoteToNotebook(user, notebookId, params.Title, params.Content); err != nil {
		log.Println(err)
		httputils.InternalServerError(w)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

type getNotesForNotebookResponse struct {
	Notes []models.Note `json:"notes"`
}

// GetNotes godoc
//
//	@Summary	Get notes for notebook
//	@Tags		notes
//	@Router		/notes/{notebookId} [get]
//	@Param		notebookId	path		string	true	"Notebook Id for new note"
//	@Success	200			{object}	handlers.getNotesForNotebookResponse
//	@Security	BearerAuth
func (n notesHandler) GetNotesForNotebook(user models.User, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("notebookId")
	if id == "" {
		httputils.BadRequestError(w)
		return
	}
	notebookId, err := uuid.Parse(id)
	if err != nil {
		httputils.BadRequestError(w)
		return
	}
	notes, err := n.notesService.FetchNotes(user, notebookId)
	if err != nil {
		log.Println(err)
		httputils.BadRequestError(w)
		return
	}
	rsp := getNotesForNotebookResponse{
		Notes: notes,
	}
	httputils.WriteJson(w, http.StatusOK, rsp)
}
