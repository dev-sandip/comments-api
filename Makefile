APP_NAME := comments-api
GO := go

.PHONY: all build run test clean fmt vet goose

all: build

build:
	@$(GO) build -o bin/$(APP_NAME) ./cmd

run:
	@$(GO) run ./cmd

test:
	@$(GO) test ./...

fmt:
	@$(GO) fmt ./...

vet:
	@$(GO) vet ./...

clean:
	@rm -rf bin

# Goose migration
goose:
	@goose -s create $(filter-out $@,$(MAKECMDGOALS)) sql

# Allow migration name as an additional argument
%:
	@: