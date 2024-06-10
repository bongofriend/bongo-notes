package api

import (
	"net/http"
)

type TestHandler struct{}

// Greet godoc
//
//	@Summary		Greet
//	@Description	Greet
//	@Tags			test
//	@Produce		json
//	@Success		200	{string}	string
//	@Failure		500	{string}	string
//	@Router			/greet [get]
func (t TestHandler) Greet(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["greet"] = "Hello there"

	writeJson(w, data)
}

func (t TestHandler) Register(a *apiMux) {
	a.HandleFunc("/greet", t.Greet)
}

func newTestHandler() TestHandler {
	return TestHandler{}
}
