package httputils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJson[T any](w http.ResponseWriter, data T) {
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

func NotFoundError(w http.ResponseWriter) {
	http.Error(w, "Not found", http.StatusNotFound)
}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func NotAuthenticatedError(w http.ResponseWriter) {
	http.Error(w, "Not authenticated", http.StatusUnauthorized)
}
