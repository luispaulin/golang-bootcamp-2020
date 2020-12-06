package registry

import (
	"github.com/luispaulin/api-challenge/interface/controller"
	ip "github.com/luispaulin/api-challenge/interface/presenter"
	ir "github.com/luispaulin/api-challenge/interface/repository"
	"github.com/luispaulin/api-challenge/usecase/interactor"
	up "github.com/luispaulin/api-challenge/usecase/presenter"
	ur "github.com/luispaulin/api-challenge/usecase/repository"
)

// NewPokemonController creation
func (r *registry) NewPokemonController() controller.PokemonController {
	return controller.NewPokemonController(r.NewPokemonInteractor())
}

// NewPokemonInteractor creation
func (r *registry) NewPokemonInteractor() interactor.PokemonInteractor {
	return interactor.NewPokemonInteractor(
		r.NewPokemonRepository(),
		r.NewPokemonPresenter(),
	)
}

// NewPokemonRepository creation
func (r *registry) NewPokemonRepository() ur.PokemonRepository {
	return ir.NewPokemonRepository(r.file, r.client)
}

// NewPokemonPresenter creation
func (r *registry) NewPokemonPresenter() up.PokemonPresenter {
	return ip.NewPokemonPresenter()
}
