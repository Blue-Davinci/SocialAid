.PHONY: help
help:
	@echo "make run/api - Run the API"

.PHONY: run/api
run/api:
	@echo "Running the API"
	@go run ./cmd/api
