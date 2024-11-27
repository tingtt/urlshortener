package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SaveTest struct {
	caseName         string
	in               SaveTestIn
	registryBehavior RegistryBehavior
	errorIs          error
}

type SaveTestIn struct {
	shortURL       string
	redirectTarget string
}

var savetests = []SaveTest{
	{
		caseName: "save",
		in:       SaveTestIn{"/path/to/short", "https://example.test/redirecttarget"},
		registryBehavior: RegistryBehavior{
			save: &RegistryBehaviorSave{err: nil},
		},
		errorIs: nil,
	},
	{
		caseName: "malformed URL",
		in: SaveTestIn{
			shortURL:       "/path/to/short",
			redirectTarget: " http://foo.com",
		},
		registryBehavior: RegistryBehavior{
			save: &RegistryBehaviorSave{err: nil},
		},
		errorIs: ErrMalformedURL,
	},
	{
		caseName: "handle unexpected error",
		in: SaveTestIn{
			shortURL:       "/path/to/short",
			redirectTarget: "https://example.test/redirecttarget",
		},
		registryBehavior: RegistryBehavior{
			save: &RegistryBehaviorSave{err: errors.New("unexpected error")},
		},
		errorIs: ErrInternal,
	},
}

func Test_handler_Save(t *testing.T) {
	t.Parallel()

	for _, tt := range savetests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			registry := new(MockRegistry)
			registry.On("Save", mock.Anything, mock.Anything).Return(tt.registryBehavior.save.err)

			h := handler{Dependencies{registry}}
			err := h.Save(tt.in.shortURL, tt.in.redirectTarget)

			assert.ErrorIs(t, err, tt.errorIs)
			if err == nil {
				registry.AssertCalled(t, "Save", tt.in.shortURL, tt.in.redirectTarget)
			}
		})
	}
}
