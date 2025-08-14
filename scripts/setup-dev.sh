#!/bin/bash

# CS2 Log Generator - Development Setup Script
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

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Main setup function
main() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════╗"
    echo "║                CS2 Log Generator Setup                   ║"
    echo "║                Development Environment                   ║"
    echo "╚══════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    # Check prerequisites
    log "Checking prerequisites..."
    
    if ! command_exists docker; then
        error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    # Check for docker compose (modern syntax) or docker-compose (legacy)
    if ! command_exists "docker compose" && ! command_exists docker-compose; then
        error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    if ! command_exists make; then
        warning "Make is not installed. Some convenience commands won't work."
    fi
    
    success "Prerequisites check passed"
    
    # Create environment file
    log "Setting up environment configuration..."
    if [ ! -f .env ]; then
        cp .env.example .env
        success "Created .env from .env.example"
    else
        warning ".env already exists, skipping"
    fi
    
    # Detect docker compose command
    if command_exists "docker compose"; then
        DOCKER_COMPOSE="docker compose"
    else
        DOCKER_COMPOSE="docker-compose"
    fi
    
    # Build Docker images
    log "Building Docker images (this may take a few minutes)..."
    $DOCKER_COMPOSE build
    success "Docker images built successfully"
    
    # Start services
    log "Starting development environment..."
    $DOCKER_COMPOSE up -d
    success "Development environment started"
    
    # Wait for services to be healthy
    log "Waiting for services to be ready..."
    sleep 10
    
    # Check service health
    log "Checking service health..."
    
    # Backend health check
    if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
        success "Backend is healthy (http://localhost:8080)"
    else
        warning "Backend health check failed - it may still be starting up"
    fi
    
    # Frontend health check
    if curl -sf http://localhost:5173 >/dev/null 2>&1; then
        success "Frontend is running (http://localhost:5173)"
    else
        warning "Frontend health check failed - it may still be starting up"
    fi
    
    # Display information
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════════════╗"
    echo "║                     Setup Complete!                     ║"
    echo "╚══════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    echo "🚀 Development environment is ready!"
    echo ""
    echo "📍 URLs:"
    echo "   Frontend: http://localhost:5173"
    echo "   Backend:  http://localhost:8080"
    echo "   API:      http://localhost:8080/api/v1"
    echo ""
    echo "🛠️  Useful commands:"
    echo "   make dev      - Start development environment"
    echo "   make down     - Stop all services"
    echo "   make logs     - View container logs"
    echo "   make help     - Show all available commands"
    echo ""
    echo "📝 Next steps:"
    echo "   1. Open your browser to http://localhost:5173"
    echo "   2. Check the backend API at http://localhost:8080/health"
    echo "   3. View logs with: make logs"
    echo ""
    echo "For more information, see README.md"
}

# Handle script termination
cleanup() {
    echo ""
    warning "Setup interrupted. You can run this script again anytime."
    exit 1
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Run main function
main "$@"