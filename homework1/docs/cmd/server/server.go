package main

import (
	"context"
	"homework1/internal/mw"
	"homework1/internal/repository"
	"homework1/internal/repository/postgres"
	"homework1/internal/usecase"
	cliserver "homework1/pkg/cli/v1"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	PsqlDSN   = "postgres://postgres:qwe@localhost:5432/postgres?sslmode=disable"
	GrpcHost  = "localhost:7001"
	HttpHost  = "localhost:7002"
	AdminHost = "localhost:7003"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, PsqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	storageFacade := newStorage(pool)

	orderRepository, err := repository.NewOrderRepository("./internal/database.json")
	if err != nil {
		log.Fatal(err)
	}

	orderUseCase := usecase.NewOrderUseCase(orderRepository, storageFacade)

	lis, err := net.Listen("tcp", GrpcHost)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mw.Logging), //mw.Auth
	)
	reflection.Register(grpcServer)
	cliserver.RegisterCliServer(grpcServer, orderUseCase)

	mux := runtime.NewServeMux()
	err = cliserver.RegisterCliHandlerFromEndpoint(ctx, mux, GrpcHost, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})

	if err != nil {
		log.Fatal(err)
	}
	//log.Println("all good")

	go func() {
		r := chi.NewRouter()
		r.Mount("/", mux)
		r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			b, _ := os.ReadFile("./pkg/cli/v1/cliserver.swagger.json")
			w.Write(b)
		})
		r.Get("/swagger-ui", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			b, _ := os.ReadFile("./cmd/server/static/index.html")
			w.Write(b)
		})
		if err := http.ListenAndServe(HttpHost, r); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("Listnig grpc on: %s and http on: %s", GrpcHost, HttpHost)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}

func newStorage(pool *pgxpool.Pool) repository.Facade {
	txManager := postgres.NewTxManager(pool)
	pgRepo := postgres.NewPgRepository(txManager)
	return repository.NewStorageFacade(*pgRepo, txManager)
}
