# Go Lambda API

A minimal Golang HTTP API built with the standard library, featuring API key authentication and designed for serverless deployment on AWS Lambda.

## Features

- ✅ Pure Go implementation using only `net/http` (no external frameworks)
- ✅ API key authentication via `X-API-Key` header
- ✅ Health check endpoint at `/health`
- ✅ Graceful shutdown support
- ✅ Request logging middleware
- ✅ Comprehensive unit tests
- ✅ AWS Lambda deployment ready with Serverless Framework
- ✅ Dual deployment: Run locally as HTTP server or deploy to AWS Lambda

## Project Structure

```
/
├── cmd/
│   ├── api/
│   │   └── main.go           # Local server entry point
│   └── lambda/
│       └── main.go           # AWS Lambda entry point
├── internal/
│   ├── handlers/
│   │   ├── health.go         # Health check handler
│   │   └── handlers_test.go  # Handler tests
│   ├── middleware/
│   │   ├── auth.go           # Authentication middleware
│   │   └── auth_test.go      # Middleware tests
│   └── server/
│       ├── server.go         # Server setup and configuration
│       └── server_test.go    # Server tests
├── serverless.yml            # Serverless Framework configuration
├── Taskfile.yml              # Task runner configuration
├── .air.toml                 # Hot reload configuration
├── Dockerfile                # Docker configuration
├── .gitignore                # Git ignore file
├── go.mod                    # Go module file
└── README.md                 # This file
```

## Local Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- [Task](https://taskfile.dev) - Task runner (recommended)
- (Optional) Serverless Framework for deployment
- (Optional) AWS CLI configured with appropriate credentials
- (Optional) Docker for containerized deployment

### Installation

1. Clone the repository:
```bash
git clone https://github.com/nicobistolfi/go-lambda-api.git
cd go-lambda-api
```

2. Install dependencies:
```bash
go mod download
# or using Task
task mod
```

3. Create a `.env` file from the example:
```bash
cp .env.example .env
# Edit .env and set your API_KEY
```

4. Install Task runner (if not already installed):
```bash
# macOS
brew install go-task/tap/go-task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

# Windows (using Scoop)
scoop install task
```

## Quick Start with Task

View all available tasks:
```bash
task --list
# or simply
task
```

Common operations:
```bash
# Run the server locally
task run

# Run tests
task test

# Run tests with coverage
task test-coverage

# Build the binary
task build

# Format code
task fmt

# Start development server with hot reload
task dev
```

## Environment Variables Configuration

The application uses the following environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `API_KEY` | API key for authentication | - | Yes |
| `PORT` | Port to run the server on | `8080` | No |

## Running the Server Locally

### Using Task (Recommended)

```bash
# Run with default dev API key
task run

# Run with custom API key
API_KEY="your-secret-api-key" task run

# Run on custom port
PORT=3000 task run
```

### Using Go directly

1. Set the required environment variables:
```bash
export API_KEY="your-secret-api-key"
export PORT="8080"  # Optional, defaults to 8080
```

2. Run the server:
```bash
go run cmd/api/main.go
```

### Using Docker

```bash
# Build and run in Docker
API_KEY="your-secret-api-key" task docker
```

The server will start on `http://localhost:8080` (or the port specified).

### Testing the API

Health check (no authentication required):
```bash
curl http://localhost:8080/health
# or using Task
curl http://localhost:8080/health | jq .
```

Expected response:
```json
{"status":"ok"}
```

With authentication (for future authenticated endpoints):
```bash
curl -H "X-API-Key: your-secret-api-key" http://localhost:8080/some-endpoint
# or using Task with custom API key
API_KEY="your-secret-api-key" task run
```

## Running Tests

### Using Task (Recommended)

```bash
# Run all tests
task test

# Run tests with coverage
task test-coverage

# Run linter
task lint

# Clean build artifacts
task clean
```

### Using Go directly

Run all tests with coverage:
```bash
go test -v -cover ./...
```

Run tests for a specific package:
```bash
go test -v ./internal/handlers
go test -v ./internal/middleware
go test -v ./internal/server
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Serverless Deployment

### Prerequisites

1. Install Serverless Framework:
```bash
npm install -g serverless
```

2. Install dependencies:
```bash
npm install
```

3. Configure AWS credentials:
```bash
aws configure
```

### Deployment Steps

#### Using Task (Recommended)

Make sure you have a `.env` file with your API_KEY set, or pass it explicitly:

```bash
# Deploy using .env file
task deploy

# Or deploy with explicit API_KEY
API_KEY="your-api-key" task deploy

# Deploy to specific stage
STAGE=production task deploy

# View logs
task logs
```

#### Using Serverless directly

1. Set your API key as an environment variable:
```bash
export API_KEY="your-production-api-key"
```

2. Deploy to AWS:
```bash
serverless deploy --stage production --region us-west-1
```

3. Deploy to a specific stage:
```bash
serverless deploy --stage dev
serverless deploy --stage staging
serverless deploy --stage production
```

### Viewing Logs

View function logs:
```bash
serverless logs -f api --tail
# or using Task
task logs
```

### Removing the Deployment

Remove the deployed service:
```bash
serverless remove --stage production
# or using Task
STAGE=production serverless remove
```

## API Documentation

### Endpoints

#### `GET /health`
Health check endpoint that returns the service status.

**Authentication**: Not required

**Response:**
- Status: `200 OK`
- Body: `{"status": "ok"}`

### Authentication

All endpoints (except `/health`) require API key authentication via the `X-API-Key` header.

**Example:**
```bash
curl -H "X-API-Key: your-api-key" https://your-api-url.com/endpoint
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid API key
  - `{"error": "Missing API key"}`
  - `{"error": "Invalid API key"}`
  - `{"error": "API key not configured"}`

## Development Guidelines

### Development Tools

Install development dependencies:
```bash
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

This installs:
- `air` - Hot reload for development
- `golangci-lint` - Linting tool

Start development server with hot reload:
```bash
task dev
```

Run code checks:
```bash
# Format code
task fmt

# Run linter (if installed)
task lint

# Run default task (format, test, build)
task default
```

Clean build artifacts:
```bash
task clean
```

### Adding New Endpoints

1. Create a new handler in `internal/handlers/`
2. Add authentication by wrapping with `middleware.AuthMiddleware()`
3. Register the route in `internal/server/server.go`
4. Write comprehensive tests

Example:
```go
// In internal/server/server.go
mux.HandleFunc("/api/users", middleware.AuthMiddleware(handlers.UsersHandler))
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions small and focused
- Write tests for all new functionality
- Use meaningful variable and function names

## Troubleshooting

### Common Issues

1. **Server fails to start**
   - Check if the port is already in use
   - Ensure all environment variables are set correctly

2. **Authentication failures**
   - Verify the `API_KEY` environment variable is set
   - Check that the `X-API-Key` header matches exactly

3. **Deployment issues**
   - Ensure AWS credentials are configured
   - Check Serverless Framework version compatibility
   - Verify the Go version matches the Lambda runtime

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.