IMAGE=svc-bityield-api:latest
DC=docker-compose

build:
	@echo "building [api]..."
	@go build -o bin/api main.go

deps:
	@go get -u ./...
	@go mod download
	@go mod tidy
	@go mod vendor

docker-build:
	@docker build --squash -t $(IMAGE) -f Dockerfile .

run:
	@PORT=8000 ./bin/api

up:
	@docker-compose up

watch:
	@PORT=8000 air

docker:
	$(DC) -f $(DC).yml up --build

down:
	-docker stop protocol-api_protocol-api_1
	-docker rm protocol-api_protocol-api_1
	-docker stop protocol-api_postgres_1
	-docker rm protocol-api_postgres_1
	-docker stop protocol-api_redis_1
	-docker rm protocol-api_redis_1