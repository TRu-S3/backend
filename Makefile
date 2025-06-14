.PHONY: build run stop clean logs

# Build the Docker image
build:
	docker-compose build

# Run the application
run:
	docker-compose up -d

# Stop the application
stop:
	docker-compose down

# Clean up containers and images
clean:
	docker-compose down --rmi all --volumes --remove-orphans

# View logs
logs:
	docker-compose logs -f

# Rebuild and run
restart:
	docker-compose down
	docker-compose build
	docker-compose up -d

# Run in development mode (with logs visible)
dev:
	docker-compose up

# Execute shell in running container
shell:
	docker-compose exec app sh

# Build and run locally (without Docker)
local-run:
	go run main.go

# Test the application
test:
	curl http://localhost:8080
	curl http://localhost:8080/health
