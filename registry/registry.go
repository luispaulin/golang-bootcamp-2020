package registry

import (
	"os"

	"github.com/luispaulin/api-challenge/interface/controller"
)

type registry struct {
	file *os.File
}

// Registry interface to initialize layers
type Registry interface {
	NewAppController() controller.AppController
}

// NewRegistry with data file
func NewRegistry(file *os.File) Registry {
	return &registry{file}
}

func (r *registry) NewAppController() controller.AppController {
	return r.NewPokemonController()
}
