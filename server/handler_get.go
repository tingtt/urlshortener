package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"
	"urlshortener/usecase"
)

func (handler *handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	redirectTargetURL, err := handler.deps.Usecase.Find(r.URL.Path)
	if err != nil && !errors.Is(err, usecase.ErrShortenedURLNotExists) {
		slog.Error(
			fmt.Sprintf("POST \"%s\"", r.URL.Path),
			slog.Int("status", http.StatusInternalServerError),
			slog.String("error", "failed to find redirect target: "+err.Error()),
		)
		http.Error(w, MsgSystemError, http.StatusInternalServerError)
		return
	}

	if isEditMode(r.URL.Query()) {
		shortURLs, err := handler.deps.Usecase.FindAll()
		if err != nil {
			slog.Error(
				fmt.Sprintf("POST \"%s\"", r.URL.Path),
				slog.Int("status", http.StatusInternalServerError),
				slog.String("error", "failed to find redirect target: "+err.Error()),
			)
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}

		html := handler.deps.UI.EditPage(
			r.URL.Path,
			redirectTargetURL,
			filterShortURLsByPrefix(shortURLs, path.Dir(r.URL.Path)),
		)
		err = html.Render(w)
		if err != nil {
			slog.Error(
				fmt.Sprintf("POST \"%s\"", r.URL.Path),
				slog.Int("status", http.StatusInternalServerError),
				slog.String("error", "failed to render html: "+err.Error()),
			)
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}

		slog.Debug(
			fmt.Sprintf("GET \"%s\"", r.URL.Path),
			slog.Int("status", http.StatusOK),
		)
		w.WriteHeader(http.StatusOK)
		return
	}

	if /* shortened URL not found */ errors.Is(err, usecase.ErrShortenedURLNotExists) {
		shortURLs, err := handler.deps.Usecase.FindAll()
		if err != nil {
			slog.Error(
				fmt.Sprintf("POST \"%s\"", r.URL.Path),
				slog.Int("status", http.StatusInternalServerError),
				slog.String("error", "failed to find redirect target: "+err.Error()),
			)
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}

		html := handler.deps.UI.RegisterPage(
			r.URL.Path,
			filterShortURLsByPrefix(shortURLs, path.Dir(r.URL.Path)),
		)
		err = html.Render(w)
		if err != nil {
			slog.Error(
				fmt.Sprintf("POST \"%s\"", r.URL.Path),
				slog.Int("status", http.StatusInternalServerError),
				slog.String("error", "failed to render html: "+err.Error()),
			)
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}

		slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusNotFound), slog.String("mode", "append"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	slog.Debug(
		fmt.Sprintf("GET \"%s\"", r.URL.Path),
		slog.Int("status", http.StatusFound),
		slog.String("location", redirectTargetURL),
	)
	http.Redirect(w, r, redirectTargetURL, http.StatusFound)
}

func filterShortURLsByPrefix(shortURLs []usecase.ShortURL, prefix string) []usecase.ShortURL {
	filtered := make([]usecase.ShortURL, 0, len(shortURLs))
	for _, shortURL := range shortURLs {
		if strings.HasPrefix(shortURL.From, prefix) {
			filtered = append(filtered, shortURL)
		}
	}
	return filtered
}

func isEditMode(query url.Values) bool {
	return query.Has(QueryKeyEditMode)
}
