up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

test:
	go test ./...