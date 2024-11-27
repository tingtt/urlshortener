package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func newRouter(h Handler) http.Handler {
	r := chi.NewRouter()
	r.Get("/*", h.HandleGet)
	r.Post("/*", h.HandlePost)
	return r
}
