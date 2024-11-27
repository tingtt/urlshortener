package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("return non-nil struct", func(t *testing.T) {
		t.Parallel()

		h := New(Dependencies{new(MockRegistry)})
		assert.NotNil(t, h)
	})

	t.Run("panic with not valid dependencies", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			New(Dependencies{nil})
		})
	})
}
