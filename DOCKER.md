# Docker Development Guide

This guide covers the Docker development environment setup for the CS2 Log Generator project.

## Quick Start

### Prerequisites

- Docker (20.10+)
- Docker Compose (2.0+)
- Make (optional, for convenience commands)

### One-Command Setup

```bash
# Run the automated setup script
./scripts/setup-dev.sh
```

### Manual Setup

```bash
# 1. Copy environment file
cp .env.example .env

# 2. Build and start development environment
make dev

# Or without make:
docker-compose up -d --build
```

## Services

### Backend (Go/Gin)
- **Port**: 8080
- **URL**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Hot Reload**: Enabled with Air
- **Volume Mount**: `./backend:/app`

### Frontend (React/Vite)
- **Port**: 5173
- **URL**: http://localhost:5173
- **Hot Reload**: Enabled with Vite HMR
- **Volume Mount**: `./frontend:/app`

### Redis (Optional)
- **Port**: 6379
- **Start with**: `make redis` or `docker-compose --profile with-cache up -d`

## Development Workflow

### Daily Development

```bash
# Start development environment
make dev

# View logs (all services)
make logs

# View specific service logs
make logs-backend
make logs-frontend

# Stop services
make down
```

### Building and Testing

```bash
# Build images
make build

# Run tests
make test

# Build production images
make build-prod
```

### Debugging

```bash
# Open shell in containers
make shell-backend
make shell-frontend

# Check service status
make status

# Health checks
make health

# Monitor resource usage
make monitor
```

## Docker Architecture

### Multi-Stage Builds

Each service uses multi-stage Dockerfiles:

#### Backend Dockerfile Stages:
1. **dev** - Development with hot reload (Air)
2. **builder** - Production build stage
3. **production** - Minimal runtime image

#### Frontend Dockerfile Stages:
1. **dev** - Development with Vite dev server
2. **builder** - Production build stage  
3. **production** - Nginx serving static files

### Development vs Production

| Aspect | Development | Production |
|--------|-------------|------------|
| **Target Stage** | `dev` | `production` |
| **Volume Mounts** | Yes (hot reload) | No |
| **Image Size** | Larger (dev tools) | Minimal |
| **Build Time** | Faster (cached) | Slower (optimized) |
| **Security** | Relaxed | Hardened |

## Environment Variables

### Backend Variables
```bash
PORT=8080                    # Server port
ENV=development             # Environment mode
LOG_LEVEL=debug             # Logging level
CORS_ORIGIN=http://localhost:5173  # CORS configuration
GIN_MODE=debug              # Gin framework mode
```

### Frontend Variables
```bash
VITE_API_URL=http://localhost:8080   # Backend API URL
VITE_WS_URL=ws://localhost:8080      # WebSocket URL
VITE_ENV=development                 # Environment mode
```

### Docker Compose Variables
```bash
COMPOSE_PROJECT_NAME=cs2-log-generator  # Project name
COMPOSE_HTTP_TIMEOUT=120               # HTTP timeout
```

## Available Make Commands

```bash
make help           # Show all available commands
make dev            # Start development environment
make build          # Build development images
make build-prod     # Build production images
make down           # Stop all services
make logs           # View logs from all services
make clean          # Remove containers and volumes
make test           # Run tests in containers
make health         # Check service health
make restart        # Restart all services
```

## Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Check what's using the port
lsof -i :8080
lsof -i :5173

# Or change ports in docker-compose.yml
```

#### Permission Issues
```bash
# Fix volume permissions (if needed)
sudo chown -R $USER:$USER ./backend ./frontend
```

#### Hot Reload Not Working
```bash
# Ensure proper volume mounts
docker-compose down
docker-compose up -d --build

# Check file watching inside container
make shell-backend
# Inside container: ls -la /app
```

#### Build Cache Issues
```bash
# Force rebuild without cache
make build
# or
docker-compose build --no-cache
```

#### Container Won't Start
```bash
# Check logs for errors
make logs

# Check service health
docker-compose ps

# Inspect specific service
docker-compose logs backend
```

### Performance Optimization

#### Volume Performance on macOS/Windows
```yaml
# In docker-compose.yml, use delegated mounting
volumes:
  - ./backend:/app:delegated
  - ./frontend:/app:delegated
```

#### Build Speed
```bash
# Use BuildKit for faster builds
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1
```

## Production Deployment

### Production Build
```bash
# Build production images
make build-prod

# Start production environment
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Production Configuration
- Uses optimized, minimal images
- No volume mounts
- Nginx serves frontend static files
- Environment variables configured for production

## Health Checks

All services include health checks:

- **Backend**: `curl http://localhost:8080/health`
- **Frontend**: HTTP check on port 5173/80
- **Redis**: `redis-cli ping`

## Networking

Services communicate through a custom bridge network:
- **Network Name**: `cs2-log-generator-network`
- **Internal Communication**: Services can reach each other by service name
- **External Access**: Only exposed ports are accessible from host

## Volumes

### Development Volumes
- Source code mounted for hot reload
- Node modules/Go vendor excluded from host mount

### Production Volumes
- Redis data persistence: `cs2-log-generator-redis-data`

## Security Considerations

### Development
- Relaxed security for development convenience
- CORS enabled for localhost

### Production
- Non-root user in containers
- Minimal base images (Alpine Linux)
- Security headers in Nginx
- Environment variables for secrets

## Best Practices

1. **Always use .env files** - Never hardcode configuration
2. **Layer caching** - Copy dependencies first, then source code
3. **Multi-stage builds** - Keep final images minimal
4. **Health checks** - Ensure services are ready before accepting traffic
5. **Volume exclusions** - Don't mount node_modules or vendor directories
6. **Graceful shutdown** - Handle SIGTERM properly

## Development Tips

1. **Use make commands** for consistency across team
2. **Check logs frequently** with `make logs`
3. **Monitor resource usage** with `make monitor`
4. **Clean up regularly** with `make clean`
5. **Test in production mode** before deploying with `make build-prod`

## Contributing

When modifying Docker configuration:

1. Test both development and production builds
2. Update this documentation
3. Test on different platforms (Linux/macOS/Windows)
4. Update make commands if needed
5. Ensure backwards compatibility