// Package cli
package main

import (
	"fmt"

	database "github.com/anderson89marques/roxb3/internal/adapters/drivens/postgres"
	"github.com/anderson89marques/roxb3/internal/infra/config"
	services "github.com/anderson89marques/roxb3/internal/services"
)

func main() {
	fmt.Println("Hello World")
	err := config.ParseEnv()
	if err != nil {
		panic(err)
	}
	repo, err := database.NewRepository(config.GetEnv())
	if err != nil {
		panic(err)
	}
	stockService := services.NewStockService(repo, config.GetEnv())
	err = stockService.IngestData()
	if err != nil {
		panic(err)
	}
}
