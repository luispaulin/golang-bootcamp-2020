package datastore

import (
	"os"

	resty "github.com/go-resty/resty/v2"
)

// NewDB contructor with csv file and external API
func NewDB() (*os.File, *resty.Client, error) {
	//TODO Set file path with a Config
	file, err := os.OpenFile(
		"datastore/pokemons.csv",
		os.O_RDWR|os.O_CREATE,
		os.ModePerm,
	)

	if err != nil {
		return nil, nil, err
	}

	client := resty.New()

	return file, client, nil
}
