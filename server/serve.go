package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	uiprovider "urlshortener/ui/provider"
	"urlshortener/usecase"
)

type Dependencies struct {
	Usecase usecase.Handler
	UI      uiprovider.Provider
}

func (deps *Dependencies) validate() {
	if deps.Usecase == nil {
		panic("handler is nil")
	}
	if deps.UI == nil {
		panic("UI provider is nil")
	}
}

type waitgroup interface {
	Add(delta int)
	Done()
}

func Serve(port uint16, deps Dependencies, shutdownCtx context.Context, wg waitgroup) error {
	server := newServer(port, newRouter(newHandler(deps)))

	slog.Info("http server started on " + server.Addr)
	return gracefulServe(server, shutdownCtx, wg)
}

func newServer(port uint16, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
}

type server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func gracefulServe(server server, shutdownCtx context.Context, wg waitgroup) error {
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-shutdownCtx.Done()

		slog.Info("http: Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		slog.Info("http: Server closed")
	}()

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
