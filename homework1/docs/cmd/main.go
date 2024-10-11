package main

import (
	"context"
	"homework1/internal/cli"
	"homework1/internal/dto"
	"homework1/internal/repository"
	"homework1/internal/repository/postgres"
	"homework1/internal/usecase"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	const psqlDSN = "postgres://postgres:qwe@localhost:5432/postgres?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	storageFacade := newStorage(pool)
	GenerateFakeData(500)

	orderRepository, err := repository.NewOrderRepository("./internal/database.json")
	if err != nil {
		log.Fatal(err)
	}

	orderUseCase := usecase.NewOrderUseCase(orderRepository, storageFacade)

	cli.Run(*orderUseCase)
}

func newStorage(pool *pgxpool.Pool) repository.Facade {
	txManager := postgres.NewTxManager(pool)
	pgRepo := postgres.NewPgRepository(txManager)
	return repository.NewStorageFacade(*pgRepo, txManager)
}

func GenerateFakeData(size int) {
	const psqlDSN = "postgres://postgres:qwe@localhost:5433/postgresFake?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	storageFacade := newStorage(pool)
	storageFacade.DropTable(ctx)

	for i := 0; i < size; i++ {
		var order dto.OrderDto
		err := gofakeit.Struct(&order)
		if err != nil {
			log.Fatal(err)
		}
		order.Id = i + 1
		if order.PackageType != "stretch" {
			x := rand.Intn(2)
			if x == 0 {
				order.AdditionalStretch = true
			} else {
				order.AdditionalStretch = false
			}
		}
		randHours := rand.Intn(201) - 100
		order.ValidTime = time.Now().Add(time.Hour * time.Duration(randHours)).Format(cli.TIMELAYOUT)
		storageFacade.AddOrder(ctx, order)
	}
}
