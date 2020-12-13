package repository

import (
	"errors"
	"net/http"
	"testing"

	"github.com/luispaulin/api-challenge/domain/model"

	resty "github.com/go-resty/resty/v2"
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

type mockedSource struct {
	mock.Mock
}

func (ml *mockedSource) Get(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {
	args := ml.Called(pokemons)

	if argsPokemon := args.Get(0); argsPokemon == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.Pokemon), args.Error(1)
}

func (ml *mockedSource) Write(pokemons []*model.Pokemon) error {
	args := ml.Called(pokemons)
	return args.Error(0)
}

func Test_PokemonRepository_FindAll(t *testing.T) {
	pokemonsEmpty := make([]*model.Pokemon, 0)

	tests := []struct {
		name                string
		pokemonsInput       []*model.Pokemon
		pokemonsOutput      []*model.Pokemon
		errorOutput         error
		localSourcePokemons []*model.Pokemon
		localSourceError    error
	}{
		{
			name:                "Not source Get error",
			pokemonsInput:       pokemonsEmpty,
			pokemonsOutput:      pokemonsSample,
			errorOutput:         nil,
			localSourcePokemons: pokemonsSample,
			localSourceError:    nil,
		},
		{
			name:                "Source Get error",
			pokemonsInput:       pokemonsEmpty,
			pokemonsOutput:      nil,
			errorOutput:         errors.New("Source get error"),
			localSourcePokemons: pokemonsSample,
			localSourceError:    errors.New("Source get error"),
		},
	}

	testLocalSource := new(mockedSource)
	testRemoteSource := new(mockedSource)
	pokemonRepository := &pokemonRepository{testLocalSource, testRemoteSource}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLocalSource.
				On("Get", tt.pokemonsInput).
				Return(tt.localSourcePokemons, tt.localSourceError).
				Once()
			result, err := pokemonRepository.FindAll(tt.pokemonsInput)
			assert.Equal(t, tt.pokemonsOutput, result)
			assert.Equal(t, tt.errorOutput, err)
		})
	}
}

func Test_PokemonRepository_Sync(t *testing.T) {
	tests := []struct {
		name         string
		statusOutput string
		codeOutput   int
		errOutput    error

		remotePokemons []*model.Pokemon
		remoteError    error

		localError error
	}{
		{
			name:           "No sources errors",
			statusOutput:   "Ok",
			codeOutput:     http.StatusOK,
			errOutput:      nil,
			remotePokemons: pokemonsSample,
			remoteError:    nil,
			localError:     nil,
		},
		{
			name:           "Error in remote source's get",
			statusOutput:   "",
			codeOutput:     0,
			errOutput:      errors.New("Remote source Get error"),
			remotePokemons: nil,
			remoteError:    errors.New("Remote source Get error"),
			localError:     nil,
		},
		{
			name:           "Not successful response from source",
			statusOutput:   "Not found",
			codeOutput:     http.StatusNotFound,
			errOutput:      nil,
			remotePokemons: nil,
			remoteError:    &errorHTTP{"Not found", http.StatusNotFound},
			localError:     nil,
		},
		{
			name:           "Error in local source's write",
			statusOutput:   "",
			codeOutput:     0,
			errOutput:      errors.New("Local source write error"),
			remotePokemons: pokemonsSample,
			remoteError:    nil,
			localError:     errors.New("Local source write error"),
		},
	}

	testLocalSource := new(mockedSource)
	testRemoteSource := new(mockedSource)
	pokemonRepository := &pokemonRepository{testLocalSource, testRemoteSource}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pokemons []*model.Pokemon
			testRemoteSource.
				On("Get", pokemons).
				Return(tt.remotePokemons, tt.remoteError).
				Once()

			testLocalSource.
				On("Write", tt.remotePokemons).
				Return(tt.localError).
				Once()

			status, code, err := pokemonRepository.Sync()
			assert.Equal(t, tt.errOutput, err)
			assert.Equal(t, tt.statusOutput, status)
			assert.Equal(t, tt.codeOutput, code)
		})
	}
}

func Test_RemoteSource_Get(t *testing.T) {
	pokemonsEmpty := make([]*model.Pokemon, 0)

	tests := []struct {
		name           string
		pokemonsInput  []*model.Pokemon
		pokemonsOutput []*model.Pokemon
		errorOutput    error
	}{
		{
			name:           "Remote source succesfull get",
			pokemonsInput:  pokemonsEmpty,
			pokemonsOutput: pokemonsSample,
			errorOutput:    nil,
		},
		{
			name:           "Remote source get not found",
			pokemonsInput:  pokemonsEmpty,
			pokemonsOutput: nil,
			errorOutput:    &errorHTTP{"404 Not Found", http.StatusNotFound},
		},
	}

	client := resty.New().SetHostURL("https://pokeapi.co/api/v2/")
	remoteSrc := remoteSource{client, "pokemon"}

	for _, tt := range tests {
		if _, ok := tt.errorOutput.(*errorHTTP); ok {
			remoteSrc.getEndpoint = "poke"
		}
		result, err := remoteSrc.Get(tt.pokemonsInput)
		assert.Equal(t, tt.errorOutput, err)
		if tt.pokemonsOutput != nil {
			assert.NotNil(t, result)
		} else {
			assert.Nil(t, result)
		}
	}
}
