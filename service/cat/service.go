package cat

import (
	"fmt"
	"image"

	httputil "github.com/MaksimTheTestTaskSolver/poketask/util/http"
)

const catApiUrl = "https://api.thecatapi.com/v1/images/search?mime_types=png"

func NewService() *Service {
	return &Service{}
}

type Service struct {
}

type CatAPIResp []struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (s *Service) GetCatImage() (image.Image, error) {
	catAPIResp := CatAPIResp{}
	err := httputil.Get(catApiUrl, &catAPIResp)
	if err != nil {
		return nil, fmt.Errorf("can't get data from cat API: %w\n", err)
	}

	if len(catAPIResp) == 0 {
		return nil, fmt.Errorf("cat API returned an empty list\n")
	}

	if catAPIResp[0].URL == "" {
		return nil, fmt.Errorf("no URL in the cat API ressponse")
	}

	catImage, err := httputil.GetImage(catAPIResp[0].URL)
	if err != nil {
		return nil, fmt.Errorf("can't get cat image: %w", err)
	}
	return catImage, nil
}
