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
