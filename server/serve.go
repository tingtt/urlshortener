package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
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

func Serve(port uint16, deps Dependencies, shutdownCtx context.Context, wg *sync.WaitGroup) error {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: newHandler(deps),
	}

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

	slog.Info("http server started on " + server.Addr)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
