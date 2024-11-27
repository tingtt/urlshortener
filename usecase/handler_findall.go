package usecase

import (
	"fmt"
	"urlshortener/utils/tree"
)

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
