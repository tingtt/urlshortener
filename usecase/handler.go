package usecase

import (
	"errors"
	"fmt"
	"net/url"
	"urlshortener/registry"
	"urlshortener/utils/tree"
)

type Dependencies struct {
	Registry registry.Registry
}

func (deps *Dependencies) validate() {
	if deps.Registry == nil {
		panic("registry is nil")
	}
}

func New(deps Dependencies) Handler {
	deps.validate()
	return &handler{deps}
}

type handler struct {
	deps Dependencies
}

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

// Find implements Handler.
func (h *handler) Find(shortURL string) (redirectTarget string, err error) {
	redirectTarget, err = h.deps.Registry.Find(shortURL)
	if errors.Is(err, registry.ErrRedirectTargetNotFound) {
		return "", fmt.Errorf("%w: %w", ErrShortenedURLNotExists, err)
	}
	return redirectTarget, nil
}

// FindAll implements Handler.
func (h *handler) FindAll() ([]ShortURL, error) {
	redirectMap, err := h.deps.Registry.FindAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// sort shortened URLs
	var root *tree.Node[ShortURL]
	cnt := 0
	for from, to := range redirectMap {
		cnt++
		root = tree.Insert(root, ShortURL{from, to}, func(new, curr ShortURL) (isLeft bool) {
			return new.From < curr.From
		})
	}
	var shortenedURLs = make([]ShortURL, 0, cnt)
	tree.InOrderTraversal(root, &shortenedURLs)
	return shortenedURLs, nil
}

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
