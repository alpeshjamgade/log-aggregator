APP_BINARY=logAggregatorApp

# Build binary
build-binary:
	@echo "Building app..."
	CGO_ENABLED=0 go build -o build/$(APP_BINARY) cmd/main.go
	@echo "Done!"

# Build and start the binary
run-build: build
	@echo "Starting app..."
	build/$(APP_BINARY)
	@echo "Done!"

# Run project
run:
	@echo "Starting app..."
	go run cmd/main.go
	@echo "Done!"

# Dockerize
dockerize:
	@echo "Building app..."
	docker build -t registry.tradelab.in/log-aggregator:$(TAG) .
	@echo "Done!"
