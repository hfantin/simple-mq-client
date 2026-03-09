# load environment variables
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

all: update run

update: 
	@go get -u ./...
	@go mod tidy

run:
	@echo "running..."
	@go run main.go
