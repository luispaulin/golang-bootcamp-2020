package presenter

import (
	"testing"

	"github.com/luispaulin/api-challenge/domain/model"

	"github.com/stretchr/testify/assert"
)

func Test_PokemonPresenter_ResponsePokemons(t *testing.T) {
	emptyPokemons := make([]*model.Pokemon, 0)
	samplePokemons := []*model.Pokemon{
		{
			Name: "bulbasaur",
			URL:  "https://pokeapi.co/api/v2/pokemon/1/",
		},
		{
			Name: "ivysaur",
			URL:  "https://pokeapi.co/api/v2/pokemon/2/",
		},
	}

	tests := []struct {
		name     string
		pokemons []*model.Pokemon
		wanted   []*model.Pokemon
	}{
		{
			name:     "Nil pokemons input",
			pokemons: nil,
			wanted:   emptyPokemons,
		},
		{
			name:     "Empty or filled pokemons input",
			pokemons: samplePokemons,
			wanted:   samplePokemons,
		},
	}

	pokemonPresenter := NewPokemonPresenter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pokemonPresenter.ResponsePokemons(tt.pokemons)
			assert.Exactly(t, tt.wanted, result)
		})
	}
}
