package v1

import (
	"encoding/json"
	"log"
	"net/http"
)

func Respond(logger *log.Logger, w http.ResponseWriter, data interface{}, statusCode int, marshalIndent *string) {
	w.WriteHeader(statusCode)
	if data == nil {
		return
	}
	switch {
	case marshalIndent != nil:
		data, err := json.MarshalIndent(data, "", *marshalIndent)
		if err != nil {
			logger.Printf("error encoding the data with the marshal indent: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		w.Write(data)
	default:
		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Printf("error encoding the data: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
