COMPOSE_FILE=./deployments/docker-compose.yml
DOCKER_FILE_PATH=./deployments/docker
PROJECT_NAME=GO-GCP-AUTH
SERVICES := auth server
TAGS := latest
BUILD_CONTEXT ?= .
DOCKER_ACCOUNT=16181181418

.PHONY: up-local up-debug up-production down restart build-images push-images logs

# Start all services in the background
up-local:
	@echo "Starting services with Docker Compose..."
	docker compose -f $(COMPOSE_FILE) --env-file ./env/.env.local up $(FLAGS)

up-debug:
	@echo "Starting services with Docker Compose..."
	docker compose -f $(COMPOSE_FILE) --env-file ./env/.env.debug up $(FLAGS)

up-production:
	@echo "Starting services with Docker Compose..."
	docker compose -f $(COMPOSE_FILE) --env-file ./env/.env.production up $(FLAGS)


# Stop all running services
down:
	@echo "Stopping services with Docker Compose..."
	docker compose -f $(COMPOSE_FILE) down

# Restart all services
restart: down up

build-images:
	@for service in $(SERVICES); do \
		for tag in $(TAGS); do \
			docker build \
				-t $(DOCKER_ACCOUNT)/go-gcp-auth_$$service:$$tag \
				-f $(DOCKER_FILE_PATH)/$$service.Dockerfile \
				$(BUILD_CONTEXT); \
		done; \
	done

push-images:
	@for service in $(SERVICES); do \
		for tag in $(TAGS); do \
			docker push $(DOCKER_ACCOUNT)/go-gcp-auth_$$service:$$tag; \
		done; \
	done


.PHONY: list
list:
	@echo "Available services:"
	@for service in $(SERVICES); do echo "  - $$service"; done
	@echo "\nAvailable tags:"
	@for tag in $(TAGS); do echo "  - $$tag"; done

# Show logs from all running services
logs:
	docker compose -f $(COMPOSE_FILE) logs -f
