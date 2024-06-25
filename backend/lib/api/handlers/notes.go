package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
	"github.com/google/uuid"
)

type notesHandler struct {
	notesService services.NotesService
}

// Register implements ApiHandler.
func (n notesHandler) Register(m *ApiMux) {
	m.AuthenticatedServiceResponseHandlerFunc("POST /notes/{notebookId}", n.CreateNewNote)
	m.AuthenticatedServiceResponseHandlerFunc("GET /notes/{notebookId}", n.GetNotesForNotebook)
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
//	@Failure	400
//	@Failure	401
//	@Failure	500
//	@Security	BearerAuth
func (n notesHandler) CreateNewNote(user models.User, r *http.Request) ServiceResponse {
	notebookIdPath := r.PathValue("notebookId")
	if notebookIdPath == "" {
		return BadRequest(nil)
	}
	notebookId, err := uuid.Parse(notebookIdPath)
	if err != nil {
		return BadRequest(err)
	}
	params, err := postCreateNewNoteParams(r)
	if err != nil {
		return BadRequest(err)
	}
	if err := n.notesService.AddNoteToNotebook(user, notebookId, params.Title, params.Content); err != nil {
		return InternalServerError(err)
	}
	return Accepted()
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
//	@Failure	400
//	@Failure	500
//	@Failure	401
//	@Security	BearerAuth
func (n notesHandler) GetNotesForNotebook(user models.User, r *http.Request) ServiceResponse {
	id := r.PathValue("notebookId")
	if id == "" {
		return BadRequest(nil)
	}
	notebookId, err := uuid.Parse(id)
	if err != nil {
		return BadRequest(err)
	}
	notes, err := n.notesService.FetchNotes(user, notebookId)
	if err != nil {
		return InternalServerError(err)
	}
	rsp := getNotesForNotebookResponse{
		Notes: notes,
	}
	return Success(http.StatusOK, rsp)
}

// TODO
func UpdateNote(user models.User, r *http.Request) ServiceResponse {
	panic("not implemented")
}
