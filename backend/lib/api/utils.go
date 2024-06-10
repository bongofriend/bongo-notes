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
