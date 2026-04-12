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

NAMESPACE       := maple-system
MAPLE_APP       := maple
AUTH_APP        := auth
SCRAPER_APP     := scraper
BASE_IMAGE_NAME := localhost/maple
VERSION         := 0.0.1
MAPLE_IMAGE     := $(BASE_IMAGE_NAME)/$(MAPLE_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)
SCRAPER_IMAGE   := $(BASE_IMAGE_NAME)/$(SCRAPER_APP):$(VERSION)
# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"

# ==============================================================================
# Detect operating system and set the appropriate open command

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	OPEN_CMD := open
else
	OPEN_CMD := xdg-open
endif

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-brew:
	brew update
	brew list pgcli || brew install pgcli
	brew list watch || brew install watch

dev-docker:
	docker pull docker.io/$(GOLANG) & \
	docker pull docker.io/$(ALPINE) & \
	docker pull docker.io/$(POSTGRES) & \
	docker pull docker.io/$(GRAFANA) & \
	docker pull docker.io/$(PROMETHEUS) & \
	docker pull docker.io/$(TEMPO) & \
	docker pull docker.io/$(LOKI) & \
	docker pull docker.io/$(PROMTAIL) & \
	wait;

# =========== BUILD CONTAINERS ===========

build: maple metrics auth scraper ## Build all containers

maple: ## Build the maple container
	docker build \
		-f zoltan/docker/dockerfile.maple \
		-t $(MAPLE_IMAGE) \
		-t $(BASE_IMAGE_NAME)/backend:dev \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

metrics: ## Build the metrics container
	docker build \
		-f zoltan/docker/dockerfile.metrics \
		-t $(METRICS_IMAGE) \
		-t $(BASE_IMAGE_NAME)/metrics:dev \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

auth: ## Build the auth container
	docker build \
		--no-cache \
		-f zoltan/docker/dockerfile.auth \
		-t $(AUTH_IMAGE) \
		-t $(BASE_IMAGE_NAME)/auth:dev \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

scraper: ## Build the scraper container
	docker build \
		-f zoltan/docker/dockerfile.scrape \
		-t $(SCRAPER_IMAGE) \
		-t $(BASE_IMAGE_NAME)/scraper:dev \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.


# ==============================================================================
# Metrics and Tracing

metrics-view-sc: ## View the metrics in a browser
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

metrics-view: ## View the metrics in a browser
	expvarmon -ports="localhost:4020" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

grafana: ## view the grafana dashboard
	$(OPEN_CMD) http://localhost:3100/

statsviz: ## view the statsviz dashboard
	$(OPEN_CMD) http://localhost:3010/debug/statsviz/


# ==============================================================================
# Audit

audit: ## Run the audit service
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...
	govulncheck ./...

# ==============================================================================
# RUN
#
.PHONY: dev
dev: build ## Run the application in dev mode locally
	docker compose -f zoltan/compose/docker-compose.yml -f zoltan/compose/docker-compose.dev.yml --env-file zoltan/compose/.env.dev up --force-recreate

.PHONY: dev-clean
dev-clean: down-dev ## Clean and restart dev environment
	make dev

.PHONY: down-dev
down-dev: ## stop the application in dev mode locally 
	docker compose -f zoltan/compose/docker-compose.yml -f zoltan/compose/docker-compose.dev.yml down -v
