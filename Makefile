.PHONY: run-dev run-docker

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	
run-dev: ## Run the application in dev mode (default)
	go run ./...

run-docker: ## Run the application in docker for development
	@docker build -t maple .
	@docker run -p 8080:8080 maple
