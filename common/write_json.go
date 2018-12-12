package common

import (
	"encoding/json"
	"net/http"
)

// WriteJSON encodes json data and sens it to  the response writer with a status code.
func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	serialized, err := json.Marshal(data)
	if err != nil {
		serialized, _ = json.Marshal(err.Error())
	}
	w.Write(serialized)
}
