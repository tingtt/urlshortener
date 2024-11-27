package server

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	uiprovider "urlshortener/ui/provider"
	"urlshortener/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type HandlerPostTest struct {
	caseName        string
	in              HandlerPostTestIn
	usecaseBehavior HandlerPostTestUsecaseBehavior
	out             HandlerGetTestOut
}

type HandlerPostTestIn struct {
	reqURL   string
	formData url.Values
}

type HandlerPostTestUsecaseBehavior struct {
	find    *UsecaseBehaviorFind
	findAll *UsecaseBehaviorFindAll
	save    *UsecaseBehaviorSave
	delete  *UsecaseBehaviorDelete
}

type UsecaseBehaviorSave struct {
	err error
}

type UsecaseBehaviorDelete struct {
	err error
}

var handlerposttest = []HandlerPostTest{
	{
		caseName: "save new shortened URL/redirect",
		in: HandlerPostTestIn{
			"https://urlshortener.example/new",
			map[string][]string{
				uiprovider.PostFormKeyRegisterShortenedURLTarget: {"https://example.test/newredirecttarget"},
				"redirect": {"on"},
			},
		},
		usecaseBehavior: HandlerPostTestUsecaseBehavior{
			save: &UsecaseBehaviorSave{nil},
		},
		out: HandlerGetTestOut{http.StatusFound, "https://example.test/newredirecttarget"},
	},
	{
		caseName: "save new shortened URL/no redirect",
		in: HandlerPostTestIn{
			"https://urlshortener.example/new",
			map[string][]string{
				uiprovider.PostFormKeyRegisterShortenedURLTarget: {"https://example.test/newredirecttarget"},
			},
		},
		usecaseBehavior: HandlerPostTestUsecaseBehavior{
			save: &UsecaseBehaviorSave{nil},
		},
		out: HandlerGetTestOut{http.StatusFound, "/new?edit"},
	},
	{
		caseName: "call save with malformed URL",
		in: HandlerPostTestIn{
			"https://urlshortener.example/new",
			map[string][]string{
				uiprovider.PostFormKeyRegisterShortenedURLTarget: {"foo.html"},
			},
		},
		usecaseBehavior: HandlerPostTestUsecaseBehavior{
			save: &UsecaseBehaviorSave{usecase.ErrMalformedURL},
		},
		out: HandlerGetTestOut{http.StatusBadRequest, ""},
	},
	{
		caseName: "handle internal error",
		in: HandlerPostTestIn{
			"https://urlshortener.example/new",
			map[string][]string{
				uiprovider.PostFormKeyRegisterShortenedURLTarget: {"https://example.test/newredirecttarget"},
			},
		},
		usecaseBehavior: HandlerPostTestUsecaseBehavior{
			save: &UsecaseBehaviorSave{usecase.ErrInternal},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "handle unexpected error",
		in: HandlerPostTestIn{
			"https://urlshortener.example/new",
			map[string][]string{
				uiprovider.PostFormKeyRegisterShortenedURLTarget: {"https://example.test/newredirecttarget"},
			},
		},
		usecaseBehavior: HandlerPostTestUsecaseBehavior{
			save: &UsecaseBehaviorSave{errors.New("unexpected error")},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "delete shortened URLs",
		in: HandlerPostTestIn{
			"https://urlshortener.example/",
			map[string][]string{
				uiprovider.PostFromKeyDeleteShortenedURLs: {"https://example.test/short1", "https://example.test/short2"},
			},
		},
		usecaseBehavior: HandlerPostTestUsecaseBehavior{
			delete: &UsecaseBehaviorDelete{nil},
		},
		out: HandlerGetTestOut{http.StatusFound, "/?edit"},
	},
	{
		caseName: "handle internal error",
		in: HandlerPostTestIn{
			"https://urlshortener.example/",
			map[string][]string{
				uiprovider.PostFromKeyDeleteShortenedURLs: {"https://example.test/short1", "https://example.test/short2"},
			},
		},
		usecaseBehavior: HandlerPostTestUsecaseBehavior{
			delete: &UsecaseBehaviorDelete{usecase.ErrInternal},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "handle unexpected error",
		in: HandlerPostTestIn{
			"https://urlshortener.example/",
			map[string][]string{
				uiprovider.PostFromKeyDeleteShortenedURLs: {"https://example.test/short1", "https://example.test/short2"},
			},
		},
		usecaseBehavior: HandlerPostTestUsecaseBehavior{
			delete: &UsecaseBehaviorDelete{errors.New("unexpected error")},
		},
		out: HandlerGetTestOut{http.StatusInternalServerError, ""},
	},
	{
		caseName: "empty form data",
		in: HandlerPostTestIn{
			"https://urlshortener.example/",
			map[string][]string{},
		},
		usecaseBehavior: HandlerPostTestUsecaseBehavior{},
		out:             HandlerGetTestOut{http.StatusBadRequest, ""},
	},
}

func Test_handler_HandlePost(t *testing.T) {
	t.Parallel()

	for _, tt := range handlerposttest {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			log.SetOutput(io.Discard)

			req := httptest.NewRequest("POST", tt.in.reqURL, strings.NewReader(tt.in.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
			if tt.usecaseBehavior.save != nil {
				usecase.On("Save", mock.Anything, mock.Anything).Return(tt.usecaseBehavior.save.err)
				defer func() {
					targetURL := tt.in.formData.Get(uiprovider.PostFormKeyRegisterShortenedURLTarget)
					usecase.AssertCalled(t, "Save", req.URL.Path, targetURL)
				}()
			}
			if tt.usecaseBehavior.delete != nil {
				usecase.On("Delete", mock.Anything).Return(tt.usecaseBehavior.delete.err)
				defer func() {
					deleteShortenedURLs := tt.in.formData[uiprovider.PostFromKeyDeleteShortenedURLs]
					usecase.AssertCalled(t, "Delete", deleteShortenedURLs)
				}()
			}

			rw := httptest.NewRecorder()

			h := handler{Dependencies{usecase, new(MockUI)}}
			h.HandlePost(rw, req)

			assert.Equal(t, tt.out.status, rw.Code)
			if tt.out.location != "" {
				assert.Equal(t, tt.out.location, rw.Header().Get("Location"))
			}
		})
	}
}

func Test_editModeURL(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "append edit query key",
			args: args{
				path: "/test",
			},
			want: "/test?edit",
		},
		{
			name: "append edit query key",
			args: args{
				path: "/test?test",
			},
			want: "/test?test&edit",
		},
		{
			name: "already contain edit query key",
			args: args{
				path: "/test?edit",
			},
			want: "/test?edit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := editModeURL(tt.args.path); got != tt.want {
				t.Errorf("editModeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
