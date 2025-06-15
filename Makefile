.PHONY: build run stop clean logs help

# Help
help:
	@echo "TRu-S3 Backend - Available commands:"
	@echo ""
	@echo "Docker commands:"
	@echo "  build                  - Build the Docker image"
	@echo "  run                    - Run the application with Docker"
	@echo "  stop                   - Stop the application"
	@echo "  clean                  - Clean up containers and images"
	@echo "  restart                - Rebuild and run"
	@echo "  dev                    - Run in development mode (with logs)"
	@echo "  logs                   - View logs"
	@echo "  shell                  - Execute shell in running container"
	@echo ""
	@echo "Local development:"
	@echo "  local-run              - Run locally without Docker"
	@echo "  local-cloudsql         - Run locally with Cloud SQL"
	@echo ""
	@echo "Cloud SQL commands:"
	@echo "  cloudsql-info          - Show Cloud SQL instance info"
	@echo "  cloudsql-setup         - Setup Cloud SQL Auth Proxy"
	@echo "  cloudsql-setup-complete - Complete automated setup"
	@echo "  cloudsql-troubleshoot  - Troubleshoot Cloud SQL issues"
	@echo "  cloudsql-local         - Start Cloud SQL Auth Proxy"
	@echo "  cloudsql-compose       - Start with Docker Compose + Cloud SQL"
	@echo ""
	@echo "Cloud Run commands:"
	@echo "  cloudrun-deploy        - Deploy to Cloud Run"
	@echo "  cloudrun-build         - Build image for Cloud Run"
	@echo ""
	@echo "Testing:"
	@echo "  test                   - Basic API test"
	@echo "  test-all               - Comprehensive API test"

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

# Cloud SQL関連タスク
cloudsql-info:
	./check-cloudsql.sh

cloudsql-setup:
	./setup-cloudsql-proxy.sh

cloudsql-setup-complete:
	./setup-cloud-sql-complete.sh

cloudsql-troubleshoot:
	./troubleshoot-cloudsql.sh

cloudsql-local:
	source .env.local && ./cloud-sql-proxy $$CLOUD_SQL_CONNECTION_NAME --port=5433

cloudsql-compose:
	docker-compose -f docker-compose.cloudsql.yml up

# Cloud Run関連タスク
cloudrun-deploy:
	./deploy-cloudrun.sh

cloudrun-build:
	gcloud builds submit --config=cloudbuild.yaml

# 環境別の起動
local-dev:
	go run main.go

local-cloudsql:
	source .env.local && go run main.go

# 包括的なテスト
test-all:
	curl -s http://localhost:8080/health | jq .
	curl -s http://localhost:8080/api/v1/files | jq .
	echo "Test file content" > test-upload.txt
	curl -X POST -F "file=@test-upload.txt" http://localhost:8080/api/v1/files
	curl -s http://localhost:8080/api/v1/files | jq .
	rm -f test-upload.txt
