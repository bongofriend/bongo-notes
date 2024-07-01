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
	m.AuthenticatedServiceResponseHandlerFunc("POST /notes/{notebookId}", n.CreateNewNote)
	m.AuthenticatedServiceResponseHandlerFunc("GET /notes/{notebookId}", n.GetNotesForNotebook)
	m.AuthenticatedServiceResponseHandlerFunc("PUT /notes/{notebookId}/{noteId}", n.UpdateNote)
	m.AuthenticatedHandlerFunc("GET /notes/{notebookId}/{noteid}", n.GetNote)
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

type updateNoteRequest struct {
	Content string `json:"content"`
}

// UpdateNote godoc
//
//	@Summary	Update note
//	@Tags		notes
//	@Router		/notes/{notebookId}/{noteId} [put]
//	@Param		notebookId		path	string						true	"Id of Notebook which Note is part of"
//	@Param		noteId			path	string						true	"Id of note to update"
//	@Param		notebookParams	body	handlers.updateNoteRequest	true	"Pramas to update Note conet"
//	@Success	200
//	@Failure	400
//	@Failure	500
//	@Failure	401
//	@Security	BearerAuth
func (n notesHandler) UpdateNote(user models.User, r *http.Request) ServiceResponse {
	notebookPathId := r.PathValue("notebookId")
	if len(notebookPathId) == 0 {
		return BadRequest(nil)
	}
	notebookId, err := uuid.Parse(notebookPathId)
	if err != nil {
		return BadRequest(nil)
	}
	notePathId := r.PathValue("noteId")
	if len(notePathId) == 0 {
		return BadRequest(nil)
	}
	noteId, err := uuid.Parse(notePathId)
	if err != nil {
		return BadRequest(err)
	}
	var reqBody updateNoteRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		return BadRequest(err)
	}
	if err := n.notesService.UpdateNote(user, notebookId, noteId, reqBody.Content); err != nil {
		return InternalServerError(err)
	}
	return Accepted()
}

// GetNote godoc
//
//	@Summary	Get note
//	@Tags		notes
//	@Router		/notes/{notebookId}/{noteId} [get]
//	@Param		notebookId	path	string	true	"Id of Notebook which Note is part of"
//	@Param		noteId		path	string	true	"Id of note to read"
//	@Param		diffId		query	string	false	"Id of diff"
//	@Produce plain
//	@Success	200
//	@Failure	400
//	@Failure	500
//	@Failure	401
//	@Security	BearerAuth
func (n notesHandler) GetNote(user models.User, w http.ResponseWriter, r *http.Request) {
	notebookPathId := r.PathValue("notebookId")
	if len(notebookPathId) == 0 {
		log.Println("No notebookPath found")
		httputils.BadRequestError(w)
		return
	}
	notebookId, err := uuid.Parse(notebookPathId)
	if err != nil {
		log.Println(err)
		httputils.BadRequestError(w)
		return
	}
	notePathId := r.PathValue("nodeId")
	if len(notePathId) == 0 {
		log.Println("No noteId found")
		httputils.BadRequestError(w)
		return
	}
	noteId, err := uuid.Parse(notePathId)
	if err != nil {
		log.Println(err)
		httputils.BadRequestError(w)
		return
	}
	diffQueryId := r.URL.Query().Get("diff")
	if len(diffQueryId) == 0 {
		content, err := n.notesService.GetNote(user, notebookId, noteId)
		if err != nil {
			log.Println(err)
			httputils.InternalServerError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(content); err != nil {
			log.Println(err)
			httputils.InternalServerError(w)
		}
	} else {
		diffId, err := uuid.Parse(diffQueryId)
		if err != nil {
			log.Println(err)
			httputils.InternalServerError(w)
			return
		}
		content, err := n.notesService.GetPatchedNote(user, notebookId, noteId, diffId)
		if err != nil {
			log.Println(err)
			httputils.InternalServerError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(content); err != nil {
			log.Println(err)
			httputils.InternalServerError(w)
		}
	}

}
