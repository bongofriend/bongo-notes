package handlers

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
	httputils "github.com/bongofriend/bongo-notes/backend/lib/api/utils"
)

type notesHandler struct {
	notesService services.NotesService
}

// Register implements ApiHandler.
func (n notesHandler) Register(m *ApiMux) {
	m.AuthenticatedHandlerFunc("POST /notes/{notebook_id}", n.CreateNewNote)
}

func NewNotesHandler(s services.ServicesContainer) ApiHandler {
	return notesHandler{
		notesService: s.NotesService(),
	}
}

type createNoteRequest struct {
	Title string                `form:"title"`
	File  *multipart.FileHeader `form:"note"`
}

func postCreateNewNoteParams(r *http.Request) (createNoteRequest, error) {
	if err := r.ParseForm(); err != nil {
		return createNoteRequest{}, err
	}
	title := r.FormValue("title")
	if title == "" {
		return createNoteRequest{}, errors.New("note title was empty")
	}
	_, header, err := r.FormFile("note")
	if err != nil {
		return createNoteRequest{}, err
	}
	return createNoteRequest{
		Title: title,
		File:  header,
	}, nil
}

func (n notesHandler) CreateNewNote(user models.User, w http.ResponseWriter, r *http.Request) {
	notebookIdPath := r.PathValue("notebook_id")
	if notebookIdPath == "" {
		httputils.NotFoundError(w)
		return
	}
	notebookId, err := strconv.Atoi(notebookIdPath)
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
	if err := n.notesService.AddNoteToNotebook(user, int32(notebookId), params.Title, params.File); err != nil {
		log.Println(err)
		httputils.InternalServerError(w)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
