check: mock test

mock:
	@echo "Generate Mock..."
	@sh mockgen.sh

test:
	@echo "Run Unit Test..."
	@sh coverage.sh

build:
	@go build cmd/api/main.go

gorun:
	@go run cmd/api/main.go

gorun-mq:
	@go run cmd/mq/main.go

docker-start:
	@docker-compose up -d

docker-stop:
	@docker-compose down

run: gorun

run-mq: docker-start gorun-mq

coverage:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover-backend.html
	go tool cover -func cover.out