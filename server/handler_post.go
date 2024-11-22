package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

func (handler *handler) HandlePost(w http.ResponseWriter, r *http.Request) {
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
	fmt.Printf("redirect: %v\n", redirect)

	if redirect == "on" {
		slog.Debug(fmt.Sprintf("POST \"%s\"", r.URL.Path), slog.Int("status", http.StatusFound), slog.String("mode", "redirect"), slog.String("to", targetURL))
		http.Redirect(w, r, targetURL, http.StatusFound)
		return
	}
	slog.Debug(fmt.Sprintf("POST \"%s\"", r.URL.Path), slog.Int("status", http.StatusFound), slog.String("mode", "redirect"), slog.String("to", r.URL.Path+"?list=true"))
	http.Redirect(w, r, r.URL.Path+"?list=true", http.StatusFound)
}
