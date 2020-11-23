package presenter

import "github.com/luispaulin/api-challenge/domain/model"

// PokemonPresenter use case interface
type PokemonPresenter interface {
	ResponsePokemons(pokemons []*model.Pokemon) []*model.Pokemon
}
