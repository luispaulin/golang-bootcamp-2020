package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo"
	"github.com/luispaulin/api-challenge/datastore"
	"github.com/luispaulin/api-challenge/registry"
	"github.com/luispaulin/api-challenge/router"
)

func main() {
	// Init file reader
	file, err := datastore.NewDB()

	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Init registry
	r := registry.NewRegistry(file)

	// Init router
	e := echo.New()
	e = router.NewRouter(e, r.NewAppController())

	fmt.Println("Server listen at http://localhost:5000")

	// Start server
	if err := e.Start(":5000"); err != nil {
		log.Fatalln(err)
	}
}
