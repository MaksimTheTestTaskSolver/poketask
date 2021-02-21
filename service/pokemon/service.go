package pokemon

import (
	"fmt"
	"image"

	httputil "github.com/MaksimTheTestTaskSolver/poketask/util/http"
)

const pokemonApiUrlPrefix = "https://pokeapi.co/api/v2/pokemon/"

func NewService() *Service {
	return &Service{}
}

type Service struct {
}

type PokemonAPIResp struct {
	Name    string `json:"name"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

func (s *Service) GetPokemonImage(pokemonID string) (image.Image, error) {
	pokemonAPIResp := PokemonAPIResp{}
	err := httputil.Get(pokemonApiUrlPrefix+pokemonID, &pokemonAPIResp)
	if err != nil {
		return nil, fmt.Errorf("can't get data from pokemon API: %w\n", err)
	}

	if pokemonAPIResp.Sprites.FrontDefault == "" {
		return nil, fmt.Errorf("no URL in the pokemon API ressponse")
	}

	pokemonImage, err := httputil.GetImage(pokemonAPIResp.Sprites.FrontDefault)
	if err != nil {
		return nil, fmt.Errorf("can't get pokemon image: %w", err)
	}
	return pokemonImage, nil
}
