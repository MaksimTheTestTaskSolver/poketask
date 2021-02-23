package cat

// cat service fetches the cat image from the cat API. It uses caching and request limiting.

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MaksimTheTestTaskSolver/poketask/imagecache"
	"github.com/MaksimTheTestTaskSolver/poketask/model"
	httputil "github.com/MaksimTheTestTaskSolver/poketask/util/http"
)

const catApiUrl = "https://api.thecatapi.com/v1/images/search?mime_types=png"

func NewService() *Service {
	return &Service{
		imageCache: imagecache.NewImageCache(),
		rlclient:   httputil.NewRLClient(10, true),
	}
}

type Service struct {
	imageCache *imagecache.ImageCache
	rlclient   *httputil.RLClient
}

type CatAPIResp []Cat

type Cat struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// GetCatImage returns a random cat image. It caches the fetched images and limits amount of parallel requests to the API
func (s *Service) GetCatImage() (*model.Image, error) {
	cat, err := s.GetCatResponse()

	if errors.Is(err, httputil.ErrQuotaReached) {
		fmt.Println("request quota for cats API reached")
		cachedImage := s.imageCache.GetRandom()
		if cachedImage == nil {
			// returned only when we have more than 10 concurrent requests without any cat images in the cache
			return nil, fmt.Errorf("too many requests")
		}
		return cachedImage, nil
	}

	if err != nil {
		return nil, err
	}

	// TODO: check why the same picture of cat has different ids
	catImage := s.imageCache.Get(cat.ID)
	if catImage != nil {
		//TODO: use logger
		fmt.Println("fetching cat from the cache")
		return catImage, nil
	}

	catImage, err = s.getCatImage(cat)
	if err != nil {
		return nil, err
	}

	s.imageCache.Set(cat.ID, catImage)

	return catImage, nil
}

func (s *Service) getCatImage(cat Cat) (*model.Image, error) {
	if cat.URL == "" {
		return nil, fmt.Errorf("no URL in the cat API ressponse")
	}

	catImage, err := httputil.GetImage(http.DefaultClient, cat.URL)
	if err != nil {
		return nil, fmt.Errorf("can't get cat image: %w", err)
	}

	return &model.Image{ID: cat.ID, Image: catImage}, nil
}

func (s *Service) GetCatResponse() (Cat, error) {
	catAPIResp := CatAPIResp{}
	err := httputil.Get(s.rlclient, catApiUrl, &catAPIResp)
	if err != nil {
		return Cat{}, fmt.Errorf("can't get data from cat API: %w\n", err)
	}

	if len(catAPIResp) == 0 {
		return Cat{}, fmt.Errorf("cat API returned an empty list\n")
	}

	firstCat := catAPIResp[0]
	if firstCat.ID == "" {
		return Cat{}, fmt.Errorf("empty id in the response from cat API")
	}
	return firstCat, nil
}
