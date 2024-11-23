package usecase

import "errors"

var (
	ErrShortenedURLNotExists = errors.New("shortened URL not exists")
	ErrMalformedURL          = errors.New("malformed URL")
	ErrInternal              = errors.New("internal error")
)

type Handler interface {
	Save(shortURL string, redirectTarget string) error
	Find(shortURL string) (redirectTarget string, err error)
	FindAll() ([]ShortURL, error)
	Delete(shortURLs ...string) error
}

type ShortURL struct {
	From string
	To   string
}
