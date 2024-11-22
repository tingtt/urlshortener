package registry

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"sync"
)

func New(dataDir string, shutdownCtx context.Context, wg *sync.WaitGroup) (Registry, error) {
	registry := &registry{persistentDataFilePath: path.Join(dataDir, "save.csv")}

	err := registry.loadToCache()
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-shutdownCtx.Done()

		slog.Info("registry: Saving data...")
		err := registry.SavePersistently()
		if err != nil {
			slog.Error("failed to save data", slog.String("err", err.Error()))
			return
		}
		slog.Info("registry: Finish to save data")
	}()

	return registry, nil
}
