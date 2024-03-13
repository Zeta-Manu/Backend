# Makefile for API

# parameters
PORT = 8080
DOCKER_IMAGE_NAME = zeta-api
DOCKER_CONTAINER_NAME = zeta-api-container

# Build Development Docker image
docker-build-dev:
	docker build -t $(DOCKER_IMAGE_NAME):development -f deployment/Dockerfile.dev .

# Run Development Docker container
docker-run-dev:
	docker run -p $(PORT):$(PORT) --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME):development

# Build and Run Development
dev:
	docker-compose build
	docker-compose up

dev-stop:
	docker-compose down

# Build Docker image
docker-build-prod:
	docker build -t $(DOCKER_IMAGE_NAME):production -f deployment/Dockerfile.prod .

# Run Docker container
docker-run-prod:
	docker run -p $(PORT):$(PORT) --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME):production

# Stop and remove Docker container
docker-stop:
	docker stop $(DOCKER_CONTAINER_NAME)
	docker rm $(DOCKER_CONTAINER_NAME)

tidy:
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest

swag:
	swag init -g cmd/app/main.go

DB_URL = "mysql://root:password@tcp(localhost:3307)/manu"
MIGRATIONS_PATH = db/migrations

migrate:
	@echo "Please specifiy 'up' or 'down' as a sub-target"
	@echo $(DB_URL)

migrate-install:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	@echo "Running migrations up..."
	migrate -database $(DB_URL) -path $(MIGRATIONS_PATH) up

migrate-down:
	@echo "Running migrations down..."
	migrate -database $(DB_URL) -path $(MIGRATIONS_PATH) down

.PHONY: migrate migrate-up migrate-down migrate-install
