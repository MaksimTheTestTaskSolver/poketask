package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/MaksimTheTestTaskSolver/poketask/handler"
	"github.com/MaksimTheTestTaskSolver/poketask/service/cat"
	"github.com/MaksimTheTestTaskSolver/poketask/service/pokemon"
)

const apiBasePath = "/api/v1"

func main() {
	pokemonService := pokemon.NewService()
	catService := cat.NewService()
	pokeCatHandler := handler.NewPokeCat(pokemonService, catService)

	router := gin.Default()
	apiV1Group := router.Group(apiBasePath)

	apiV1Group.GET("/pokemon/:pokemonId", pokeCatHandler.Handle)

	fmt.Println(router.Run(":8080"))
}

