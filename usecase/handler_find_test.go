package usecase

import (
	"errors"
	"testing"
	"urlshortener/registry"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FindTest struct {
	caseName         string
	in               FindTestIn
	registryBehavior RegistryBehavior
	redirectTarget   string
	errorIs          error
}

type FindTestIn struct {
	shortURL string
}

var findtests = []FindTest{
	{
		caseName: "save",
		in:       FindTestIn{"/path/to/short"},
		registryBehavior: RegistryBehavior{
			find: &RegistryBehaviorFind{"https://example.test/redirecttarget", nil},
		},
		redirectTarget: "https://example.test/redirecttarget",
		errorIs:        nil,
	},
	{
		caseName: "not found",
		in:       FindTestIn{"/path/to/short"},
		registryBehavior: RegistryBehavior{
			find: &RegistryBehaviorFind{"", registry.ErrRedirectTargetNotFound},
		},
		redirectTarget: "",
		errorIs:        ErrShortenedURLNotExists,
	},
	{
		caseName: "handle unexpected error",
		in:       FindTestIn{"/path/to/short"},
		registryBehavior: RegistryBehavior{
			find: &RegistryBehaviorFind{"", errors.New("unexpected error")},
		},
		redirectTarget: "",
		errorIs:        ErrInternal,
	},
}

func Test_handler_Find(t *testing.T) {
	t.Parallel()

	for _, tt := range findtests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			registry := new(MockRegistry)
			registry.On("Find", mock.Anything).Return(tt.registryBehavior.find.redirectTarget, tt.registryBehavior.find.err)

			h := handler{Dependencies{registry}}
			got, err := h.Find(tt.in.shortURL)

			assert.Equal(t, tt.redirectTarget, got)
			assert.ErrorIs(t, err, tt.errorIs)
			if err == nil {
				registry.AssertCalled(t, "Find", tt.in.shortURL)
			}
		})
	}
}
