package usecase

import (
	"urlshortener/registry"
)

type Dependencies struct {
	Registry registry.Registry
}

func (deps *Dependencies) validate() {
	if deps.Registry == nil {
		panic("registry is nil")
	}
}

func New(deps Dependencies) Handler {
	deps.validate()
	return &handler{deps}
}

type handler struct {
	deps Dependencies
}
