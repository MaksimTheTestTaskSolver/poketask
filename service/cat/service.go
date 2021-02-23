package cat
// cat service fetches the cat image from the cat API. It uses caching and request limiting.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/MaksimTheTestTaskSolver/poketask/imagecache"
	"github.com/MaksimTheTestTaskSolver/poketask/model"
	"github.com/MaksimTheTestTaskSolver/poketask/requestlimiter"
	httputil "github.com/MaksimTheTestTaskSolver/poketask/util/http"
)

const catApiUrl = "https://api.thecatapi.com/v1/images/search?mime_types=png"

func NewService() *Service {
	return &Service{
		imageCache:     imagecache.NewImageCache(),
		requestLimiter: requestlimiter.NewRequestLimiter(10),
	}
}

type Service struct {
	imageCache     *imagecache.ImageCache
	requestLimiter *requestlimiter.RequestLimiter
}

type CatAPIResp []Cat

type Cat struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// GetCatImage returns a random cat image. It caches the fetched images and limits amount of parallel requests to the API
func (s *Service) GetCatImage() (*model.Image, error) {
	lockKey := strconv.Itoa(rand.Int())

	err := s.requestLimiter.AcquireLock(lockKey)
	if err != nil {
		cachedImage := s.imageCache.GetRandom()
		if cachedImage == nil {
			// returned only when we have more than 10 concurrent requests without any cat images in the cache
			return nil, fmt.Errorf("too many requests")
		}
		return cachedImage, nil
	}

	defer s.requestLimiter.FreeLock(lockKey)

	fmt.Println("calling cat API")

	cat, err := s.GetCatResponse()
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

	catImage, err := httputil.GetImage(cat.URL)
	if err != nil {
		return nil, fmt.Errorf("can't get cat image: %w", err)
	}

	return &model.Image{ID: cat.ID, Image: catImage}, nil
}

func (s *Service) GetCatResponse() (Cat, error) {
	catAPIResp := CatAPIResp{}
	err := httputil.Get(catApiUrl, &catAPIResp)
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
