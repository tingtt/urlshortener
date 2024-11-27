package server

import (
	"net/http"
)

type Handler interface {
	HandleGet(w http.ResponseWriter, r *http.Request)
	HandlePost(w http.ResponseWriter, r *http.Request)
}

func newHandler(deps Dependencies) Handler {
	deps.validate()
	return &handler{deps}
}

type handler struct {
	deps Dependencies
}

const (
	MsgSystemError = "System Error. Please contact administrator."

	QueryKeyEditMode = "edit"
)
