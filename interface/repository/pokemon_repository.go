package repository

import (
	"fmt"
	"net/http"
	"os"

	"github.com/luispaulin/api-challenge/domain/model"

	resty "github.com/go-resty/resty/v2"
	"github.com/gocarina/gocsv"
)

type response struct {
	Results *[]*model.Pokemon `json:"results"`
}

type errorHTTP struct {
	status string
	code   int
}

// Error string returned
func (he *errorHTTP) Error() string {
	return fmt.Sprintf("Http error: %v, %v", he.status, he.code)
}

type pokemonRepository struct {
	localSource  LocalSource
	remoteSource RemoteSource
}

type localSource struct {
	file *os.File
}

type remoteSource struct {
	client      *resty.Client
	getEndpoint string
}

// PokemonRepository for interface
type PokemonRepository interface {
	FindAll(pokemons []*model.Pokemon) ([]*model.Pokemon, error)
	Sync() (string, int, error)
}

// LocalSource for reading and writing
type LocalSource interface {
	Get(pokemons []*model.Pokemon) ([]*model.Pokemon, error)
	Write(pokemons []*model.Pokemon) error
}

// RemoteSource for just reading
type RemoteSource interface {
	Get(pokemons []*model.Pokemon) ([]*model.Pokemon, error)
}

// Get pokemons from CSV file
func (ls *localSource) Get(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {
	// Set reader at file's beginning
	if _, err := ls.file.Seek(0, 0); err != nil {
		return nil, err
	}

	// Check if empty file
	fileInfo, err := ls.file.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		return nil, nil
	}

	// Parses csv file info to pokemon slice
	if err := gocsv.UnmarshalFile(ls.file, &pokemons); err != nil {
		return nil, err
	}

	return pokemons, nil
}

// Write pokemons to CSV file
func (ls *localSource) Write(pokemons []*model.Pokemon) error {
	// Delete previous content
	if err := ls.file.Truncate(0); err != nil {
		return err
	}

	// Set writer at file's beginning
	if _, err := ls.file.Seek(0, 0); err != nil {
		return err
	}

	// Write data in file
	if err := gocsv.MarshalFile(&pokemons, ls.file); err != nil {
		return err
	}

	return nil
}

// Get pokemons from external API
func (rs *remoteSource) Get(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {
	// Create response struct
	result := response{&pokemons}

	// TODO Handle better place for api url
	// Request to external API
	resp, err := rs.client.R().
		EnableTrace().
		SetQueryString("limit=2000").
		ForceContentType("application/json").
		SetResult(&result).
		Get(rs.getEndpoint)

	if err != nil {
		return nil, err
	}

	// Check if request status successfull
	if !resp.IsSuccess() {
		return nil, &errorHTTP{resp.Status(), resp.StatusCode()}
	}

	return pokemons, nil
}

// NewPokemonRepository for interface
func NewPokemonRepository(db *os.File, client *resty.Client) PokemonRepository {
	return &pokemonRepository{
		&localSource{db},
		&remoteSource{client, "pokemon"},
	}
}

// FindAll pokemons from local source
func (pr *pokemonRepository) FindAll(pokemons []*model.Pokemon) ([]*model.Pokemon, error) {
	pokemons, err := pr.localSource.Get(pokemons)

	if err != nil {
		return nil, err
	}

	return pokemons, nil
}

// Sync Clocal source with external source
func (pr *pokemonRepository) Sync() (string, int, error) {
	var pokemons []*model.Pokemon

	pokemons, err := pr.remoteSource.Get(pokemons)

	if e, ok := err.(*errorHTTP); ok {
		return e.status, e.code, e
	} else if err != nil {
		return "", 0, err
	}

	if err := pr.localSource.Write(pokemons); err != nil {
		return "", 0, err
	}

	return "Ok", http.StatusOK, nil
}
