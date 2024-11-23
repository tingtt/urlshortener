package registry

import (
	"errors"
)

type Registry interface {
	Save(path string, redirectTarget string) error
	Find(path string) (redirectTarget string, err error)
	FindAll() (map[string]string, error)
	Remove(path string) error
	SavePersistently() error
	loadToCache() error
}

var (
	ErrURLNotEnough           = errors.New("url not enough")
	ErrRedirectTargetNotFound = errors.New("redirect target not found")
)
