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


depgraph-install:
	go install github.com/kisielk/godepgraph@latest

depgraph-build:
	godepgraph ./cmd | dot -Tpng -o godepgraph.png

depgraph:
	make depgraph-install
	make depgraph-build







