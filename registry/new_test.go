package registry

import (
	"context"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("will return non-nil struct", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		wg := &sync.WaitGroup{}

		registry, err := New(dir, context.Background(), wg)

		assert.NotNil(t, registry)
		assert.NoError(t, err)
	})

	t.Run("load data to cache", func(t *testing.T) {
		t.Parallel()

		for _, tt := range registrytests {
			t.Run(tt.caseName, func(t *testing.T) {
				t.Parallel()

				dir := t.TempDir()
				err := os.WriteFile(path.Join(dir, "save.csv"), []byte(tt.rawData), os.ModePerm)
				if err != nil {
					t.Fatal("failed to write file: " + err.Error())
				}
				wg := &sync.WaitGroup{}

				r, err := New(dir, context.Background(), wg)

				if tt.invalidRawData {
					assert.ErrorIs(t, err, ErrInvalidShortURL)
				} else {
					assert.NotNil(t, r)
					assert.NoError(t, err)
					assert.Equal(t, r.(*registry).data, tt.data)
				}
			})
		}
	})
}
