package registry

import (
	"github.com/luispaulin/api-challenge/interface/controller"
	ip "github.com/luispaulin/api-challenge/interface/presenter"
	ir "github.com/luispaulin/api-challenge/interface/repository"
	"github.com/luispaulin/api-challenge/usecase/interactor"
	up "github.com/luispaulin/api-challenge/usecase/presenter"
	ur "github.com/luispaulin/api-challenge/usecase/repository"
)

func (r *registry) NewPokemonController() controller.PokemonController {
	return controller.NewPokemonController(r.NewPokemonInteractor())
}

func (r *registry) NewPokemonInteractor() interactor.PokemonInteractor {
	return interactor.NewPokemonInteractor(r.NewPokemonRepository(), r.NewPokemonPresenter())
}

func (r *registry) NewPokemonRepository() ur.PokemonRepository {
	return ir.NewPokemonRepository(r.file)
}

func (r *registry) NewPokemonPresenter() up.PokemonPresenter {
	return ip.NewPokemonPresenter()
}
