APP=cli.exe
APP_EXECUTABLE="./cmd/$(APP)"
SHELL := /bin/bash 


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


build:
	make lint
	go build -o $(APP_EXECUTABLE) ./cmd/main.go
	@echo "Build passed"

run: 
	make build
	$(APP_EXECUTABLE)

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