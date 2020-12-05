package repository

import (
	"os"

	resty "github.com/go-resty/resty/v2"
	"github.com/gocarina/gocsv"
	"github.com/luispaulin/api-challenge/domain/model"
)

type pokemonRepository struct {
	file   *os.File
	client *resty.Client
}

// PokemonRepository for interface
type PokemonRepository interface {
	FindAll(pokemons []*model.Pokemon) ([]*model.Pokemon, error)
	Sync() (string, int, error)
}

// NewPokemonRepository for interface
func NewPokemonRepository(db *os.File, client *resty.Client) PokemonRepository {
	return &pokemonRepository{db, client}
}

func (pr *pokemonRepository) FindAll(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {

	// Set reader t file's beginning
	if _, err := pr.file.Seek(0, 0); err != nil {
		return nil, err
	}

	// Parses csv file info to pokemon slice
	if err := gocsv.UnmarshalFile(pr.file, &pokemons); err != nil {
		return nil, err
	}

	return pokemons, nil
}

func (pr *pokemonRepository) Sync() (string, int, error) {

	var pokemons []*model.Pokemon
	// TODO this struct here?
	type Response struct {
		Results *[]*model.Pokemon `json:"results"`
	}

	// Create response struct
	result := Response{&pokemons}

	// TODO Handle better place for api url
	// Request to external API
	resp, err := pr.client.R().
		EnableTrace().
		SetQueryString("limit=2000").
		ForceContentType("application/json").
		SetResult(&result).
		Get("https://pokeapi.co/api/v2/pokemon")

	if err != nil {
		return "", 0, err
	}

	// Check if request status successfull
	if !resp.IsSuccess() {
		return resp.Status(), resp.StatusCode(), nil
	}

	if err := pr.file.Truncate(0); err != nil {
		return "", 0, err
	}

	// Set writer at file's beginning
	if _, err := pr.file.Seek(0, 0); err != nil {
		return "", 0, err
	}

	// Write collection into csv format
	if err := gocsv.MarshalFile(&pokemons, pr.file); err != nil {
		return "", 0, err
	}

	return resp.Status(), resp.StatusCode(), nil
}
