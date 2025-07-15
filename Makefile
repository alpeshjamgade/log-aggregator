APP_BINARY=logAggregatorApp
DATABASE_NAME=log_aggregator
CLICKHOUSE_URL="clickhouse://default:@localhost:9000"
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

# ClickHouse Migration Commands (using clickhouse-cli or ch-migrate)
migration_create:
	@echo "Creating migration not supported via CLI for ClickHouse."
	@echo "Please create a new SQL file manually in ${MIGRATION_PATH}"

migration_up:
	@echo "Running migration..."
	clickhouse-client --database=$(DATABASE_NAME) --query="$(shell cat ${MIGRATION_PATH}/*.up.sql | tr '\n' ' ')"
	@echo "Done!"

migration_down:
	@echo "Rolling back migration..."
	clickhouse-client --database=$(DATABASE_NAME) --query="$(shell cat ${MIGRATION_PATH}/*.down.sql | tr '\n' ' ')"
	@echo "Done!"

migration_fix:
	@echo "ClickHouse does not track migration versions like other DBs."
	@echo "Use a manual table to manage migration state if needed."
