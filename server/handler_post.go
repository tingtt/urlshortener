package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func (handler *handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.PostForm.Has("delete") {
		for _, deleteShortenedURL := range r.PostForm["delete"] {
			err := handler.deps.Registry.Remove(deleteShortenedURL)
			if err != nil {
				slog.Debug(fmt.Sprintf("POST \"%s\"", r.URL.Path), slog.Int("status", http.StatusInternalServerError))
				slog.Error("failed to remove shortened URL", slog.String("path", r.URL.RawPath), slog.String("err", err.Error()))
				http.Error(w, MsgSystemError, http.StatusInternalServerError)
				return
			}
		}
		slog.Debug(fmt.Sprintf("POST \"%s\"", r.URL.Path), slog.Int("status", http.StatusFound), slog.String("to", r.URL.Path+"?list=true"), slog.String("removed", strings.Join(r.PostForm["delete"], ",")))
		http.Redirect(w, r, r.URL.Path+"?list=true", http.StatusFound)
		return
	}

	if r.PostForm.Has("target_url") {
		targetURL := r.PostFormValue("target_url")
		_, err := url.Parse(targetURL)
		if err != nil {
			slog.Debug(fmt.Sprintf("POST \"%s\"", r.URL.Path), slog.Int("status", http.StatusBadRequest), slog.String("form_value", targetURL))
			http.Error(w, "malformed URL", http.StatusBadRequest)
			return
		}

		err = handler.deps.Registry.Append(r.URL.Path, targetURL)
		if err != nil {
			slog.Debug(fmt.Sprintf("POST \"%s\"", r.URL.Path), slog.Int("status", http.StatusInternalServerError))
			slog.Error("failed to save redirect target", slog.String("path", r.URL.RawPath), slog.String("err", err.Error()))
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}

		redirect := r.PostFormValue("redirect")

		if redirect == "on" {
			slog.Debug(fmt.Sprintf("POST \"%s\"", r.URL.Path), slog.Int("status", http.StatusFound), slog.String("to", targetURL), slog.String("to", r.URL.Path+"?list=true"), slog.String("registered", r.URL.Path+" -> "+targetURL))
			http.Redirect(w, r, targetURL, http.StatusFound)
			return
		}
		slog.Debug(fmt.Sprintf("POST \"%s\"", r.URL.Path), slog.Int("status", http.StatusFound), slog.String("to", r.URL.Path+"?list=true"), slog.String("registered", r.URL.Path+" -> "+targetURL))
		http.Redirect(w, r, r.URL.Path+"?list=true", http.StatusFound)
		return
	}

	// malformed body received.
	slog.Debug(fmt.Sprintf("POST \"%s\"", r.URL.Path), slog.Int("status", http.StatusBadRequest), slog.String("cause", "malformed body"))
	http.Error(w, "malformed body", http.StatusBadRequest)
}
