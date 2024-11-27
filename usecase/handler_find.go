package usecase

import (
	"errors"
	"fmt"
	"urlshortener/registry"
)

// Find implements Handler.
func (h *handler) Find(shortURL string) (redirectTarget string, err error) {
	redirectTarget, err = h.deps.Registry.Find(shortURL)
	if errors.Is(err, registry.ErrRedirectTargetNotFound) {
		return "", fmt.Errorf("%w: %w", ErrShortenedURLNotExists, err)
	}
	if err != nil {
		return "", ErrInternal
	}
	return redirectTarget, nil
}
