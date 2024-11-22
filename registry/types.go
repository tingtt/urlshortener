package registry

import (
	"errors"
)

type Registry interface {
	Append(path string, redirectTarget string) error
	Find(path string) (redirectTarget string, err error)
	FindAll() (map[string]string, error)
	SavePersistently() error
	loadToCache() error
}

var (
	ErrURLNotEnough           = errors.New("url not enough")
	ErrRedirectTargetNotFound = errors.New("redirect target not found")
)
