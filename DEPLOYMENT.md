# Deployment Strategy

This project supports both **monolith** and **microservices** deployment strategies based on your traffic and budget requirements.

## Architecture Overview

### Shared Components
- `domains/` - Shared entities, DTOs, repositories, database logic, and helpers
- Database schema and migrations
- JWT authentication system
- Configuration management

### Services
1. **User Service** (`pkg/user/`) - Authentication and session management
2. **Short-URL Service** (`pkg/short-url/`) - URL shortening functionality
3. **Monolith** (`pkg/main.go`) - Combined deployment of both services

## Deployment Strategies

### 1. Monolith Deployment (Recommended for Start)

**When to Use:**
- Limited cloud budget
- Low to medium traffic
- Single team development
- Simplified deployment and operations

**How to Deploy:**
```bash
# Navigate to monolith directory
cd pkg/

# Install dependencies
go mod tidy

# Run database migrations
go run ../cmd/main.go -d migrate

# Seed initial data
go run ../cmd/main.go -d seed

# Start monolith server
go run main.go
```

**Endpoints:**
- Health: `http://localhost:8080/health`
- User API: `http://localhost:8080/api/v1/user/*`
- Short URL API: `http://localhost:8080/api/v1/url/*`
- Public redirects: `http://localhost:8080/url/{shortCode}`


**How to Deploy:**

**User Service:**
```bash
cd pkg/user/
go mod tidy
go run main.go
# Runs on port 8081
```

**Short-URL Service:**
```bash
cd pkg/short-url/
go mod tidy  
go run main.go
# Runs on port 8080
```

**Endpoints:**
- User Service: `http://localhost:8081/api/v1/user/*`
- Short URL Service: `http://localhost:8080/api/v1/url/*`
- Public redirects: `http://localhost:8080/url/{shortCode}`

Both deployment strategies use the same configuration system (`domains/config/config.go`):

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=shorturl
DB_SSL_MODE=disable

# Security
ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
RATE_LIMIT_DURATION=1m

# Server
PORT=8080  # or 8081 for user service
ENVIRONMENT=development  # or production
```






