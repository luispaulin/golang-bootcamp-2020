package registry

import (
	"os"

	"github.com/luispaulin/api-challenge/interface/controller"

	resty "github.com/go-resty/resty/v2"
)

type registry struct {
	file   *os.File
	client *resty.Client
}

// Registry interface to initialize layers
type Registry interface {
	NewAppController() controller.AppController
}

// NewRegistry with data file
func NewRegistry(file *os.File, client *resty.Client) Registry {
	return &registry{file, client}
}

func (r *registry) NewAppController() controller.AppController {
	return r.NewPokemonController()
}
