package main

import (
	"homework1/internal/cli"
	"homework1/internal/repository"
	"homework1/internal/usecase"
	"log"
)

func main() {
	orderRepository, err := repository.NewOrderRepository("./internal/database.json")
	if err != nil {
		log.Fatal(err)
	}
	orderUseCase := usecase.NewOrderUseCase(orderRepository)

	cli.Run(*orderUseCase)
}
