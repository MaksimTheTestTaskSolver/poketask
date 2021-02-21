package handler

import (
	"fmt"
	"image"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"

	"github.com/MaksimTheTestTaskSolver/poketask/service/cat"
	"github.com/MaksimTheTestTaskSolver/poketask/service/pokemon"
)

func NewPokeCat(pokemonService *pokemon.Service, catService *cat.Service) *PokeCat {
	return &PokeCat{
		pokemonService: pokemonService,
		catService: catService,
	}
}

type PokeCat struct {
	pokemonService *pokemon.Service
	catService *cat.Service
}

func (p *PokeCat) Handle(c *gin.Context) {
	pokemonImage, err := p.pokemonService.GetPokemonImage(c.Param("pokemonId"))

	if err != nil {
		fmt.Printf("can't get pokemon image: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	catImage, err := p.catService.GetCatImage()
	if err != nil {
		fmt.Printf("can't get cat image: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	//TODO: find better approach to scale down the pokemon image with different aspect ration of cats image
	pokemonImageResizedWidth := catImage.Bounds().Dx() / 5

	pokemonImageSmall := imaging.Resize(pokemonImage, pokemonImageResizedWidth, 0, imaging.Lanczos)

	resultImage := imaging.Overlay(
		catImage,
		pokemonImageSmall,
		image.Pt(0, catImage.Bounds().Dy()-pokemonImageSmall.Bounds().Dy()),
		1,
	)

	err = imaging.Encode(c.Writer, resultImage, imaging.PNG)
	if err != nil {
		fmt.Printf("can't encode resulting image: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}
}
