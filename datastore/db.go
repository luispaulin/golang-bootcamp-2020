package datastore

import (
	"os"

	resty "github.com/go-resty/resty/v2"
)

// NewDB contructor with csv file
func NewDB() (*os.File, *resty.Client, error) {
	//TODO Set file path with a Config
	file, err := os.OpenFile("datastore/pokemons.csv", os.O_RDWR, os.ModePerm)

	//TODO HTTP Client here?
	client := resty.New()

	if err != nil {
		return nil, nil, err
	}

	return file, client, nil
}
