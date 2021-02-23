package handler

import (
	"fmt"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"

	"github.com/MaksimTheTestTaskSolver/poketask/service/cat"
	"github.com/MaksimTheTestTaskSolver/poketask/service/imagemerger"
	"github.com/MaksimTheTestTaskSolver/poketask/service/pokemon"
)

func NewPokeCat() *PokeCat {
	return &PokeCat{
		pokemonService: pokemon.NewService(),
		catService: cat.NewService(),
		imageMergerService: imagemerger.NewService(),
	}
}

type PokeCat struct {
	pokemonService     *pokemon.Service
	catService         *cat.Service
	imageMergerService *imagemerger.Service
}

func (p *PokeCat) Handle(c *gin.Context) {
	pokemonID := c.Param("pokemonId")
	if pokemonID == "" {
		fmt.Printf("empty pokemonId\n")
		c.Status(http.StatusBadRequest)
		return
	}

	catImage, err := p.catService.GetCatImage()
	if err != nil {
		fmt.Printf("can't get cat image: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	pokemonImage, err := p.pokemonService.GetPokemonImage(pokemonID)
	if err != nil {
		fmt.Printf("can't get pokemon image: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	result, err := p.imageMergerService.MergeImages(catImage, pokemonImage)
	if err != nil {
		fmt.Printf("can't merge images: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = imaging.Encode(c.Writer, result.Image, imaging.PNG)
	if err != nil {
		fmt.Printf("can't encode resulting image: %s", err)
		c.Status(http.StatusInternalServerError)
	}
}
