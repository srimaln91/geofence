package handlers

import (
	"encoding/json"
	"net/http"
)

// Response struct
type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseJSON responds with a json response
func (resp *Response) ResponseJSON(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(resp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(response))

}

// ReponseError responds json response with error status code
func (resp *Response) ReponseError(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(resp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(400)
	w.Write([]byte(response))

}
