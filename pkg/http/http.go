package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func EncodeResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err, ok := data.(error); ok {
		// TODO: derive the http status from the error
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}

func GetRouteVariable(r *http.Request, key string) string {
	vars := mux.Vars(r)

	if _, ok := vars[key]; !ok {
		return ""
	}

	return vars[key]
}
