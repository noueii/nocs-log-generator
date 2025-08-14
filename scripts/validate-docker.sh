#!/bin/bash

# CS2 Log Generator - Docker Configuration Validation Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Validation function
validate_config() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════╗"
    echo "║            Docker Configuration Validation              ║"
    echo "╚══════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    # Check if files exist
    log "Checking required Docker files..."
    
    files=(
        "docker-compose.yml"
        "docker-compose.prod.yml"
        "backend/Dockerfile"
        "frontend/Dockerfile"
        "backend/.dockerignore"
        "frontend/.dockerignore"
        ".env.example"
        "Makefile"
    )
    
    for file in "${files[@]}"; do
        if [[ -f "$file" ]]; then
            success "$file ✓"
        else
            error "$file ✗"
            exit 1
        fi
    done
    
    # Validate docker-compose.yml syntax
    log "Validating docker-compose.yml syntax..."
    if docker compose -f docker-compose.yml config >/dev/null 2>&1; then
        success "docker-compose.yml syntax is valid"
    else
        error "docker-compose.yml syntax validation failed"
        docker compose -f docker-compose.yml config 2>&1 | head -5
        exit 1
    fi
    
    # Validate production compose file
    log "Validating production docker-compose configuration..."
    if docker compose -f docker-compose.yml -f docker-compose.prod.yml config >/dev/null 2>&1; then
        success "Production compose configuration is valid"
    else
        error "Production compose configuration validation failed"
        docker compose -f docker-compose.yml -f docker-compose.prod.yml config 2>&1 | head -5
        exit 1
    fi
    
    # Check Dockerfile syntax (basic)
    log "Checking Dockerfile syntax..."
    
    # Backend Dockerfile
    if grep -q "FROM golang:" backend/Dockerfile && grep -q "EXPOSE 8080" backend/Dockerfile; then
        success "Backend Dockerfile looks valid"
    else
        error "Backend Dockerfile appears invalid"
        exit 1
    fi
    
    # Frontend Dockerfile
    if grep -q "FROM node:" frontend/Dockerfile && grep -q "EXPOSE" frontend/Dockerfile; then
        success "Frontend Dockerfile looks valid"
    else
        error "Frontend Dockerfile appears invalid"
        exit 1
    fi
    
    # Check .env.example
    log "Validating environment variables template..."
    if grep -q "PORT=8080" .env.example && grep -q "VITE_API_URL" .env.example; then
        success "Environment variables template looks complete"
    else
        warning "Environment variables template may be incomplete"
    fi
    
    # Check Makefile
    log "Validating Makefile..."
    if grep -q "dev:" Makefile && grep -q "build:" Makefile && grep -q "help:" Makefile; then
        success "Makefile contains required targets"
    else
        error "Makefile is missing required targets"
        exit 1
    fi
    
    # Summary
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════════════╗"
    echo "║                 Validation Complete!                    ║"
    echo "╚══════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    echo "✅ All Docker configuration files are present and valid"
    echo ""
    echo "🚀 Ready to start development:"
    echo "   ./scripts/setup-dev.sh   # Automated setup"
    echo "   make dev                 # Manual start"
    echo ""
    echo "🏗️  Build commands:"
    echo "   make build               # Development images"
    echo "   make build-prod          # Production images"
    echo ""
    echo "📋 All available commands:"
    echo "   make help                # Show all commands"
}

# Check if docker compose is available
if ! command -v docker >/dev/null 2>&1; then
    error "Docker is not available. Please install Docker first."
    exit 1
fi

# Check compose command
if ! docker compose version >/dev/null 2>&1 && ! docker-compose --version >/dev/null 2>&1; then
    error "Docker Compose is not available. Skipping syntax validation."
    warning "Will only perform basic file existence checks."
    echo ""
fi

# Run validation
validate_config