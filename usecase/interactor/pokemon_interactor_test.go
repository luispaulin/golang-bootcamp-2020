package interactor

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

type mockedRepository struct {
	mock.Mock
}

type mockedPresenter struct {
	mock.Mock
}

func (mr *mockedRepository) FindAll(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {
	args := mr.Called(pokemons)

	if argsPokemon := args.Get(0); argsPokemon == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.Pokemon), args.Error(1)
}

func (mr *mockedRepository) Sync() (string, int, error) {
	args := mr.Called()
	return args.String(0), args.Int(1), args.Error(2)
}

func (mp *mockedPresenter) ResponsePokemons(pokemons []*model.Pokemon) []*model.Pokemon {
	args := mp.Called(pokemons)

	if argsPokemon := args.Get(0); argsPokemon == nil {
		return nil
	}

	return args.Get(0).([]*model.Pokemon)
}

func Test_PokemonInteractor_Get(t *testing.T) {
	tests := []struct {
		name                    string
		pokemonsInput           []*model.Pokemon
		pokemonsOutput          []*model.Pokemon
		errorOutput             error
		repoPokemonsOutput      []*model.Pokemon
		repoErrorOutput         error
		presenterPokemonsOutput []*model.Pokemon
	}{
		{
			name:                    "Not errors",
			pokemonsInput:           nil,
			pokemonsOutput:          pokemonsSample,
			errorOutput:             nil,
			repoPokemonsOutput:      pokemonsSample,
			repoErrorOutput:         nil,
			presenterPokemonsOutput: pokemonsSample,
		},
		{
			name:                    "Error in repository",
			pokemonsInput:           nil,
			pokemonsOutput:          nil,
			errorOutput:             errors.New("Repository error"),
			repoPokemonsOutput:      nil,
			repoErrorOutput:         errors.New("Repository error"),
			presenterPokemonsOutput: nil,
		},
	}
	pokemonRepository := new(mockedRepository)
	pokemonPresenter := new(mockedPresenter)
	pokemonInteractor := NewPokemonInteractor(pokemonRepository, pokemonPresenter)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pokemonRepository.
				On("FindAll", tt.pokemonsInput).
				Return(tt.repoPokemonsOutput, tt.repoErrorOutput).
				Once()
			pokemonPresenter.
				On("ResponsePokemons", tt.repoPokemonsOutput).
				Return(tt.presenterPokemonsOutput).
				Once()

			result, err := pokemonInteractor.Get(tt.pokemonsInput)
			assert.Equal(t, tt.pokemonsOutput, result)
			assert.Equal(t, tt.errorOutput, err)
		})
	}
}

func Test_PokemonInteractor_Refresh(t *testing.T) {
	tests := []struct {
		name             string
		statusOutput     string
		codeOutput       int
		errorOutput      error
		repoStatusOutput string
		repoCodeOutput   int
		repoErrorOutput  error
	}{
		{
			name:             "Not errors",
			statusOutput:     "Ok",
			codeOutput:       http.StatusOK,
			errorOutput:      nil,
			repoStatusOutput: "Ok",
			repoCodeOutput:   http.StatusOK,
			repoErrorOutput:  nil,
		},
		{
			name:             "Repo sync error",
			statusOutput:     "",
			codeOutput:       0,
			errorOutput:      errors.New("Repository error"),
			repoStatusOutput: "",
			repoCodeOutput:   0,
			repoErrorOutput:  errors.New("Repository error"),
		},
	}

	pokemonRepository := new(mockedRepository)
	pokemonPresenter := new(mockedPresenter)
	pokemonInteractor := NewPokemonInteractor(pokemonRepository, pokemonPresenter)

	for _, tt := range tests {
		pokemonRepository.
			On("Sync").
			Return(tt.repoStatusOutput, tt.repoCodeOutput, tt.repoErrorOutput).
			Once()

		status, code, err := pokemonInteractor.Refresh()
		assert.Equal(t, tt.statusOutput, status)
		assert.Equal(t, tt.codeOutput, code)
		assert.Equal(t, tt.errorOutput, err)
	}
}
