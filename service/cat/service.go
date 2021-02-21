package cat

import (
	"fmt"
	"image"

	"github.com/MaksimTheTestTaskSolver/poketask/imagecache"
	"github.com/MaksimTheTestTaskSolver/poketask/requestlimiter"
	httputil "github.com/MaksimTheTestTaskSolver/poketask/util/http"
)

const catApiUrl = "https://api.thecatapi.com/v1/images/search?mime_types=png"

func NewService() *Service {
	return &Service{
		imageCache: imagecache.NewImageCache(),
		requestLimiter: requestlimiter.NewRequestLimiter(10),
	}
}

type Service struct {
	imageCache *imagecache.ImageCache
	requestLimiter *requestlimiter.RequestLimiter
}

type CatAPIResp []struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
}

func (s *Service) GetCatImage() (image image.Image, catID string, err error) {
	catAPIResp := CatAPIResp{}
	err = httputil.Get(catApiUrl, &catAPIResp)
	if err != nil {
		return nil, "", fmt.Errorf("can't get data from cat API: %w\n", err)
	}

	if len(catAPIResp) == 0 {
		return nil, "", fmt.Errorf("cat API returned an empty list\n")
	}

	firstCat := catAPIResp[0]
	if firstCat.ID == "" {
		return nil, "", fmt.Errorf("empty id in the response from cat API")
	}

	catImage := s.imageCache.Get(firstCat.ID)
	if catImage != nil {
		//TODO: use logger
		fmt.Println("fetching cat from the cache")
		return catImage, "", nil
	}

	err = s.requestLimiter.AcquireLock(firstCat.ID)
	if err == requestlimiter.ErrQuotaReached {
		catID, cachedImage := s.imageCache.GetRandom()
		return cachedImage, catID, err
	}

	if err == requestlimiter.ErrLockAlreadyAcquired {
		fmt.Println("was in a waiting queue")
		return s.imageCache.Get(firstCat.ID), firstCat.ID, nil
	}

	fmt.Println("calling cat API")
	defer s.requestLimiter.FreeLock(firstCat.ID)

	if firstCat.URL == "" {
		return nil, "", fmt.Errorf("no URL in the cat API ressponse")
	}

	catImage, err = httputil.GetImage(firstCat.URL)
	if err != nil {
		return nil, "", fmt.Errorf("can't get cat image: %w", err)
	}

	s.imageCache.Set(firstCat.ID, catImage)

	return catImage, firstCat.ID, nil
}
