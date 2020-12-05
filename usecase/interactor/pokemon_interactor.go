package interactor

import (
	"github.com/luispaulin/api-challenge/domain/model"
	"github.com/luispaulin/api-challenge/usecase/presenter"
	"github.com/luispaulin/api-challenge/usecase/repository"
)

type pokemonInteractor struct {
	pokemonRepository repository.PokemonRepository
	pokemonPresenter  presenter.PokemonPresenter
}

// PokemonInteractor for pokemons use cases
type PokemonInteractor interface {
	Get(pokemons []*model.Pokemon) ([]*model.Pokemon, error)
	Refresh() (string, error)
}

// NewPokemonInteractor constructor
func NewPokemonInteractor(
	r repository.PokemonRepository,
	p presenter.PokemonPresenter) PokemonInteractor {
	return &pokemonInteractor{r, p}
}

func (po *pokemonInteractor) Get(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {
	pokemons, err := po.pokemonRepository.FindAll(pokemons)

	if err != nil {
		return nil, err
	}

	return po.pokemonPresenter.ResponsePokemons(pokemons), nil
}

//
func (po *pokemonInteractor) Refresh() (string, error) {
	message, err := po.pokemonRepository.Sync()

	if err != nil {
		return "", err
	}

	return message, nil
}
