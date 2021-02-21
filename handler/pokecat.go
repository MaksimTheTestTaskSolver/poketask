package handler

import (
	"fmt"
	"image"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"

	"github.com/MaksimTheTestTaskSolver/poketask/imageCache"
	"github.com/MaksimTheTestTaskSolver/poketask/service/cat"
	"github.com/MaksimTheTestTaskSolver/poketask/service/pokemon"
)

func NewPokeCat(pokemonService *pokemon.Service, catService *cat.Service, imageCache *imageCache.ImageCache) *PokeCat {
	return &PokeCat{
		pokemonService: pokemonService,
		catService: catService,
		imageCache: imageCache,
	}
}

type PokeCat struct {
	pokemonService *pokemon.Service
	catService *cat.Service
	imageCache *imageCache.ImageCache
}

func (p *PokeCat) Handle(c *gin.Context) {
	pokemonID := c.Param("pokemonId")
	if pokemonID == "" {
		fmt.Printf("empty pokemonId\n")
		c.Status(http.StatusBadRequest)
		return
	}

	catImage, catID, err := p.catService.GetCatImage()
	if err != nil {
		fmt.Printf("can't get cat image: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	resultImageID := catID + pokemonID

	resultImage := p.imageCache.Get(resultImageID)
	if resultImage != nil {
		//TODO: use logger
		fmt.Println("fetching result image from the cache")
		p.encodeAndWriteToResponse(c, resultImage)
		return
	}

	pokemonImage, err := p.pokemonService.GetPokemonImage(c.Param("pokemonId"))
	if err != nil {
		fmt.Printf("can't get pokemon image: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	resultImage = p.mergeImages(catImage, pokemonImage)

	p.imageCache.Set(resultImageID, resultImage)

	p.encodeAndWriteToResponse(c, resultImage)
}

func (p *PokeCat) mergeImages(catImage image.Image, pokemonImage image.Image) image.Image {
	//TODO: find better approach to scale down the pokemon image with different aspect ration of cats image
	pokemonImageResizedWidth := catImage.Bounds().Dx() / 5

	pokemonImageSmall := imaging.Resize(pokemonImage, pokemonImageResizedWidth, 0, imaging.Lanczos)

	return imaging.Overlay(
		catImage,
		pokemonImageSmall,
		image.Pt(0, catImage.Bounds().Dy()-pokemonImageSmall.Bounds().Dy()),
		1,
	)
}

func (p *PokeCat) encodeAndWriteToResponse(c *gin.Context, resultImage image.Image) {
	err := imaging.Encode(c.Writer, resultImage, imaging.PNG)
	if err != nil {
		fmt.Printf("can't encode resulting image: %s", err)
		c.Status(http.StatusInternalServerError)
	}
}
