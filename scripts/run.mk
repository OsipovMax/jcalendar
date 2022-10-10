all: run

run: ## Run application
run:
	go run ./cmd/${APP_NAME}

.PHONY: all run
