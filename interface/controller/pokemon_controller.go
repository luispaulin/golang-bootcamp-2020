package controller

import (
	"net/http"

	"github.com/luispaulin/api-challenge/domain/model"
	"github.com/luispaulin/api-challenge/usecase/interactor"
)

type pokemonController struct {
	pokemonInteractor interactor.PokemonInteractor
}

// PokemonController interface
type PokemonController interface {
	GetPokemons(c Context) error
	SyncPokemons(c Context) error
}

// NewPokemonController constructor
func NewPokemonController(pin interactor.PokemonInteractor) PokemonController {
	return &pokemonController{pin}
}

func (pc *pokemonController) GetPokemons(c Context) error {
	var pokemons []*model.Pokemon

	pokemons, err := pc.pokemonInteractor.Get(pokemons)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pokemons)
}

func (pc *pokemonController) SyncPokemons(c Context) error {
	message, code, err := pc.pokemonInteractor.Refresh()

	if err != nil {
		return err
	}

	return c.JSON(code, message)
}
