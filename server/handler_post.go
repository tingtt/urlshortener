package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	uiprovider "urlshortener/ui/provider"
	"urlshortener/usecase"
)

func (h *handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.PostForm.Has(uiprovider.PostFormKeyRegisterShortenedURLTarget) {
		targetURL := r.PostFormValue(uiprovider.PostFormKeyRegisterShortenedURLTarget)
		redirectAfterRegisterEnabled := r.PostFormValue(uiprovider.PostFormKeyRedirectAfterRegister) == "on"

		err := h.deps.Usecase.Save(r.URL.Path, targetURL)
		if err != nil {
			if errors.Is(err, usecase.ErrMalformedURL) {
				slog.Debug(
					fmt.Sprintf("POST \"%s\"", r.URL.Path),
					slog.Int("status", http.StatusBadRequest),
					slog.String("cause", fmt.Sprintf("malformed url (\"%s\")", targetURL)),
				)
				http.Error(w, "malformed URL", http.StatusBadRequest)
				return
			}
			slog.Error(
				fmt.Sprintf("POST \"%s\"", r.URL.Path),
				slog.Int("status", http.StatusInternalServerError),
				slog.String("error", "failed to save redirect target: "+err.Error()),
			)
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}

		if redirectAfterRegisterEnabled {
			slog.Debug(
				fmt.Sprintf("POST \"%s\"", r.URL.Path),
				slog.Int("status", http.StatusFound),
				slog.String("location", targetURL),
				slog.String("registered", r.URL.Path+" -> "+targetURL),
			)
			http.Redirect(w, r, targetURL, http.StatusFound)
			return
		}
		slog.Debug(
			fmt.Sprintf("POST \"%s\"", r.URL.Path),
			slog.Int("status", http.StatusFound),
			slog.String("location", editModeURL(r.URL.Path)),
		)
		http.Redirect(w, r, editModeURL(r.URL.Path), http.StatusFound)
		return
	}

	if r.PostForm.Has(uiprovider.PostFromKeyDeleteShortenedURLs) {
		deleteShortenedURLs := r.PostForm[uiprovider.PostFromKeyDeleteShortenedURLs]

		err := h.deps.Usecase.Delete(deleteShortenedURLs...)
		if err != nil {
			slog.Error(
				fmt.Sprintf("POST \"%s\"", r.URL.Path),
				slog.Int("status", http.StatusInternalServerError),
				slog.String("error", "failed to remove shortened URL: "+err.Error()),
			)
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}

		slog.Debug(
			fmt.Sprintf("POST \"%s\"", r.URL.Path),
			slog.Int("status", http.StatusFound),
			slog.String("location", editModeURL(r.URL.Path)),
		)
		http.Redirect(w, r, editModeURL(r.URL.Path), http.StatusFound)
		return
	}

	// malformed body received.
	slog.Debug(
		fmt.Sprintf("POST \"%s\"", r.URL.Path),
		slog.Int("status", http.StatusBadRequest),
		slog.String("cause", "malformed body"),
	)
	http.Error(w, "malformed body", http.StatusBadRequest)
}

func editModeURL(path string) string {
	if strings.Contains(path, "?"+QueryKeyEditMode) || strings.Contains(path, "&"+QueryKeyEditMode) {
		// already edit mode
		return path
	}
	if strings.Contains(path, "?") {
		return path + "&" + QueryKeyEditMode
	}
	return path + "?" + QueryKeyEditMode
}
