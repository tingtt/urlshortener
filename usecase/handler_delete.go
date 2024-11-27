package usecase

import "fmt"

// Delete implements Handler.
func (h *handler) Delete(shortURLs ...string) error {
	for _, deleteShortenedURL := range shortURLs {
		err := h.deps.Registry.Remove(deleteShortenedURL)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInternal, err)
		}
	}
	return nil
}
