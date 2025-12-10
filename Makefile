.PHONY: build lint run test docker-start docker-stop docker-clean

build:
	go build ./cmd/main/main.go

lint:
	golangci-lint fmt
	golangci-lint run

run:
	go run ./cmd/main/main.go

test:
	go test ./...

watch-test:
	watchexec -e go "go test ./..."

docker-start:
	docker compose up --build

docker-stop:
	docker compose down

docker-clean:
	docker compose down --volumes
