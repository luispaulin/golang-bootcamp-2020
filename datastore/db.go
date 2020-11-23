package datastore

import (
	"os"
)

// NewDB contructor with csv file
func NewDB() (*os.File, error) {
	//TODO Set file path with a Config
	file, err := os.Open("datastore/pokemons.csv")

	if err != nil {
		return nil, err
	}

	return file, nil
}
