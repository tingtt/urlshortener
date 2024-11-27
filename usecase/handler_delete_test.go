package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DeleteTest struct {
	caseName         string
	in               []string
	registryBehavior RegistryBehavior
	errorIs          error
}

var deletetests = []DeleteTest{
	{
		caseName: "remove",
		in:       []string{"/0", "/a", "/b"},
		registryBehavior: RegistryBehavior{
			remove: &RegistryBehaviorRemove{nil},
		},
		errorIs: nil,
	},
	{
		caseName: "handle unexpected error",
		in:       []string{"/0", "/a", "/b"},
		registryBehavior: RegistryBehavior{
			remove: &RegistryBehaviorRemove{errors.New("unexpected error")},
		},
		errorIs: ErrInternal,
	},
}

func Test_handler_Delete(t *testing.T) {
	t.Parallel()

	for _, tt := range deletetests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			registry := new(MockRegistry)
			registry.On("Remove", mock.Anything).Return(tt.registryBehavior.remove.err)

			h := handler{Dependencies{registry}}
			err := h.Delete(tt.in...)

			assert.ErrorIs(t, err, tt.errorIs)
			if err == nil {
				registry.AssertCalled(t, "Remove", mock.Anything)
			}
		})
	}
}
