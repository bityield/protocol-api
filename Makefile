IMAGE=svc-bityield-api:latest

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