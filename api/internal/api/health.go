package api

import (
	"encoding/json"
	"net/http"
)

func (a *api) checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "available",
	}

	enc, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(enc)
}
