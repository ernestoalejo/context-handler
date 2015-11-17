package handler

import (
	"encoding/json"
	"net/http"

	"github.com/juju/errors"
)

// ReadJSON checks the request to see if it's a POST one; and reads the JSON data.
func ReadJSON(r *http.Request, data interface{}) error {
	if r.Method != "POST" {
		return errors.New("bad method")
	}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return errors.Trace(err)
	}

	return nil
}

// WriteJSON to the response using the correct Content-Type header.
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return errors.Trace(err)
	}

	return nil
}
