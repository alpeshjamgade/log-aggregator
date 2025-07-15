APP_BINARY=logAggregatorApp
DATABASE_URL="postgresql://postgres:postgres@localhost:5432/log_aggregator?sslmode=disable"
MIGRATION_PATH="./private/migrations/"

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

migration_create:
	migrate create -ext sql -dir ${MIGRATION_PATH} -seq ${name}

migration_up:
	migrate -path ${MIGRATION_PATH} -database ${DATABASE_URL} -verbose up ${NUMBER}

migration_down:
	migrate -path ${MIGRATION_PATH} -database ${DATABASE_URL} -verbose down ${NUMBER}

migration_fix:
	migrate -path ${MIGRATION_PATH} -database ${DATABASE_URL} force ${VERSION}
