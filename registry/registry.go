package registry

import (
	"os"

	"github.com/luispaulin/api-challenge/interface/controller"
)

type registry struct {
	file *os.File
}

// TODO: better comment
type Registry interface {
	NewAppController() controller.AppController
}

// TODO: better comment
func NewRegistry(file *os.File) Registry {
	return &registry{file}
}

func (r *registry) NewAppController() controller.AppController {
	return r.NewPokemonController()
}
