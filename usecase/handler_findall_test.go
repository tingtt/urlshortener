package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FindAllTest struct {
	caseName         string
	registryBehavior RegistryBehavior
	expect           FindAllTestOut
}

type FindAllTestOut struct {
	shortURLs []ShortURL
	errorIs   error
}

var findalltests = []FindAllTest{
	{
		caseName: "return sorted short URL list",
		registryBehavior: RegistryBehavior{
			findAll: &RegistryBehaviorFindAll{map[string]string{
				"/a": "https://example.test/target/a",
				"/b": "https://example.test/target/b",
				"/0": "https://example.test/target/0",
			}, nil},
		},
		expect: FindAllTestOut{[]ShortURL{
			{"/0", "https://example.test/target/0"},
			{"/a", "https://example.test/target/a"},
			{"/b", "https://example.test/target/b"},
		}, nil},
	},
	{
		caseName: "handle unexpected error",
		registryBehavior: RegistryBehavior{
			findAll: &RegistryBehaviorFindAll{nil, errors.New("unexpected error")},
		},
		expect: FindAllTestOut{nil, ErrInternal},
	},
}

func Test_handler_FindAll(t *testing.T) {
	t.Parallel()

	for _, tt := range findalltests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			registry := new(MockRegistry)
			registry.On("FindAll", mock.Anything).Return(tt.registryBehavior.findAll.shortURLs, tt.registryBehavior.findAll.err)

			h := handler{Dependencies{registry}}
			got, err := h.FindAll()

			assert.Equal(t, tt.expect.shortURLs, got)
			assert.ErrorIs(t, err, tt.expect.errorIs)
			if err == nil {
				registry.AssertCalled(t, "FindAll")
			}
		})
	}
}
