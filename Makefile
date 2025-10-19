.PHONY: build run docker-build docker-run docker-compose-up docker-compose-down k8s-deploy

# Local development
build:
	go build -o bin/scraper-api ./cmd/api
	go build -o bin/goscraper-cli ./cmd/cli

run:
	go run ./cmd/api/main.go

run-cli:
	go run ./cmd/cli/main.go

test:
	go test ./...

# Configuration management
init-config:
	go run ./cmd/cli init

setup:
	go run ./cmd/cli setup

validate-config:
	go run ./cmd/cli validate

show-config:
	go run ./cmd/cli config

# Docker commands
docker-build:
	docker build -t scraper-api:latest .

docker-run:
	docker run -p 8080:8080 scraper-api:latest

# Docker Compose
docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

docker-compose-logs:
	docker-compose logs -f scraper-api

# Kubernetes
k8s-deploy:
	kubectl apply -f k8s/

k8s-delete:
	kubectl delete -f k8s/

k8s-logs:
	kubectl logs -f deployment/scraper-api

# Health check
health:
	curl http://localhost:8080/health

# Test scraping
test-scrape:
	curl -X POST http://localhost:8080/api/scrape \
		-H "Content-Type: application/json" \
		-d '{"url": "https://example.com"}'

# Production deployment
deploy-prod:
	docker build -t your-registry/scraper-api:$(shell git rev-parse --short HEAD) .
	docker push your-registry/scraper-api:$(shell git rev-parse --short HEAD)
	kubectl set image deployment/scraper-api scraper-api=your-registry/scraper-api:$(shell git rev-parse --short HEAD)