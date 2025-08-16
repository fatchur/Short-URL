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

**Advantages:**
- ✅ Single deployment unit
- ✅ Lower infrastructure costs
- ✅ Simplified monitoring and logging
- ✅ No network latency between services
- ✅ Easier development and testing

**Disadvantages:**
- ❌ Single point of failure
- ❌ Harder to scale individual components
- ❌ Technology coupling

### 2. Microservices Deployment (For Scale)

**When to Use:**
- High traffic requiring independent scaling
- Multiple team development
- Need for technology diversity
- Sufficient budget for infrastructure complexity

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

**Advantages:**
- ✅ Independent scaling
- ✅ Technology diversity
- ✅ Team autonomy
- ✅ Fault isolation
- ✅ Independent deployments

**Disadvantages:**
- ❌ Higher infrastructure costs
- ❌ Network latency between services
- ❌ Complex monitoring and debugging
- ❌ Data consistency challenges

## Migration Path

### Phase 1: Start with Monolith
1. Deploy monolith (`pkg/main.go`)
2. Monitor traffic and resource usage
3. Identify scaling bottlenecks

### Phase 2: Extract Services (When Ready)
1. **Traffic Triggers:**
   - > 1000 requests/minute
   - Different scaling patterns between user auth and URL shortening
   - Team growing beyond 5-8 developers

2. **Budget Considerations:**
   - Can afford separate databases/Redis instances
   - Can afford load balancers and service discovery
   - Can afford monitoring and logging infrastructure

3. **Migration Steps:**
   - Set up separate deployments for each service
   - Implement service-to-service authentication
   - Split databases if needed
   - Update load balancer configuration
   - Gradually migrate traffic

## Configuration

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

## Production Considerations

### Monolith Production
- Use Docker for consistent deployments
- Implement horizontal scaling with load balancer
- Use managed PostgreSQL and Redis
- Set up monitoring and alerting

### Microservices Production
- Use Kubernetes or Docker Swarm
- Implement API Gateway
- Use service mesh for communication
- Implement distributed tracing
- Set up centralized logging

## Monitoring

### Key Metrics to Track
- Request latency and throughput
- Error rates
- Database connection pool usage
- Memory and CPU usage
- User authentication success rates
- URL shortening success rates

### When to Scale Out
- CPU usage > 80% consistently
- Memory usage > 85%
- Response time > 500ms for 95th percentile
- Database connection pool exhaustion
- Different scaling needs between services

## Recommended Timeline

1. **Months 1-6:** Monolith deployment
2. **Months 6-12:** Monitor and optimize monolith
3. **Month 12+:** Consider microservices if traffic justifies

The architecture is designed to support both strategies with minimal code changes, allowing you to evolve based on actual business needs rather than premature optimization.