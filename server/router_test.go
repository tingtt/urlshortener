package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type RouterTest struct {
	caseName          string
	method            string
	reqURL            string
	callHandlerMethod string
}

var routertests = []RouterTest{
	{
		caseName:          "\"GET /new\" handles Handler.HandleGet",
		method:            http.MethodGet,
		reqURL:            "https://urlshortener.example/new",
		callHandlerMethod: "HandleGet",
	},
	{
		caseName:          "\"POST /exists\" handles Handler.HandlePost",
		method:            http.MethodPost,
		reqURL:            "https://urlshortener.example/exists",
		callHandlerMethod: "HandlePost",
	},
}

func Test_newRouter(t *testing.T) {
	t.Parallel()

	for _, tt := range routertests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, tt.reqURL, nil)

			mockHandler := new(MockHandler)
			mockHandler.On("HandleGet", mock.Anything, mock.Anything)
			mockHandler.On("HandlePost", mock.Anything, mock.Anything)

			router := newRouter(mockHandler)
			router.ServeHTTP(httptest.NewRecorder(), req)

			if tt.callHandlerMethod == "HandleGet" {
				mockHandler.AssertNumberOfCalls(t, "HandleGet", 1)
			}
			if tt.callHandlerMethod == "HandlePost" {
				mockHandler.AssertNumberOfCalls(t, "HandlePost", 1)
			}
		})
	}
}
