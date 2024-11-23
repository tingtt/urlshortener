package uiprovider

import (
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

var _ layoutI = &layout{}

type layout struct{}

// Base implements layoutI.
func (l *layout) Base(children ...gomponents.Node) gomponents.Node {
	return html.HTML(
		html.H1(gomponents.Text("URL Shortener")),
		gomponents.Group(children),
	)
}
