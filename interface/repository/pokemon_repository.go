package repository

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/luispaulin/api-challenge/domain/model"
)

type pokemonRepository struct {
	file *os.File
}

// PokemonRepository for interface
type PokemonRepository interface {
	FindAll(pokemons []*model.Pokemon) ([]*model.Pokemon, error)
}

// NewPokemonRepository for interface
func NewPokemonRepository(db *os.File) PokemonRepository {
	return &pokemonRepository{db}
}

func (pr *pokemonRepository) FindAll(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {

	// Parses csv file info to pokemon slice
	if err := gocsv.UnmarshalFile(pr.file, &pokemons); err != nil {
		return nil, err
	}

	return pokemons, nil
}
