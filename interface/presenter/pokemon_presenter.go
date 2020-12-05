package presenter

import (
	"github.com/luispaulin/api-challenge/domain/model"
)

type pokemonPresenter struct{}

// PokemonPresenter for interface
type PokemonPresenter interface {
	ResponsePokemons(pokemons []*model.Pokemon) []*model.Pokemon
}

// NewPokemonPresenter for interface
func NewPokemonPresenter() PokemonPresenter {
	return &pokemonPresenter{}
}

func (pp *pokemonPresenter) ResponsePokemons(pokemons []*model.Pokemon) []*model.Pokemon {

	if pokemons == nil {
		pokemons = make([]*model.Pokemon, 0)
	}

	return pokemons
}
