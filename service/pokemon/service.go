// pokemon service fetches the pokemon image from the pokemon API. It uses caching and request limiting.
package pokemon

import (
	"fmt"
	"net/http"

	"github.com/MaksimTheTestTaskSolver/poketask/imagecache"
	"github.com/MaksimTheTestTaskSolver/poketask/model"
	httputil "github.com/MaksimTheTestTaskSolver/poketask/util/http"
)

const pokemonApiUrlPrefix = "https://pokeapi.co/api/v2/pokemon/"

func NewService() *Service {
	return &Service{
		imageCache: imagecache.NewImageCache(),
	}
}

type Service struct {
	imageCache *imagecache.ImageCache
}

type PokemonAPIResp struct {
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

// GetPokemonImage returns the pokemon image by given pokemonID. It caches the fetched images and limits amount of
// parallel requests to the API
func (s *Service) GetPokemonImage(pokemonID string) (*model.Image, error) {
	pokemonImage, unlock, err := s.imageCache.GetWithLock(pokemonID)
	if err != nil {
		return nil, fmt.Errorf("can't get the pokemon image from cache with lock: %w", err)
	}

	defer unlock()

	if pokemonImage != nil {
		//TODO: use logger
		fmt.Println("fetching pokemon from the cache")
		return pokemonImage, nil
	}

	fmt.Println("pokemon image cache miss")

	pokemonImage, err = s.GetImage(pokemonID)
	if err != nil {
		return nil, err
	}

	s.imageCache.Set(pokemonID, pokemonImage)

	return pokemonImage, nil
}

func (s *Service) GetImage(pokemonID string) (*model.Image, error) {
	pokemonAPIResp := PokemonAPIResp{}
	err := httputil.Get(http.DefaultClient, pokemonApiUrlPrefix+pokemonID, &pokemonAPIResp)
	if err != nil {
		return nil, fmt.Errorf("can't get data from pokemon API: %w\n", err)
	}

	if pokemonAPIResp.Sprites.FrontDefault == "" {
		return nil, fmt.Errorf("no URL in the pokemon API ressponse")
	}

	pokemonImage, err := httputil.GetImage(http.DefaultClient, pokemonAPIResp.Sprites.FrontDefault)
	if err != nil {
		return nil, fmt.Errorf("can't get pokemon image: %w", err)
	}

	return &model.Image{ID: pokemonID, Image: pokemonImage}, nil
}
