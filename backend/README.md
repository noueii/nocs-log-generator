# CS2 Log Generator Backend

Go backend service for generating CS2 match logs and parsing demo files.

## Project Structure

```
backend/
├── cmd/server/          # Main application entry point
│   └── main.go         # HTTP server with Gin framework
├── pkg/
│   ├── api/            # HTTP handlers and routes
│   ├── generator/      # Match log generation logic
│   ├── parser/         # Demo file parsing logic
│   ├── models/         # Data structures and types
│   └── utils/          # Shared utilities
└── go.mod              # Go module definition
```

## Quick Start

### Prerequisites
- Go 1.21 or later
- Git

### Development

1. **Navigate to backend directory**:
   ```bash
   cd backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the server**:
   ```bash
   go run cmd/server/main.go
   ```

4. **Test the API**:
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # API status
   curl http://localhost:8080/api/v1/status
   ```

### Available Endpoints

- `GET /health` - Health check endpoint
- `GET /api/v1/status` - API status information

## Environment Variables

- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Gin mode: debug/release (default: debug)

## Development Notes

### Current Status
- ✅ Basic server setup with Gin framework
- ✅ Health check endpoints
- ✅ CORS middleware for frontend development
- ✅ Project structure established
- 🔄 Match generation logic (upcoming)
- 🔄 Demo parsing integration (upcoming)

### API Design Principles
- RESTful endpoints
- JSON request/response format
- Proper HTTP status codes
- Error handling with descriptive messages
- CORS enabled for development

### Code Structure
- **Package naming**: lowercase, single word
- **Files**: snake_case.go
- **Exported functions**: PascalCase
- **Private functions**: camelCase
- **Constants**: PascalCase

## Future Implementation

### Phase 2: Core Features
- Match generation algorithm
- Event simulation logic
- Log formatting system
- Demo parser integration with demoinfocs-golang

### Phase 3: Enhancement
- WebSocket streaming for real-time logs
- Advanced configuration options
- Performance optimization
- Comprehensive error handling