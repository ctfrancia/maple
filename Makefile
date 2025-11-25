# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	
.PHONY: run-dev
run-dev: ## Run the application in dev mode (default)
	go run ./...

run-db-dev: ## Run a local instance of postgres for development
	docker build -f Dockerfile.postgres -t my-postgres .
	docker run -d -p 5432:5432 --name postgres-dev my-postgres

.PHONY: run-docker
run-docker: ## Run the application in docker for development
	@docker build -t maple .
	@docker run -p 8080:8080 maple


# =========== DEPENDENCIES ===========
GOLANG          := golang:1.25
ALPINE          := alpine:3.22
POSTGRES        := postgres:18.0
GRAFANA         := grafana/grafana:12.2.0
PROMETHEUS      := prom/prometheus:v3.7.0
TEMPO           := grafana/tempo:2.9.0
LOKI            := grafana/loki:3.5.0
PROMTAIL        := grafana/promtail:3.5.0

NAMESPACE       := chess-system
CHESS_APP       := chess
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/maple
VERSION         := 0.0.1
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)

# =========== BUILD CONTAINERS ===========

build: maple metrics auth ## Build all containers

maple: ## Build the maple container

metrics: ## Build the metrics container

auth: ## Build the auth container
