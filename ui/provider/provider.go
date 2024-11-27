package uiprovider

import (
	"urlshortener/usecase"

	"maragu.dev/gomponents"
)

func New() Provider {
	return &provider{&layout{}, &component{}}
}

type provider struct {
	layout    layoutI
	component componentI
}

// EditPage implements Provider.
func (p *provider) EditPage(reqURL, redirectTargetURL string, shortURLs []usecase.ShortURL) gomponents.Node {
	return p.layout.Base(
		p.component.PathInfo(reqURL, redirectTargetURL),
		gomponents.If(redirectTargetURL == "",
			p.component.RegisterForm(),
		),
		gomponents.If(redirectTargetURL != "",
			p.component.UpdateForm(),
		),
		p.component.ShortURLList(shortURLs, reqURL),
	)
}

// RegisterPage implements Provider.
func (p *provider) RegisterPage(reqURL string, shortURLs []usecase.ShortURL) gomponents.Node {
	return p.layout.Base(
		p.component.PathInfo(reqURL, "" /* redirect target not exists */),
		p.component.RegisterForm(),
		p.component.ShortURLList(shortURLs),
	)
}
