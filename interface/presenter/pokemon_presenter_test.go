package presenter

import (
	"testing"

	"github.com/luispaulin/api-challenge/domain/model"

	"github.com/stretchr/testify/assert"
)

// TODO I read that it could be Test_pokemonPresenter_ResponsePokemons
func TestResponsePokemons(t *testing.T) {
	assert := assert.New(t)

	// Test nil to empty Pokemon list
	pokemonPresenter := NewPokemonPresenter()
	pokemons := pokemonPresenter.ResponsePokemons(nil)
	emptyPokemons := make([]*model.Pokemon, 0)
	assert.Exactly(emptyPokemons, pokemons)

	// Test empty pokemon list
	pokemons = pokemonPresenter.ResponsePokemons(
		[]*model.Pokemon{},
	)
	assert.Exactly(emptyPokemons, pokemons)

	// Test pokemon filled list
	pokemons = append(
		pokemons,
		&model.Pokemon{Name: "charmander", URL: "https://pokeapi.co/api/v2/pokemon/4/"},
		&model.Pokemon{Name: "charmeleon", URL: "https://pokeapi.co/api/v2/pokemon/5/"},
	)
	pokemonsCopy := pokemonPresenter.ResponsePokemons(pokemons)
	assert.Exactly(pokemons, pokemonsCopy)
}
