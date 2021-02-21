package pokemon

import (
	"fmt"
	"image"

	"github.com/MaksimTheTestTaskSolver/poketask/imageCache"
	httputil "github.com/MaksimTheTestTaskSolver/poketask/util/http"
)

const pokemonApiUrlPrefix = "https://pokeapi.co/api/v2/pokemon/"

func NewService(imageCache *imageCache.ImageCache) *Service {
	return &Service{
		imageCache: imageCache,
	}
}

type Service struct {
	imageCache *imageCache.ImageCache
}

type PokemonAPIResp struct {
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

func (s *Service) GetPokemonImage(pokemonID string) (image.Image, error) {
	pokemonImage := s.imageCache.Get(pokemonID)
	if pokemonImage != nil {
		//TODO: use logger
		fmt.Println("fetching pokemon from the cache")
		return pokemonImage, nil
	}

	pokemonAPIResp := PokemonAPIResp{}
	err := httputil.Get(pokemonApiUrlPrefix+pokemonID, &pokemonAPIResp)
	if err != nil {
		return nil, fmt.Errorf("can't get data from pokemon API: %w\n", err)
	}

	if pokemonAPIResp.Sprites.FrontDefault == "" {
		return nil, fmt.Errorf("no URL in the pokemon API ressponse")
	}

	pokemonImage, err = httputil.GetImage(pokemonAPIResp.Sprites.FrontDefault)
	if err != nil {
		return nil, fmt.Errorf("can't get pokemon image: %w", err)
	}

	s.imageCache.Set(pokemonID, pokemonImage)

	return pokemonImage, nil
}
