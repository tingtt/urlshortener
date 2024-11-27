package dependencies

import (
	"context"
	"os"
	"sync"
	"urlshortener/registry"
	"urlshortener/server"
	uiprovider "urlshortener/ui/provider"
	"urlshortener/usecase"
)

func Load(persistentDataDirPath string, shutdownCtx context.Context, wg *sync.WaitGroup) (server.Dependencies, error) {
	err := os.MkdirAll(persistentDataDirPath, 0755)
	if err != nil {
		return server.Dependencies{}, err
	}

	registry, err := registry.Init(persistentDataDirPath, shutdownCtx, wg)
	return server.Dependencies{
		Usecase: usecase.New(usecase.Dependencies{Registry: registry}),
		UI:      uiprovider.New(),
	}, err
}
