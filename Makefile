# CS2 Log Generator - Development Makefile
.PHONY: help dev build build-prod up down logs clean test backend frontend health status restart shell-backend shell-frontend

# Default target
.DEFAULT_GOAL := help

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

help: ## Show this help message
	@echo "$(CYAN)CS2 Log Generator - Available Commands:$(RESET)"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2}'
	@echo
	@echo "$(YELLOW)Environment:$(RESET)"
	@echo "  Development: make dev"
	@echo "  Production:  make build-prod && make up-prod"
	@echo
	@echo "$(YELLOW)URLs:$(RESET)"
	@echo "  Frontend: http://localhost:5173"
	@echo "  Backend:  http://localhost:8080"
	@echo "  API:      http://localhost:8080/api/v1"

dev: ## Start development environment with hot-reload
	@echo "$(CYAN)Starting development environment...$(RESET)"
	docker compose up -d
	@echo "$(GREEN)Development environment started!$(RESET)"
	@echo "Frontend: http://localhost:5173"
	@echo "Backend:  http://localhost:8080"
	@echo "Use 'make logs' to view output"

up: dev ## Alias for dev

build: ## Build development Docker images
	@echo "$(CYAN)Building development images...$(RESET)"
	docker compose build --no-cache

build-prod: ## Build production Docker images
	@echo "$(CYAN)Building production images...$(RESET)"
	docker compose -f docker-compose.yml -f docker-compose.prod.yml build --no-cache

up-prod: ## Start production environment
	@echo "$(CYAN)Starting production environment...$(RESET)"
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
	@echo "$(GREEN)Production environment started!$(RESET)"
	@echo "Frontend: http://localhost"
	@echo "Backend:  http://localhost:8080"

down: ## Stop and remove all containers
	@echo "$(CYAN)Stopping containers...$(RESET)"
	docker compose down
	@echo "$(GREEN)Containers stopped$(RESET)"

down-prod: ## Stop production environment
	docker compose -f docker-compose.yml -f docker-compose.prod.yml down

logs: ## View logs from all containers
	docker compose logs -f

logs-backend: ## View backend logs only
	docker compose logs -f backend

logs-frontend: ## View frontend logs only
	docker compose logs -f frontend

clean: ## Remove ONLY this project's containers and volumes
	@echo "$(YELLOW)Cleaning up CS2 Log Generator Docker resources...$(RESET)"
	docker compose down -v --remove-orphans
	@echo "$(GREEN)CS2 Log Generator cleanup complete$(RESET)"

clean-images: ## Remove ONLY this project's Docker images
	@echo "$(YELLOW)Removing CS2 Log Generator images...$(RESET)"
	docker compose down
	docker rmi nocs-log-generator-backend nocs-log-generator-frontend 2>/dev/null || true
	@echo "$(GREEN)CS2 Log Generator images removed$(RESET)"

clean-all: ## Remove ONLY this project's containers, volumes, and images
	@echo "$(YELLOW)Removing all CS2 Log Generator resources...$(RESET)"
	docker compose down -v --remove-orphans
	docker rmi nocs-log-generator-backend nocs-log-generator-frontend 2>/dev/null || true
	@echo "$(GREEN)CS2 Log Generator full cleanup complete$(RESET)"

# DANGEROUS: System-wide cleanup commands (hidden from help)
.PHONY: system-prune-careful
system-prune-careful: ## WARNING: Removes ALL stopped containers system-wide
	@echo "$(RED)⚠️  WARNING: This will remove ALL stopped containers from ALL projects!$(RESET)"
	@echo "$(RED)This is NOT project-specific and will affect other projects!$(RESET)"
	@read -p "Are you ABSOLUTELY sure? Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		docker system prune -f; \
	else \
		echo "Cancelled - no changes made"; \
	fi

test: ## Run tests in containers
	@echo "$(CYAN)Running backend tests...$(RESET)"
	docker compose exec backend go test ./... || true
	@echo "$(CYAN)Running frontend tests...$(RESET)"
	docker compose exec frontend npm test || true

test-backend: ## Run backend tests only
	docker compose exec backend go test ./...

test-frontend: ## Run frontend tests only
	docker compose exec frontend npm test

backend: ## Start only backend service
	docker compose up -d backend
	@echo "$(GREEN)Backend started: http://localhost:8080$(RESET)"

frontend: ## Start only frontend service  
	docker compose up -d frontend
	@echo "$(GREEN)Frontend started: http://localhost:5173$(RESET)"

health: ## Check health status of all services
	@echo "$(CYAN)Checking service health...$(RESET)"
	@echo "Backend Health:"
	@curl -s http://localhost:8080/health | jq . 2>/dev/null || curl -s http://localhost:8080/health || echo "Backend not responding"
	@echo
	@echo "Frontend Health:"
	@curl -s http://localhost:5173/health 2>/dev/null || echo "Frontend responding" || echo "Frontend not responding"

status: ## Show container status
	docker compose ps

restart: ## Restart all services
	@echo "$(CYAN)Restarting services...$(RESET)"
	docker compose restart
	@echo "$(GREEN)Services restarted$(RESET)"

restart-backend: ## Restart backend service
	docker compose restart backend

restart-frontend: ## Restart frontend service
	docker compose restart frontend

shell-backend: ## Open shell in backend container
	docker compose exec backend sh

shell-frontend: ## Open shell in frontend container
	docker compose exec frontend sh

# Redis commands (when using --profile with-cache)
redis: ## Start with Redis cache
	docker compose --profile with-cache up -d
	@echo "$(GREEN)Services started with Redis cache$(RESET)"

redis-cli: ## Connect to Redis CLI
	docker compose exec redis redis-cli

# Quick development workflows
quick-restart: down dev ## Quick restart of development environment

rebuild: down build dev ## Rebuild and restart development environment

# Docker maintenance
pull: ## Pull latest base images
	docker compose pull

images: ## List project images
	docker images | grep cs2-log-generator

# Development helpers
fmt-backend: ## Format Go code
	docker compose exec backend go fmt ./...

lint-backend: ## Lint Go code (requires golangci-lint in container)
	docker compose exec backend golangci-lint run ./... || echo "golangci-lint not available"

lint-frontend: ## Lint frontend code
	docker compose exec frontend npm run lint

# Environment file creation
env: ## Create .env files from examples
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(GREEN)Created .env from .env.example$(RESET)"; \
	else \
		echo "$(YELLOW).env already exists$(RESET)"; \
	fi

# Monitoring and debugging
monitor: ## Monitor container resource usage
	docker stats cs2-log-generator-backend cs2-log-generator-frontend

inspect-backend: ## Inspect backend container
	docker compose exec backend env

inspect-frontend: ## Inspect frontend container  
	docker compose exec frontend env