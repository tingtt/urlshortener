package uiprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("will return non-nil struct with non-nil fields", func(t *testing.T) {
		t.Parallel()

		provider := New().(*provider)

		assert.NotNil(t, provider)
		assert.NotNil(t, provider.layout)
		assert.NotNil(t, provider.component)
	})
}
