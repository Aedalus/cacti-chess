ONESHELL:
.DEFAULT_GOAL := help
.SHELLFLAGS = -ec

export GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
export GIT_HASH := $(shell git rev-parse HEAD)

export NAMESPACE := $(shell basename $$PWD)
export ARTIFACT_ID := $(NAMESPACE).$(GIT_HASH)

export env ?= dev

# --- Formatting --------------------------------------------------------------
export RED ?= '\033[0;31m'
export GREEN ?= '\033[0;32m'
export NO_COLOR ?= '\033[0m'

.PHONY : help
help:
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run Unit Tests
	go test -coverprofile cover.out ./...

coverage: ## Run Unit Tests With Coverage
	go tool cover -html=cover.out

build: ## Builds the binary
	rm -rf ./bin
	mkdir -p ./bin
	cd cmd
	go build -o ./bin/cacti-chess ./cmd
	go build -o ./bin/cacti-chess-uci ./uci
	chmod +x ./bin/cacti-chess
	chmod +x ./bin/cacti-chess-uci

build-lichess: ## Builds the lichess docker image
	docker build -f lichess.dockerfile . -t lichess-bot

run-lichess: ## Run the lickess bot
	go run ./lichess-bot lichess.toml