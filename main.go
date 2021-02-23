package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/MaksimTheTestTaskSolver/poketask/handler"
)

const apiBasePath = "/api/v1"

func main() {
	pokeCatHandler := handler.NewPokeCat()

	router := gin.Default()
	apiV1Group := router.Group(apiBasePath)

	apiV1Group.GET("/pokemon/:pokemonId", pokeCatHandler.Handle)

	fmt.Println(router.Run(":8080"))
}

