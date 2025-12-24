.PHONY: run build test migrate-up migrate-down

run:
	go run cmd/api/main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

migrate-up:
	migrate -path internal/infrastructure/database/postgres/migrations -database "${DATABASE_URL}" up

migrate-down:
	migrate -path internal/infrastructure/database/postgres/migrations -database "${DATABASE_URL}" down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

lint:
	golangci-lint run

mod:
	go mod tidy && go mod download