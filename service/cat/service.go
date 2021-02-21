package cat

import (
	"fmt"
	"image"

	"github.com/MaksimTheTestTaskSolver/poketask/imageCache"
	httputil "github.com/MaksimTheTestTaskSolver/poketask/util/http"
)

const catApiUrl = "https://api.thecatapi.com/v1/images/search?mime_types=png"

func NewService(imageCache *imageCache.ImageCache) *Service {
	return &Service{
		imageCache: imageCache,
	}
}

type Service struct {
	imageCache *imageCache.ImageCache
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
