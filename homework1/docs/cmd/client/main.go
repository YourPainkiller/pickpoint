package main

import (
	"homework1/internal/cli"
	cliserver "homework1/pkg/cli/v1"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	psqlDSN  = "postgres://postgres:qwe@localhost:5432/postgres?sslmode=disable"
	grpcHost = "localhost:7001"
	httpHost = "localhost:7002"
)
const GrpcHost = "localhost:7001"

func main() {
	conn, err := grpc.NewClient(GrpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer conn.Close()
	cliclient := cliserver.NewCliClient(conn)

	cli.Run(cliclient)
}
