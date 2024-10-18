OUT_PATH:=$(CURDIR)/pkg
LOCAL_BIN:=$(CURDIR)/bin

APP1=client.exe
APP_EXECUTABLE1="./cmd/client/$(APP1)"

APP2=server.exe
APP_EXECUTABLE2="./cmd/server/$(APP2)"

SHELL := /bin/bash 

set-gobin:
	set GOBIN=$(LOCAL_BIN)

check-quality: 
	make fmt
	make vet
	make lint


lint: 
	gocognit -over 5 . || exit 0
	gocyclo -over 5 . || exit 0
	squawk ./migrations/*.sql || exit 0

cover:
	go test -cover ./internal/usecase
	

vet: 
	go vet ./...

fmt: 
	go fmt ./...

tidy: 
	go get -u
	go mod tidy


# ---------------------------
# Запуск сервера и клиента CLI 
# ---------------------------
build-client:
	make lint
	go build -o $(APP_EXECUTABLE1) ./cmd/client/main.go
	@echo "Client Build passed"

run-client: 
	make build-client
	$(APP_EXECUTABLE1)

build-server:
	go build -o $(APP_EXECUTABLE2) ./cmd/server/server.go
	@echo "Server Build passed"

run-server: 
	make build-server
	$(APP_EXECUTABLE2)

# ---------------------------
# Установка зависимостей, генирация и запуск сервера
# ---------------------------

all:
	@if [ -d vendor.protogen ]; then rm -rf vendor.protogen/*; fi
	make bin-deps
	make generate
	make run-server

test-int:
	make compose-up
	make goose-down-fake
	make goose-up-fake
	go test ./test/cmd


depgraph-install:
	go install github.com/kisielk/godepgraph@latest

depgraph-build:
	godepgraph ./cmd | dot -Tpng -o godepgraph.png

depgraph:
	make depgraph-install
	make depgraph-build


# ---------------------------
# Запуск базы данных в Docker
# ---------------------------

compose-up:
	docker-compose up -d 

compose-down:
	docker-compose down

compose-stop:
	docker-compose stop 

compose-start:
	docker-compose start 

compose-ps:
	docker-compose ps 

# ---------------------------
# Запуск миграций через Goose
# ---------------------------

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres "postgres://postgres:qwe@localhost:5432/postgres?sslmode=disable" create rename_me sql

goose-up:
	goose -dir ./migrations postgres "postgres://postgres:qwe@localhost:5432/postgres?sslmode=disable" up

goose-down:
	goose -dir ./migrations postgres "postgres://postgres:qwe@localhost:5432/postgres?sslmode=disable" down

goose-status:
	goose -dir ./migrations postgres "postgres://postgres:qwe@localhost:5432/postgres?sslmode=disable" status


#fake data
goose-add-fake:
	goose -dir ./migrations postgres "postgres://postgres:qwe@localhost:5433/postgresFake?sslmode=disable" create rename_me sql

goose-up-fake:
	goose -dir ./migrations postgres "postgres://postgres:qwe@localhost:5433/postgresFake?sslmode=disable" up

goose-down-fake:
	goose -dir ./migrations postgres "postgres://postgres:qwe@localhost:5433/postgresFake?sslmode=disable" down

goose-status-fake:
	goose -dir ./migrations postgres "postgres://postgres:qwe@localhost:5433/postgresFake?sslmode=disable" status

bin-deps: .vendor-proto
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@latest

generate:
	mkdir -p ${OUT_PATH}
	protoc --proto_path api --proto_path vendor.protogen \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=${OUT_PATH} --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=${OUT_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ${OUT_PATH} --grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=${OUT_PATH} \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:${OUT_PATH}" \
		./api/cli/v1/cliserver.proto

.vendor-proto: .vendor-proto/google/protobuf .vendor-proto/google/api .vendor-proto/protoc-gen-openapiv2/options .vendor-proto/validate

.vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

.vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf &&\
		cd vendor.protogen/protobuf &&\
		git sparse-checkout set --no-cone src/google/protobuf &&\
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/protobuf/src/google/protobuf vendor.protogen/google
		rm -rf vendor.protogen/protobuf

.vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
 		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p  vendor.protogen/google
		mv vendor.protogen/googleapis/google/api vendor.protogen/google
		rm -rf vendor.protogen/googleapis

.vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
		cd vendor.protogen/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.protogen/validate
		mv vendor.protogen/tmp/validate vendor.protogen/
		rm -rf vendor.protogen/tmp
