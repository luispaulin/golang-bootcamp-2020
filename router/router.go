package router

import (
	"github.com/luispaulin/api-challenge/interface/controller"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewRouter to register endpoints
func NewRouter(e *echo.Echo, c controller.AppController) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Get pokemons list
	e.GET(
		"/pokemons",
		func(context echo.Context) error { return c.GetPokemons(context) },
	)

	// Start sync process
	e.GET(
		"/pokemons/sync",
		func(context echo.Context) error { return c.SyncPokemons(context) },
	)

	return e
}
