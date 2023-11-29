package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type ResponseController struct {
}

type ResponseError struct {
	HttpStatus int    `json:"httpStatus"`
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
}

var (
	ErrInvalidParam = errors.New("invalid parameters")
)

// send a payload of JSON content
func (r *ResponseController) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logs.Error("Failed to marshal payload, error: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// send a JSON error message
func (r *ResponseController) respondWithError(w http.ResponseWriter, code int, message string) {
	r.respondWithJSON(w, code, &ResponseError{
		HttpStatus: code,
		Message:    message,
	})
}
