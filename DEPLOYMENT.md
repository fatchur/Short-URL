# Docker Deployment Guide

This guide covers building Docker images and running services using Docker Compose.

## Prerequisites

- Docker and Docker Compose installed
- Make utility

## Building Docker Images

### Build All Services

```bash
# Build monolith service
make build-monolith

# Build user service
make build-user

# Build short-url service
make build-short-url
```

### Individual Build Commands

```bash
# Monolith (combines user + short-url services)
docker build -t short-url-monolith -f pkg/Dockerfile .

# User service only
docker build -t user-service -f pkg/user/Dockerfile .

# Short-URL service only
docker build -t short-url-service -f pkg/short-url/Dockerfile .
```

## Running Services with Docker Compose

### 1. Database Services (PostgreSQL + Redis)

Start the database services first:

```bash
# Start databases
make up-db

# Stop databases
make down-db
```

**Manual command:**
```bash
docker-compose -f docker-compose.db.yml up -d
docker-compose -f docker-compose.db.yml down
```

### 2. Monolith Service

The monolith combines both user and short-url services in a single container.

```bash
# Start monolith (includes DB migration)
make up-monolith

# Stop monolith
make down-monolith
```

**Manual command:**
```bash
docker-compose -f docker-compose.monolith.yml up -d
docker-compose -f docker-compose.monolith.yml down
```

**Monolith includes:**
- PostgreSQL database
- Redis cache
- Automatic database migration
- Combined user + short-url APIs
- Runs on port 8080

### 3. User Service

Run the user service independently:

```bash
# Start user service (includes DB migration)
make up-user

# Stop user service
make down-user
```

**Manual command:**
```bash
docker-compose -f docker-compose.user.yml up -d
docker-compose -f docker-compose.user.yml down
```

**User service includes:**
- PostgreSQL database
- Automatic database migration
- User authentication APIs
- Runs on port 8080

### 4. Short-URL Service

Run the short-url service independently:

```bash
# Start short-url service (includes DB migration)
make up-short-url

# Stop short-url service
make down-short-url
```

**Manual command:**
```bash
docker-compose -f docker-compose.short-url.yml up -d
docker-compose -f docker-compose.short-url.yml down
```

**Short-URL service includes:**
- PostgreSQL database
- Redis cache
- Automatic database migration
- URL shortening APIs
- Runs on port 8080

## Service Architecture

### Monolith Service
- **Image**: `short-url-monolith`
- **Port**: 8080
- **Dependencies**: PostgreSQL, Redis
- **APIs**: 
  - User API: `/api/v1/user`
  - Short URL API: `/api/v1/url`
  - Health check: `/health`

### User Service
- **Image**: `user-service`
- **Port**: 8080
- **Dependencies**: PostgreSQL
- **APIs**: User authentication and management

### Short-URL Service
- **Image**: `short-url-service`
- **Port**: 8080
- **Dependencies**: PostgreSQL, Redis
- **APIs**: URL shortening and management

## Network Configuration

All services use `network_mode: host` for simplicity, allowing direct access to `localhost:5432` (PostgreSQL) and `localhost:6379` (Redis).

## Database Migration

Each service automatically runs database migrations before starting:
- Migration runs using `cd cmd && go run . -d migrate`
- Services wait for migration completion before starting
- Migration runs once per deployment

## Health Checks

All services include health checks:
- PostgreSQL: `pg_isready -U postgres`
- Redis: `redis-cli ping`
- Application: `/health` endpoint

## Example Workflow

1. **Build all images:**
   ```bash
   make build-monolith
   make build-user
   make build-short-url
   ```

2. **Option A - Run monolith (recommended for development):**
   ```bash
   make up-monolith
   # Test: curl http://localhost:8080/health
   ```

3. **Option B - Run individual services:**
   ```bash
   # Start databases
   make up-db
   
   # Start user service
   make up-user
   # Test: curl http://localhost:8080/api/v1/user
   
   # Or start short-url service instead
   make up-short-url
   # Test: curl http://localhost:8080/api/v1/url
   ```

4. **Cleanup:**
   ```bash
   make down-monolith
   # or
   make down-user
   make down-short-url
   make down-db
   ```

## API Contracts

### Health Check
```
GET /health
```
**Authorization:** None required  
**Rate Limiting:** None  

**cURL Example:**
```bash
curl -X GET http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "service": "short-url-monolith",
  "version": "1.0.0"
}
```

### User Authentication API

#### Create Session (Login)
```
POST /api/v1/user/session
```
**Authorization:** None required  
**Rate Limiting:** **Strict** - 5 requests per 15 minutes per IP  

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/user/session \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "device_info": "Chrome/MacOS",
    "ip_address": "192.168.1.1"
  }'
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "device_info": "Chrome/MacOS",
  "ip_address": "192.168.1.1"
}
```
**Validation Rules:**
- `email`: Required, must be valid email format
- `password`: Required, minimum 8 characters
- `device_info`: Optional string
- `ip_address`: Optional string

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Session created successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_at": "2024-01-01T10:00:00Z"
  }
}
```
**Error Response (401 Unauthorized):**
```json
{
  "status": "error",
  "message": "Invalid credentials"
}
```
**Error Response (400 Bad Request):**
```json
{
  "status": "error",
  "message": "Email and password are required"
}
```
**Error Response (429 Too Many Requests):**
```json
{
  "error": "Too many login attempts. Please try again in 15 minutes."
}
```

### Short URL API

#### Create Short URL
```
POST /api/v1/url/
Authorization: Bearer <access_token>
```
**Authorization:** **Required** - Valid JWT Bearer token  
**Rate Limiting:** **Flexible** - 100 requests per minute per IP  

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/url/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "long_url": "https://example.com/very/long/url/path"
  }'
```

**Request Body:**
```json
{
  "long_url": "https://example.com/very/long/url/path"
}
```
**Validation Rules:**
- `long_url`: Required, must be valid URL format

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Short URL created successfully",
  "data": {
    "id": 1,
    "short_code": "abc123",
    "long_url": "https://example.com/very/long/url/path",
    "user_id": 1
  }
}
```
**Error Response (401 Unauthorized):**
```json
{
  "status": "error",
  "message": "User authentication required"
}
```
**Error Response (400 Bad Request):**
```json
{
  "status": "error",
  "message": "Invalid request body"
}
```
**Error Response (429 Too Many Requests):**
```json
{
  "error": "Too many requests, please try again later"
}
```

#### Get/Redirect Short URL (Authenticated)
```
GET /api/v1/url/{shortCode}
Authorization: Bearer <access_token>
Accept: application/json
```
**Authorization:** **Required** - Valid JWT Bearer token  
**Rate Limiting:** **Flexible** - 100 requests per minute per IP  

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/url/abc123 \
  -H "Accept: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Path Parameters:**
- `shortCode`: Required, the short code identifier

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Short URL retrieved successfully",
  "data": {
    "short_code": "abc123",
    "long_url": "https://example.com/very/long/url/path",
    "user_id": 1
  }
}
```
**Error Response (401 Unauthorized):**
```json
{
  "status": "error",
  "message": "User authentication required"
}
```
**Error Response (404 Not Found):**
```json
{
  "status": "error",
  "message": "Short URL not found or access denied"
}
```
**Error Response (429 Too Many Requests):**
```json
{
  "error": "Too many requests, please try again later"
}
```

#### Public Redirect (No Auth Required)
```
GET /url/{shortCode}
```
**Authorization:** None required  
**Rate Limiting:** None  

**cURL Example:**
```bash
# Follow redirect automatically
curl -L http://localhost:8080/url/abc123

# See redirect response without following
curl -I http://localhost:8080/url/abc123
```

**Path Parameters:**
- `shortCode`: Required, the short code identifier

**Response:** 302 Found redirect to the long URL  
**Error Response (404 Not Found):**
```json
{
  "status": "error",
  "message": "Short URL not found or access denied"
}
```

### Error Response Format
All API errors follow this format:
```json
{
  "status": "error",
  "message": "Error description"
}
```

### Rate Limiting
- **Strict Limiter** (Login): 5 requests per 15 minutes
- **Flexible Limiter** (Other APIs): 100 requests per minute

### Authentication
- Use Bearer token in Authorization header: `Authorization: Bearer <access_token>`
- Token expires as indicated in the login response
- All Short URL APIs require authentication except public redirects

### Complete Workflow Example

**1. Login and get access token:**
```bash
# Login
curl -X POST http://localhost:8080/api/v1/user/session \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response will contain access_token - copy it for next steps
```

**2. Create a short URL:**
```bash
# Replace YOUR_TOKEN with the access_token from step 1
curl -X POST http://localhost:8080/api/v1/url/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "long_url": "https://www.google.com"
  }'

# Response will contain short_code - copy it for next steps
```

**3. Get URL info (authenticated):**
```bash
# Replace YOUR_TOKEN and SHORT_CODE from previous steps
curl -X GET http://localhost:8080/api/v1/url/SHORT_CODE \
  -H "Accept: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**4. Test public redirect:**
```bash
# Replace SHORT_CODE from step 2
curl -L http://localhost:8080/url/SHORT_CODE
# This will redirect to the original long URL
```

## Troubleshooting

- **Port conflicts**: Only run one service at a time since they all use port 8080
- **Database connection**: Services use host networking to connect to `localhost:5432`
- **Migration issues**: Check logs with `docker logs <container-name>`
- **Build failures**: Ensure all dependencies are available and Docker has sufficient resources






