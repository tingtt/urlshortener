package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"urlshortener/cmd/clioption"
	"urlshortener/cmd/dependencies"
	"urlshortener/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
		return
	}
}

func run() error {
	cliOption, err := clioption.Load()
	if err != nil {
		return err
	}

	shutdownCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	wg := &sync.WaitGroup{}

	deps, err := dependencies.Load(cliOption.PersistentDataDirPath, shutdownCtx, wg)
	if err != nil {
		return err
	}

	err = server.Serve(cliOption.Port, deps, shutdownCtx, wg)
	wg.Wait()
	return err
}
