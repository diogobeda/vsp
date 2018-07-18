package web

import (
	"io"
	"net/http"
)

type iHandler interface {
	Ok(w io.Writer)
	Created(w io.Writer)
	BadRequest(w io.Writer, message string)
	Internal(w io.Writer, message string)
	NotFound(w io.Writer, message string)
}

type WebHandler struct {
	iHandler
}

func (wh WebHandler) Created(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func (wh WebHandler) BadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, message)
}

func (wh WebHandler) Internal(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, message)
}

func (wh WebHandler) NotFound(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, message)
}

func (wh WebHandler) Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
