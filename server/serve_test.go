package server

import (
	"context"
	"io"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_newServer(t *testing.T) {
	t.Run("server with specified port", func(t *testing.T) {
		log.SetOutput(io.Discard)
		s := newServer(8080, newRouter(newHandler(Dependencies{new(MockUsecase), new(MockUI)})))
		assert.Equal(t, ":8080", s.Addr)
	})
}

var _ server = new(MockServer)

type MockServer struct {
	mock.Mock
}

// ListenAndServe implements server.
func (m *MockServer) ListenAndServe() error {
	args := m.Called()
	return args.Error(0)
}

// Shutdown implements server.
func (m *MockServer) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
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

func Test_gracefulServe(t *testing.T) {
	t.Parallel()

	t.Run("serve and shutdown server, if shutdown context received", func(t *testing.T) {
		t.Parallel()
		if testing.Short() {
			t.SkipNow()
		}
		log.SetOutput(io.Discard)

		server := new(MockServer)
		server.On("ListenAndServe").Return(nil)
		server.On("Shutdown", mock.Anything).Return(nil)

		ctx, cancel := context.WithCancel(context.Background())

		wg := new(MockWaitGroup)
		wg.On("Add", mock.Anything)
		wg.On("Done")

		err := gracefulServe(server, ctx, wg)

		wg.AssertCalled(t, "Add", 1)
		server.AssertNotCalled(t, "Shutdown")
		assert.NoError(t, err)

		cancel()
		time.Sleep(5*time.Second + time.Millisecond)

		wg.AssertNumberOfCalls(t, "Done", 1)
		server.AssertNumberOfCalls(t, "Shutdown", 1)
	})

	t.Run("serve", func(t *testing.T) {
		t.Parallel()

		log.SetOutput(io.Discard)

		server := new(MockServer)
		server.On("ListenAndServe").Return(nil)
		server.On("Shutdown", mock.Anything).Return(nil)

		ctx, cancel := context.WithCancel(context.Background())

		wg := new(MockWaitGroup)
		wg.On("Add", mock.Anything)
		wg.On("Done")

		err := gracefulServe(server, ctx, wg)

		wg.AssertCalled(t, "Add", 1)
		server.AssertNotCalled(t, "Shutdown")
		assert.NoError(t, err)

		cancel()
	})
}
