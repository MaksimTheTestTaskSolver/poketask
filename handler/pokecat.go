package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"

	"github.com/MaksimTheTestTaskSolver/poketask/model"
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

	catImage, pokemonImage, err := p.getImages(pokemonID)
	if err != nil {
		fmt.Println(err)
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
		return
	}
}

func (p *PokeCat) getImages(pokemonID string) (catImage, pokemonImage *model.Image, err error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var catError, pokemonError error

	go func() {
		defer wg.Done()
		catImage, catError = p.catService.GetCatImage()
	}()

	go func() {
		defer wg.Done()
		pokemonImage, pokemonError = p.pokemonService.GetPokemonImage(pokemonID)
	}()

	wg.Wait()

	if catError != nil {
		return nil, nil, fmt.Errorf("can't get cat image: %s", catError)
	}

	if pokemonError != nil {
		return nil, nil, fmt.Errorf("can't get pokemon image: %s", pokemonError)
	}

	return catImage, pokemonImage, nil
}