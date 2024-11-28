package uiprovider

import (
	"slices"
	"strings"
	"urlshortener/usecase"

	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

var _ componentI = &component{}

type component struct{}

// PathInfo implements componentI.
func (c *component) PathInfo(reqPath, redirectTargetURL string) gomponents.Node {
	return html.Div(
		html.P(
			gomponents.Text("Request Path: "),
			html.A(
				gomponents.Attr("href", reqPath),
				gomponents.Text(reqPath),
			),
			html.Span(gomponents.Text("  ->  ")),
			gomponents.If(redirectTargetURL == "",
				html.Span(gomponents.Text("(none)")),
			),
			gomponents.If(redirectTargetURL != "",
				html.A(
					gomponents.Attr("href", redirectTargetURL),
					gomponents.Text(redirectTargetURL),
				),
			),
		),
	)
}

// RegisterForm implements componentI.
func (c *component) RegisterForm() gomponents.Node {
	return html.Form(
		gomponents.Attr("method", "POST"),
		html.Input(
			gomponents.Attr("type", "text"),
			gomponents.Attr("name", PostFormKeyRegisterShortenedURLTarget),
			gomponents.Attr("required"),
			gomponents.Attr("autofocus"),
		),
		html.Button(
			gomponents.Attr("type", "submit"),
			gomponents.Text("Register"),
		),
		html.Br(),
		html.Input(
			gomponents.Attr("type", "checkbox"),
			gomponents.Attr("name", PostFormKeyRedirectAfterRegister),
		),
		html.Label(
			gomponents.Attr("for", PostFormKeyRedirectAfterRegister),
			gomponents.Text("Redirect"),
		),
	)
}

// UpdateForm implements componentI.
func (c *component) UpdateForm() gomponents.Node {
	return html.Form(
		gomponents.Attr("method", "POST"),
		html.Input(
			gomponents.Attr("type", "text"),
			gomponents.Attr("name", PostFormKeyRegisterShortenedURLTarget),
			gomponents.Attr("required"),
			gomponents.Attr("autofocus"),
		),
		html.Button(
			gomponents.Attr("type", "submit"),
			gomponents.Text("Update"),
		),
		html.Br(),
		html.Input(
			gomponents.Attr("type", "checkbox"),
			gomponents.Attr("name", PostFormKeyRedirectAfterRegister),
		),
		html.Label(
			gomponents.Attr("for", PostFormKeyRedirectAfterRegister),
			gomponents.Text("Redirect"),
		),
	)
}

// ShortURLList implements componentI.
func (c *component) ShortURLList(shortURLs []usecase.ShortURL, selectedShortURLs ...string) gomponents.Node {
	if len(shortURLs) == 0 {
		return gomponents.Group{
			html.H2(gomponents.Text("Redirects")),
			html.P(gomponents.Text("No redirects yet.")),
		}
	}

	return gomponents.Group{
		html.H2(gomponents.Text("Redirects")),
		html.Form(
			gomponents.Attr("method", "POST"),
			html.Button(
				gomponents.Attr("type", "submit"), html.Style("color:red;"),
				gomponents.Text("Delete selected"),
			),
			html.Table(
				html.THead(
					html.Tr(
						html.Td(),
						html.Td(gomponents.Text("From")),
						html.Td(gomponents.Text("To")),
					),
				),
				html.TBody(
					gomponents.Map(shortURLs, func(shortURL usecase.ShortURL) gomponents.Node {
						return html.Tr(
							html.Td(
								html.Input(
									gomponents.Attr("type", "checkbox"),
									gomponents.Attr("name", PostFromKeyDeleteShortenedURLs),
									gomponents.Attr("value", shortURL.From),
									gomponents.If(slices.Contains(selectedShortURLs, shortURL.From), html.Checked()),
								),
							),
							html.Td(
								html.A(
									gomponents.Attr("href", shortURL.From),
									html.Class("pr"),
									gomponents.Text(strings.Split(shortURL.From, "?")[0]),
								),
							),
							html.Td(
								html.A(
									gomponents.Attr("href", shortURL.To),
									gomponents.Attr("tabindex", "-1"),
									gomponents.Text(shortURL.To),
								),
							),
						)
					}),
				),
				html.StyleEl(
					gomponents.Text(`
						.pr {
							padding-right: 16px;
						}
						thead {
							font-weight: bold;
						}
					`),
				),
			),
		),
	}
}
