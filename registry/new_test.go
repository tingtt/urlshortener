package registry

import (
	"context"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("will return non-nil struct", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		wg := &sync.WaitGroup{}

		registry, err := Init(dir, context.Background(), wg)

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

				r, err := Init(dir, context.Background(), wg)

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

var _ waitgroup = new(MockWaitGroup)

type MockWaitGroup struct {
	mock.Mock
	*sync.WaitGroup
}

func (m *MockWaitGroup) Add(delta int) {
	m.Called(delta)
}

func (m *MockWaitGroup) Done() {
	m.Called()
}

func Test_standbyGraceful(t *testing.T) {
	t.Parallel()

	t.Run("call goroutine that await shutdown", func(t *testing.T) {
		t.Parallel()

		registry := new(MockFSRegistry)
		registry.On("loadToCache").Return(nil)
		registry.On("savePersistently").Return(nil)

		ctx, cancel := context.WithCancel(context.Background())

		wg := new(MockWaitGroup)
		wg.On("Add", mock.Anything)
		wg.On("Done")

		err := standbyGraceful(registry, ctx, wg)

		wg.AssertCalled(t, "Add", 1)
		registry.AssertNotCalled(t, "savePersistently")
		assert.NoError(t, err)

		cancel()
		time.Sleep(time.Millisecond)

		wg.AssertNumberOfCalls(t, "Done", 1)
		registry.AssertNumberOfCalls(t, "savePersistently", 1)
	})
}
