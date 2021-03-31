DC=docker-compose
IG=svc-protocol-api:latest

GO_SRC_DIRS := $(shell \
	find . -name "*.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)

GO_TEST_DIRS := $(shell \
	find . -name "*_test.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)

.PHONY: all build compose-down compose-up deps run test watch

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

deploy:
	@echo "Deploying new version..."
	git fetch
	git merge origin/master
	sudo service protocol-api restart

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

lint:
	@golint ./...

run:
	@PORT=8000 ./bin/api

test:
	@go test -v ./...

watch:
	@air