package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"urlshortener/registry"
	"urlshortener/utils/tree"

	"maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func (handler *handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	mode := inspeceMode(r)

	if mode == modeList {
		redirectMap, err := handler.deps.Registry.FindAll()
		if err != nil {
			slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusInternalServerError))
			slog.Error("failed to find redirect target", slog.String("path", r.URL.Path), slog.String("err", err.Error()))
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}
		redirects := sortAndFilterRedirects(r.URL.Path, redirectMap)

		html := layoutHTML(listHTML(redirects))
		err = html.Render(w)
		if err != nil {
			slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusInternalServerError))
			slog.Error("failed render html", slog.String("path", r.URL.Path), slog.String("err", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusOK), slog.String("mode", "list"))
		return
	}

	targetURL, err := handler.deps.Registry.Find(r.URL.Path)
	if /* not found or mode is replace */ errors.Is(err, registry.ErrRedirectTargetNotFound) || mode == modeReplace {
		if mode == modeReplace {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

		redirectMap, err := handler.deps.Registry.FindAll()
		if err != nil {
			slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusInternalServerError))
			slog.Error("failed to find redirect target", slog.String("path", r.URL.Path), slog.String("err", err.Error()))
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}
		redirects := sortAndFilterRedirects(r.URL.Path, redirectMap)

		html := layoutHTML(registerFormHTML(r.URL.Path, mode == modeReplace), listHTML(redirects))
		err = html.Render(w)
		if err != nil {
			slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusInternalServerError))
			slog.Error("failed render html", slog.String("path", r.URL.Path), slog.String("err", err.Error()))
			http.Error(w, MsgSystemError, http.StatusInternalServerError)
			return
		}
		if mode == modeReplace {
			slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusOK), slog.String("mode", "replace"))
		} else {
			slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusNotFound), slog.String("mode", "append"))
		}
		return
	}
	if err != nil {
		slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusInternalServerError))
		slog.Error("failed to find redirect target", slog.String("path", r.URL.Path), slog.String("err", err.Error()))
		http.Error(w, MsgSystemError, http.StatusInternalServerError)
		return
	}

	slog.Debug(fmt.Sprintf("GET \"%s\"", r.URL.Path), slog.Int("status", http.StatusFound), slog.String("mode", "redirect"), slog.String("to", targetURL))
	http.Redirect(w, r, targetURL, http.StatusFound)
}

const (
	modeNone int = iota
	modeReplace
	modeList
)

func inspeceMode(r *http.Request) int {
	query := r.URL.Query()

	listQuery := query.Get("list")
	if listQuery == "1" || listQuery == "true" {
		return modeList
	}

	replaceQuery := query.Get("replace")
	if replaceQuery == "1" || replaceQuery == "true" {
		return modeReplace
	}
	return modeNone
}

func sortAndFilterRedirects(reqPath string, redirectMap map[string]string) []redirect {
	var root *tree.Node[redirect]
	cnt := 0
	for from, to := range redirectMap {
		if !strings.HasPrefix(from, path.Dir(reqPath)) {
			continue
		}
		cnt++
		root = tree.Insert(root, redirect{from, to}, func(new, curr redirect) (isLeft bool) {
			return new.from < curr.from
		})
	}
	var redirects = make([]redirect, 0, cnt)
	tree.InOrderTraversal(root, &redirects)
	return redirects
}

func layoutHTML(children ...gomponents.Node) gomponents.Node {
	return html.HTML(
		html.H1(gomponents.Text("URL Shortener")),
		gomponents.Group(children),
	)
}

func registerFormHTML(reqPath string, isReplaceMode bool) gomponents.Node {
	return gomponents.Group{
		html.Div(
			html.P(gomponents.Textf("Request Path: "), html.A(gomponents.Attr("href", reqPath), gomponents.Text(reqPath))),
			gomponents.If(!isReplaceMode,
				html.P(gomponents.Textf("Redirect target not found.")),
			),
		),
		html.Form(gomponents.Attr("method", "POST"),
			html.Input(gomponents.Attr("type", "text"), gomponents.Attr("name", "target_url"), gomponents.Attr("required"), gomponents.Attr("autofocus")),
			html.Button(gomponents.Attr("type", "submit"),
				gomponents.Text("Save"),
			),
			html.Br(),
			html.Input(gomponents.Attr("type", "checkbox"), gomponents.Attr("name", "redirect")),
			html.Label(gomponents.Attr("for", "redirect"), gomponents.Text("Redirect")),
		),
	}
}

type redirect struct {
	from string
	to   string
}

func listHTML(redirects []redirect) gomponents.Node {
	if len(redirects) == 0 {
		return gomponents.Group{
			html.H2(gomponents.Text("Redirects")),
			html.P(gomponents.Text("No redirects yet.")),
		}
	}

	return html.Table(
		html.THead(
			html.Tr(
				html.Td(gomponents.Text("From")),
				html.Td(gomponents.Text("To")),
			),
		),
		html.TBody(
			gomponents.Map(redirects, func(redirect redirect) gomponents.Node {
				return html.Tr(
					html.Td(
						html.A(gomponents.Attr("href", redirect.from), html.Class("pr"),
							gomponents.Text(redirect.from),
						),
					),
					html.Td(
						html.A(gomponents.Attr("href", redirect.to), gomponents.Attr("tabindex", "-1"),
							gomponents.Text(redirect.to),
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
	)
}
