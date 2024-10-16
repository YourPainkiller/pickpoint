package main

import (
	"context"
	"homework1/internal/cli"
	"homework1/internal/dto"
	"homework1/internal/mw"
	"homework1/internal/repository"
	"homework1/internal/repository/postgres"
	"homework1/internal/usecase"
	cliserver "homework1/pkg/cli/v1"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	psqlDSN  = "postgres://postgres:qwe@localhost:5432/postgres?sslmode=disable"
	grpcHost = "localhost:7001"
	httpHost = "localhost:7002"
)

func main() {
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

	lis, err := net.Listen("tcp", grpcHost)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.Logging),
	)
	reflection.Register(grpcServer)
	cliserver.RegisterCliServer(grpcServer, orderUseCase)

	mux := runtime.NewServeMux()
	err = cliserver.RegisterCliHandlerFromEndpoint(ctx, mux, grpcHost, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("all good")

	go func() {
		if err := http.ListenAndServe(httpHost, mux); err != nil {
			log.Fatal(err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

	//cli.Run(*orderUseCase)
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
