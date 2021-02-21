package pokemon

import (
	"fmt"
	"image"

	"github.com/MaksimTheTestTaskSolver/poketask/imagecache"
	"github.com/MaksimTheTestTaskSolver/poketask/requestlimiter"
	httputil "github.com/MaksimTheTestTaskSolver/poketask/util/http"
)

const pokemonApiUrlPrefix = "https://pokeapi.co/api/v2/pokemon/"

func NewService() *Service {
	return &Service{
		imageCache: imagecache.NewImageCache(),
		requestLimiter: requestlimiter.NewRequestLimiter(0),
	}
}

type Service struct {
	imageCache *imagecache.ImageCache
	requestLimiter *requestlimiter.RequestLimiter
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

	err := s.requestLimiter.AcquireLock(pokemonID)
	if err == requestlimiter.ErrQuotaReached {
		return nil, err
	}

	if err == requestlimiter.ErrLockAlreadyAcquired {
		fmt.Println("was in a waiting queue")
		return s.imageCache.Get(pokemonID), nil
	}

	fmt.Println("calling pokemon API")
	defer s.requestLimiter.FreeLock(pokemonID)

	pokemonAPIResp := PokemonAPIResp{}
	err = httputil.Get(pokemonApiUrlPrefix+pokemonID, &pokemonAPIResp)
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
