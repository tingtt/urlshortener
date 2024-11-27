package registry

import (
	"errors"
)

type Registry interface {
	Save(path string, redirectTarget string) error
	Find(path string) (redirectTarget string, err error)
	FindAll() (map[string]string, error)
	Remove(path string) error
}

var (
	ErrInvalidShortURL        = errors.New("shortened URL needs to start with \"/\"")
	ErrRedirectTargetNotFound = errors.New("redirect target not found")
)
