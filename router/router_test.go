package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luispaulin/api-challenge/interface/controller"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedController struct {
	mock.Mock
}

func (mc *mockedController) GetPokemons(c controller.Context) error {
	args := mc.Called(c)
	return args.Error(0)
}

func (mc *mockedController) SyncPokemons(c controller.Context) error {
	args := mc.Called(c)
	return args.Error(0)
}

func Test_EndpointsStatus(t *testing.T) {
	tests := []struct {
		name         string
		endpointCall string
		responseCode int

		controllerErrorOutput map[string]error
	}{
		{
			name:         "Get pokemons successful call",
			endpointCall: "/pokemons",
			responseCode: http.StatusOK,
			controllerErrorOutput: map[string]error{
				"GetPokemons": nil,
			},
		},
		{
			name:         "Sync pokemons successful call",
			endpointCall: "/pokemons/sync",
			responseCode: http.StatusOK,
			controllerErrorOutput: map[string]error{
				"SyncPokemons": nil,
			},
		},
		{
			name:         "Get pokemons error call",
			endpointCall: "/pokemons",
			responseCode: http.StatusInternalServerError,
			controllerErrorOutput: map[string]error{
				"GetPokemons": errors.New("Endpoint error"),
			},
		},
		{
			name:         "Sync pokemons error call",
			endpointCall: "/pokemons/sync",
			responseCode: http.StatusInternalServerError,
			controllerErrorOutput: map[string]error{
				"SyncPokemons": errors.New("Endpoint error"),
			},
		},
		{
			name:                  "Not existing endpoint",
			endpointCall:          "/poke",
			responseCode:          http.StatusNotFound,
			controllerErrorOutput: map[string]error{},
		},
	}
	e := echo.New()
	controller := new(mockedController)
	router := NewRouter(e, controller)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key := range tt.controllerErrorOutput {
				controller.
					On(key, mock.AnythingOfType("*echo.context")).
					Return(tt.controllerErrorOutput[key]).Once()
			}
			request, err := http.NewRequest("GET", tt.endpointCall, nil)
			assert.Nil(t, err)
			if err != nil {
				return
			}

			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, tt.responseCode, response.Code)
		})
	}
}
