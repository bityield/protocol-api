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
	@PORT=3000 ./bin/api

watch:
	@air