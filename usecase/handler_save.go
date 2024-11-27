package usecase

import (
	"fmt"
	"net/url"
)

// Save implements Handler.
func (h *handler) Save(shortURL string, redirectTarget string) error {
	_, err := url.Parse(redirectTarget)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMalformedURL, err)
	}
	err = h.deps.Registry.Save(shortURL, redirectTarget)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return nil
}
