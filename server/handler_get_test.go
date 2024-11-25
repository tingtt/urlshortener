package server

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortener/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"maragu.dev/gomponents"
)

type HandlerGetTest struct {
	caseName        string
	in              HandlerGetTestIn
	usecaseBehavior HandlerGetTestUsecaseBehavior
	uiExpectCall    UIExpectCall
	uiRenderError   error
	out             HandlerGetTestOut
}

type HandlerGetTestIn struct {
	reqURL string
}

type HandlerGetTestOut struct {
	status   int
	location string
}

type HandlerGetTestUsecaseBehavior struct {
	find    *UsecaseBehaviorFind
	findAll *UsecaseBehaviorFindAll
}

type UsecaseBehaviorFind struct {
	redirectTarget string
	err            error
}

type UsecaseBehaviorFindAll struct {
	shortURLs []usecase.ShortURL
	err       error
}

type UIExpectCall struct {
	registerPage *UIExpectCallRegisterPage
	editPage     *UIExpectCallEditPage
}

type UIExpectCallRegisterPage struct {
	reqPath   string
	shortURLs []usecase.ShortURL
}

type UIExpectCallEditPage struct {
	reqPath           string
	redirectTargetURL string
	shortURLs         []usecase.ShortURL
}

var handlergettests = []HandlerGetTest{
	{
		caseName: "redirect",
		in:       HandlerGetTestIn{"https://urlshortener.example/exists"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find: &UsecaseBehaviorFind{"https://example.test/foundredirecttarget", nil},
		},
		out: HandlerGetTestOut{http.StatusFound, "https://example.test/foundredirecttarget"},
	},
	{
		caseName: "open register page with short url list that filtered with prefix matching, if short url does not exist",
		in:       HandlerGetTestIn{"https://urlshortener.example/test/none"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find: &UsecaseBehaviorFind{"", usecase.ErrShortenedURLNotExists},
			findAll: &UsecaseBehaviorFindAll{
				[]usecase.ShortURL{
					{
						From: "/test",
						To:   "https://example.test",
					},
					{
						From: "/test/exists",
						To:   "https://example.test/foundredirecttarget",
					},
					{
						From: "/testexists",
						To:   "https://example.test",
					},
					{
						From: "/tes",
						To:   "https://example.test",
					},
					{
						From: "/ignore",
						To:   "https://example.test",
					},
				},
				nil,
			},
		},
		uiExpectCall: UIExpectCall{
			registerPage: &UIExpectCallRegisterPage{
				"/test/none",
				[]usecase.ShortURL{
					{
						From: "/test",
						To:   "https://example.test",
					},
					{
						From: "/test/exists",
						To:   "https://example.test/foundredirecttarget",
					},
					{
						From: "/testexists",
						To:   "https://example.test",
					},
				},
			},
		},
		out: HandlerGetTestOut{200, ""},
	},
	{
		caseName: "handle render error/register page",
		in:       HandlerGetTestIn{"https://urlshortener.example/test/none"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find:    &UsecaseBehaviorFind{"", usecase.ErrShortenedURLNotExists},
			findAll: &UsecaseBehaviorFindAll{[]usecase.ShortURL{}, nil},
		},
		uiExpectCall: UIExpectCall{
			registerPage: &UIExpectCallRegisterPage{"/test/none", []usecase.ShortURL{}},
		},
		uiRenderError: errors.New("failed to render"),
		out:           HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "handle render error/edit page",
		in:       HandlerGetTestIn{"https://urlshortener.example/test/exists?edit"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find: &UsecaseBehaviorFind{"https://example.test/foundredirecttarget", nil},
			findAll: &UsecaseBehaviorFindAll{
				[]usecase.ShortURL{
					{
						From: "/test/exists",
						To:   "https://example.test/foundredirecttarget",
					},
				},
				nil,
			},
		},
		uiExpectCall: UIExpectCall{
			editPage: &UIExpectCallEditPage{
				"/test/exists", "https://example.test/foundredirecttarget",
				[]usecase.ShortURL{
					{
						From: "/test/exists",
						To:   "https://example.test/foundredirecttarget",
					},
				},
			},
		},
		uiRenderError: errors.New("failed to render"),
		out:           HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "open edit page with short url list that filtered with prefix matching, if contains \"edit\" query key",
		in:       HandlerGetTestIn{"https://urlshortener.example/test/exists?edit"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find: &UsecaseBehaviorFind{"https://example.test/foundredirecttarget", nil},
			findAll: &UsecaseBehaviorFindAll{
				[]usecase.ShortURL{
					{
						From: "/test",
						To:   "https://example.test",
					},
					{
						From: "/test/exists",
						To:   "https://example.test/foundredirecttarget",
					},
					{
						From: "/testexists",
						To:   "https://example.test",
					},
					{
						From: "/tes",
						To:   "https://example.test",
					},
					{
						From: "/ignore",
						To:   "https://example.test",
					},
				},
				nil,
			},
		},
		uiExpectCall: UIExpectCall{
			editPage: &UIExpectCallEditPage{
				"/test/exists", "https://example.test/foundredirecttarget",
				[]usecase.ShortURL{
					{
						From: "/test",
						To:   "https://example.test",
					},
					{
						From: "/test/exists",
						To:   "https://example.test/foundredirecttarget",
					},
					{
						From: "/testexists",
						To:   "https://example.test",
					},
				},
			},
		},
		out: HandlerGetTestOut{200, ""},
	},
	{
		caseName: "handle internal error",
		in:       HandlerGetTestIn{"https://urlshortener.example/exists"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find: &UsecaseBehaviorFind{"", usecase.ErrInternal},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "handle internal error",
		in:       HandlerGetTestIn{"https://urlshortener.example/test/none"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find:    &UsecaseBehaviorFind{"", usecase.ErrShortenedURLNotExists},
			findAll: &UsecaseBehaviorFindAll{nil, usecase.ErrInternal},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "handle internal error",
		in:       HandlerGetTestIn{"https://urlshortener.example/exists?edit"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find:    &UsecaseBehaviorFind{"https://example.test/foundredirecttarget", nil},
			findAll: &UsecaseBehaviorFindAll{nil, usecase.ErrInternal},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "handle unexpected error",
		in:       HandlerGetTestIn{"https://urlshortener.example/exists"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find: &UsecaseBehaviorFind{"", errors.New("unexpected error")},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "handle unexpected error",
		in:       HandlerGetTestIn{"https://urlshortener.example/test/none"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find:    &UsecaseBehaviorFind{"", usecase.ErrShortenedURLNotExists},
			findAll: &UsecaseBehaviorFindAll{nil, errors.New("unexpected error")},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "handle unexpected error",
		in:       HandlerGetTestIn{"https://urlshortener.example/exists?edit"},
		usecaseBehavior: HandlerGetTestUsecaseBehavior{
			find:    &UsecaseBehaviorFind{"https://example.test/foundredirecttarget", nil},
			findAll: &UsecaseBehaviorFindAll{nil, errors.New("unexpected error")},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
}

func Test_handler_HandleGet(t *testing.T) {
	t.Parallel()

	for _, tt := range handlergettests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", tt.in.reqURL, &bytes.Buffer{})

			usecase := new(MockUsecase)
			if tt.usecaseBehavior.find != nil {
				usecase.On("Find", mock.Anything).Return(tt.usecaseBehavior.find.redirectTarget, tt.usecaseBehavior.find.err)
				defer func() {
					usecase.AssertCalled(t, "Find", req.URL.Path)
				}()
			}
			if tt.usecaseBehavior.findAll != nil {
				usecase.On("FindAll").Return(tt.usecaseBehavior.findAll.shortURLs, tt.usecaseBehavior.findAll.err)
			}

			rw := httptest.NewRecorder()

			ui := new(MockUI)

			node := gomponents.Text(t.Name())

			mockNode := new(MockNode)
			mockNode.On("Render", mock.Anything).Run(func(args mock.Arguments) {
				if tt.uiRenderError == nil {
					node.Render(args.Get(0).(io.Writer))
				}
			}).Return(tt.uiRenderError)

			ui.On("RegisterPage", mock.Anything, mock.Anything).Return(mockNode)
			if tt.uiExpectCall.registerPage != nil {
				defer func() {
					ui.AssertCalled(t, "RegisterPage",
						tt.uiExpectCall.registerPage.reqPath,
						tt.uiExpectCall.registerPage.shortURLs,
					)

					if tt.uiRenderError == nil {
						buf := &bytes.Buffer{}
						node.Render(buf)
						assert.Equal(t, rw.Body.String(), buf.String())
					}
				}()
			}
			ui.On("EditPage", mock.Anything, mock.Anything, mock.Anything).Return(mockNode)
			if tt.uiExpectCall.editPage != nil {
				defer func() {
					ui.AssertCalled(t, "EditPage",
						tt.uiExpectCall.editPage.reqPath,
						tt.uiExpectCall.editPage.redirectTargetURL,
						tt.uiExpectCall.editPage.shortURLs,
					)

					if tt.uiRenderError == nil {
						buf := &bytes.Buffer{}
						node.Render(buf)
						assert.Equal(t, rw.Body.String(), buf.String())
					}
				}()
			}

			h := handler{Dependencies{usecase, ui}}
			h.HandleGet(rw, req)

			assert.Equal(t, tt.out.status, rw.Code)
			if tt.out.location != "" {
				assert.Equal(t, tt.out.location, rw.Header().Get("Location"))
			}
		})
	}
}
