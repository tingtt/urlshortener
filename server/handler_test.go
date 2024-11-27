package server

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newHandler(t *testing.T) {
	t.Parallel()

	t.Run("return non-nil struct", func(t *testing.T) {
		t.Parallel()

		h := newHandler(Dependencies{new(MockUsecase), new(MockUI)})
		assert.NotNil(t, h)
	})

	t.Run("panic with not valid dependencies", func(t *testing.T) {
		t.Parallel()

		testNotValiDependencies := []Dependencies{
			{new(MockUsecase), nil},
			{nil, new(MockUI)},
			{nil, nil},
		}

		for i, deps := range testNotValiDependencies {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()

				assert.Panics(t, func() {
					newHandler(deps)
				})
			})
		}
	})
}
