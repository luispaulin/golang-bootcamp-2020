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
	Sync() (string, error)
}

// NewPokemonRepository for interface
func NewPokemonRepository(db *os.File, client *resty.Client) PokemonRepository {
	return &pokemonRepository{db, client}
}

func (pr *pokemonRepository) FindAll(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {

	if _, err := pr.file.Seek(0, 0); err != nil {
		return nil, err
	}

	// Parses csv file info to pokemon slice
	if err := gocsv.UnmarshalFile(pr.file, &pokemons); err != nil {
		return nil, err
	}

	return pokemons, nil
}

func (pr *pokemonRepository) Sync() (string, error) {

	var pokemons []*model.Pokemon
	// TODO this struct here?
	type Response struct {
		Results *[]*model.Pokemon `json:"results"`
	}

	result := Response{&pokemons}

	// TODO Handle better place for api url
	resp, err := pr.client.R().
		EnableTrace().
		SetQueryString("limit=2000").
		ForceContentType("application/json").
		SetResult(&result).
		Get("https://pokeapi.co/api/v2/pokemon")

	if err != nil {
		return "", err
	}

	// TODO correct http error status raising?
	if !resp.IsSuccess() {
		return resp.Status(), nil
	}

	if err != nil {
		return "", err
	}

	err = pr.file.Truncate(0)

	if err != nil {
		return "", err
	}

	csvContent, err := gocsv.MarshalString(&pokemons)

	if err != nil {
		return "", err
	}

	if _, err := pr.file.Seek(0, 0); err != nil {
		return "", err
	}

	//err = gocsv.MarshalFile(&pokemons, pr.file)
	_, err = pr.file.WriteString(csvContent)

	if err != nil {
		return "", err
	}

	return resp.Status(), nil
}
