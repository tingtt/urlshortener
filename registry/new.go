package registry

import (
	"context"
	"fmt"
	"log/slog"
	"path"
)

type waitgroup interface {
	Add(delta int)
	Done()
}

func Init(dataDir string, shutdownCtx context.Context, wg waitgroup) (Registry, error) {
	registry := &registry{persistentDataFilePath: path.Join(dataDir, "save.csv")}

	err := standbyGraceful(registry, shutdownCtx, wg)
	if err != nil {
		return nil, err
	}

	return registry, nil
}

func standbyGraceful(r fsregistry, shutdownCtx context.Context, wg waitgroup) error {
	err := r.loadToCache()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-shutdownCtx.Done()

		slog.Info("registry: Saving data...")
		err := r.savePersistently()
		if err != nil {
			slog.Error("failed to save data", slog.String("err", err.Error()))
			return
		}
		slog.Info("registry: Finish to save data")
	}()

	return nil
}
