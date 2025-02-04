.PHONY: help
help:
	@echo "make run/api - Run the API"
	@echo "make build/api - Build the API"

.PHONY: run/api
run/api:
	@echo "Running the API"
	@go run ./cmd/api

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags '-s' -o ./bin/api.exe ./cmd/api