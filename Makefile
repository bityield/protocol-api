DC=docker-compose
IG=svc-bityield-api:latest

.PHONY: all build compose-down compose-up deps run watch

all: build

b: build
cd: compose-down
cu: compose-up
d: deps
r: run
w: watch

build:
	@echo "building [api]..."
	@go build -o bin/api main.go

compose-down:
	-docker stop protocol-api_protocol-api_1
	-docker rm protocol-api_protocol-api_1
	-docker stop protocol-api_postgres_1
	-docker rm protocol-api_postgres_1
	-docker stop protocol-api_redis_1
	-docker rm protocol-api_redis_1

compose-up:
	@$(DC) -f $(DC).yml up --build

deps:
	@go get -u ./...
	@go mod download
	@go mod tidy
	@go mod vendor

run:
	@PORT=8000 ./bin/api

watch:
	@air