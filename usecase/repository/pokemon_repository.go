package repository

import (
	"github.com/luispaulin/api-challenge/domain/model"
)

// PokemonRepository use case inerface
type PokemonRepository interface {
	FindAll(pokemons []*model.Pokemon) ([]*model.Pokemon, error)
}
