package controller

import (
	"errors"
	"net/http"
	"testing"

	"github.com/luispaulin/api-challenge/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var pokemonsSample = []*model.Pokemon{
	{
		Name: "bulbasaur",
		URL:  "https://pokeapi.co/api/v2/pokemon/1/",
	},
	{
		Name: "ivysaur",
		URL:  "https://pokeapi.co/api/v2/pokemon/2/",
	},
}

type mockedContext struct {
	mock.Mock
}

func (mc *mockedContext) JSON(code int, i interface{}) error {
	args := mc.Called(code, i)
	return args.Error(0)
}

func (mc *mockedContext) Bind(i interface{}) error {
	args := mc.Called(i)
	return args.Error(0)
}

type mockedInteractor struct {
	mock.Mock
}

func (mi *mockedInteractor) Get(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {
	args := mi.Called(pokemons)

	if argsPokemon := args.Get(0); argsPokemon == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.Pokemon), args.Error(1)
}

func (mi *mockedInteractor) Refresh() (string, int, error) {
	args := mi.Called()
	return args.String(0), args.Int(1), args.Error(2)
}

func Test_PokemonController_GetPokemons(t *testing.T) {
	tests := []struct {
		name        string
		errorOutput error

		interactPokemonsOutput []*model.Pokemon
		interactErrorOutput    error

		contextJSONErrorOutput error
	}{
		{
			name:                   "Not errors",
			errorOutput:            nil,
			interactPokemonsOutput: pokemonsSample,
			interactErrorOutput:    nil,
			contextJSONErrorOutput: nil,
		},
		{
			name:                   "Interactor error",
			errorOutput:            errors.New("Interactor error"),
			interactPokemonsOutput: nil,
			interactErrorOutput:    errors.New("Interactor error"),
			contextJSONErrorOutput: nil,
		},
		{
			name:                   "Context JSON error",
			errorOutput:            errors.New("Context JSON error"),
			interactPokemonsOutput: pokemonsSample,
			interactErrorOutput:    nil,
			contextJSONErrorOutput: errors.New("Context JSON error"),
		},
	}

	var pokemons []*model.Pokemon
	context := new(mockedContext)
	pokemonInteractor := new(mockedInteractor)
	pokemonController := NewPokemonController(pokemonInteractor)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pokemonInteractor.
				On("Get", pokemons).
				Return(tt.interactPokemonsOutput, tt.interactErrorOutput).
				Once()
			context.
				On("JSON", http.StatusOK, tt.interactPokemonsOutput).
				Return(tt.contextJSONErrorOutput).
				Once()

			err := pokemonController.GetPokemons(context)
			assert.Equal(t, tt.errorOutput, err)
		})

	}
}

func Test_PokemonController_SyncPokemons(t *testing.T) {
	tests := []struct {
		name        string
		errorOutput error

		interactMessageOutput string
		interactCodeOutput    int
		interactErrorOutput   error

		contextErrorOutput error
	}{
		{
			name:                  "Not errors",
			errorOutput:           nil,
			interactMessageOutput: "Ok",
			interactCodeOutput:    http.StatusOK,
			interactErrorOutput:   nil,
			contextErrorOutput:    nil,
		},
		{
			name:                  "Interactor error",
			errorOutput:           errors.New("Interactor error"),
			interactMessageOutput: "",
			interactCodeOutput:    0,
			interactErrorOutput:   errors.New("Interactor error"),
			contextErrorOutput:    nil,
		},
		{
			name:                  "Context JSON error",
			errorOutput:           errors.New("Context JSON error"),
			interactMessageOutput: "Ok",
			interactCodeOutput:    http.StatusOK,
			interactErrorOutput:   nil,
			contextErrorOutput:    errors.New("Context JSON error"),
		},
	}

	context := new(mockedContext)
	pokemonInteractor := new(mockedInteractor)
	pokemonController := NewPokemonController(pokemonInteractor)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pokemonInteractor.
				On("Refresh").
				Return(
					tt.interactMessageOutput,
					tt.interactCodeOutput,
					tt.interactErrorOutput,
				).Once()

			context.
				On("JSON", tt.interactCodeOutput, tt.interactMessageOutput).
				Return(tt.contextErrorOutput).
				Once()

			err := pokemonController.SyncPokemons(context)
			assert.Equal(t, tt.errorOutput, err)
		})
	}
}
