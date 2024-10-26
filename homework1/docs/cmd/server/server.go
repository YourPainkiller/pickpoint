package main

import (
	"context"
	"homework1/internal/imcache"
	"homework1/internal/infra/kafka"
	"homework1/internal/infra/kafka/producer"
	"homework1/internal/mw"
	"homework1/internal/repository"
	"homework1/internal/repository/postgres"
	"homework1/internal/usecase"
	cliserver "homework1/pkg/cli/v1"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	kafkaConfig := newConfig()
	prod, err := producer.NewSyncProducer(kafkaConfig,
		producer.WithIdempotent(),
		producer.WithRequiredAcks(sarama.WaitForAll),
		producer.WithMaxOpenRequests(1),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
		// producer.WithProducerPartitioner(sarama.NewManualPartitioner),
		// producer.WithProducerPartitioner(sarama.NewRoundRobinPartitioner),
		// producer.WithProducerPartitioner(sarama.NewRandomPartitioner),
		producer.WithProducerPartitioner(sarama.NewHashPartitioner), // default
	)
	if err != nil {
		log.Fatal("Error in creating producer:", err)
	}
	defer prod.Close()

	orderUseCase := usecase.NewOrderUseCase(orderRepository, storageFacade, prod)

	lis, err := net.Listen("tcp", GrpcHost)
	if err != nil {
		log.Fatal("Error in listning tcp:", err)
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
		log.Fatal("Error in handling grpc:", err)
	}

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
		r.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(HttpHost, r); err != nil {
			log.Fatal("Error in listnign swagger", err)
		}
	}()

	log.Printf("Listnig grpc on: %s and http on: %s", GrpcHost, HttpHost)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Error in serving grpc:", err)
	}

}

func newStorage(pool *pgxpool.Pool) repository.Facade {
	txManager := postgres.NewTxManager(pool)
	pgRepo := postgres.NewPgRepository(txManager)
	ttlOrdersCache := imcache.NewOrdersCache(60 * time.Second)
	return repository.NewStorageFacade(*pgRepo, txManager, ttlOrdersCache)
}

func newConfig() kafka.Config {
	return kafka.Config{
		Brokers: []string{
			"localhost:9092",
		},
	}
}
