package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJson[T any](w http.ResponseWriter, data T) {
	buf, err := json.Marshal(data)

	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func notFoundError(w http.ResponseWriter) {
	http.Error(w, "Not found", http.StatusNotFound)
}

func internalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func notAuthenticatedError(w http.ResponseWriter) {
	http.Error(w, "Not authenticated", http.StatusUnauthorized)
}
