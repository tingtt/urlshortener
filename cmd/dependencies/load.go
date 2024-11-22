package dependencies

import (
	"context"
	"os"
	"sync"
	"urlshortener/registry"
	"urlshortener/server"
)

func Load(persistentDataDirPath string, shutdownCtx context.Context, wg *sync.WaitGroup) (server.Dependencies, error) {
	err := os.MkdirAll(persistentDataDirPath, 0755)
	if err != nil {
		return server.Dependencies{}, err
	}

	registry, err := registry.New(persistentDataDirPath, shutdownCtx, wg)
	return server.Dependencies{Registry: registry}, err
}
