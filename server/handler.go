package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	deps Dependencies
}

func newHandler(deps Dependencies) http.Handler {
	h := handler{deps}

	r := chi.NewRouter()
	r.Get("/*", h.HandleGet)
	r.Post("/*", h.HandlePost)

	return r
}

var (
	MsgSystemError = "System Error. Please contact administrator."
)
