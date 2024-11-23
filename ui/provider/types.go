package uiprovider

import (
	"urlshortener/usecase"

	"maragu.dev/gomponents"
)

const (
	PostFromKeyDeleteShortenedURLs        = "delete"
	PostFormKeyRegisterShortenedURLTarget = "target_url"
	PostFormKeyRedirectAfterRegister      = "redirect"
)

type Provider interface {
	RegisterPage(reqPath string, shortURLs []usecase.ShortURL) gomponents.Node
	EditPage(reqPath, redirectTargetURL string, shortURLs []usecase.ShortURL) gomponents.Node
}

type layoutI interface {
	Base(children ...gomponents.Node) gomponents.Node
}

type componentI interface {
	PathInfo(reqPath, redirectTargetURL string) gomponents.Node
	RegisterForm() gomponents.Node
	UpdateForm() gomponents.Node
	ShortURLList(shortURLs []usecase.ShortURL) gomponents.Node
}
