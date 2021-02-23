package imagemerger

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"

	"github.com/MaksimTheTestTaskSolver/poketask/imagecache"
	"github.com/MaksimTheTestTaskSolver/poketask/model"
)

func NewService() *Service {
	return &Service{
		imageCache: imagecache.NewImageCache(),
	}
}

type Service struct {
	imageCache *imagecache.ImageCache
}

func (s *Service) MergeImages(backgroundImage, foregroundImage *model.Image) (*model.Image, error) {
	resultImageID := backgroundImage.ID + "+" + foregroundImage.ID

	result := s.imageCache.Get(resultImageID)
	if result != nil {
		//TODO: use logger
		fmt.Println("fetching result image from the cache")
		return result, nil
	}

	resultImage := s.mergeImages(backgroundImage.Image, foregroundImage.Image)

	result = &model.Image{ID: resultImageID, Image: resultImage}

	s.imageCache.Set(resultImageID, result)

	return result, nil
}

func (s *Service) mergeImages(backgroundImage image.Image, foregroundImage image.Image) image.Image {
	foregroundImageHeight := backgroundImage.Bounds().Dy() / 5

	//TODO: check resample filter
	pokemonImageSmall := imaging.Resize(foregroundImage, 0, foregroundImageHeight, imaging.Lanczos)

	return imaging.Overlay(
		backgroundImage,
		pokemonImageSmall,
		image.Pt(0, backgroundImage.Bounds().Dy()-pokemonImageSmall.Bounds().Dy()),
		1,
	)
}
