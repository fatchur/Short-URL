# Inventory Service API

Inventory management service for handling inventory items and distributors.

## Base URL
```
http://localhost:8081
```

## Authentication

All inventory API endpoints require JWT authentication. You need to obtain a token from the User Service first.

### Get Authentication Token

**Endpoint:** `POST http://localhost:8080/api/v1/user/session`

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/user/session \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123",
    "device_info": "curl",
    "ip_address": "127.0.0.1"
  }'
```

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123",
  "device_info": "postman",
  "ip_address": "127.0.0.1"
}
```

**Response:**
```json
{
  "success": true,
  "status": 201,
  "message": "Session created successfully",
  "api_version": "v1",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_at": "2025-10-24T21:34:06.814238+07:00"
  }
}
```

**Use the `access_token` in subsequent requests:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## API Endpoints

### 1. Create Inventory

**Endpoint:** `POST /api/v1/inventory`

**cURL Example:**
```bash
curl -X POST http://localhost:8081/api/v1/inventory \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "distributor_id": 1,
    "name": "Dell XPS 13",
    "description": "Dell XPS 13 laptop with Intel i7 processor",
    "sku": "DELL-XPS13-I7-512",
    "category_id": "electronics",
    "quantity": 15,
    "min_quantity": 3,
    "unit_price": 25000000.00
  }'
```

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <your_token>
```

**Request Body:**
```json
{
  "distributor_id": 1,
  "name": "Dell XPS 13",
  "description": "Dell XPS 13 laptop with Intel i7 processor",
  "sku": "DELL-XPS13-I7-512",
  "category_id": "electronics",
  "quantity": 15,
  "min_quantity": 3,
  "unit_price": 25000000.00
}
```

**Response:**
```json
{
  "success": true,
  "status": 201,
  "message": "Inventory created successfully",
  "api_version": "v1",
  "data": {
    "id": 6,
    "distributor_id": 1,
    "name": "Dell XPS 13",
    "description": "Dell XPS 13 laptop with Intel i7 processor",
    "sku": "DELL-XPS13-I7-512",
    "category_id": "electronics",
    "quantity": 15,
    "min_quantity": 3,
    "unit_price": 25000000,
    "distributor": {
      "id": 1,
      "name": "Tech Solutions Ltd",
      "email": "contact@techsolutions.com",
      "phone_number": "021-1234567",
      "address": "123 Tech Street, Jakarta, Indonesia"
    }
  }
}
```

### 2. Update Inventory

**Endpoint:** `PUT /api/v1/inventory`

**cURL Example:**
```bash
curl -X PUT http://localhost:8081/api/v1/inventory \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "id": 6,
    "distributor_id": 1,
    "name": "Dell XPS 13 Updated",
    "description": "Dell XPS 13 laptop with Intel i7 processor - Updated",
    "sku": "DELL-XPS13-I7-512",
    "category_id": "electronics",
    "quantity": 20,
    "min_quantity": 5,
    "unit_price": 24000000.00
  }'
```

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <your_token>
```

**Request Body:**
```json
{
  "id": 6,
  "distributor_id": 1,
  "name": "Dell XPS 13 Updated",
  "description": "Dell XPS 13 laptop with Intel i7 processor - Updated",
  "sku": "DELL-XPS13-I7-512",
  "category_id": "electronics",
  "quantity": 20,
  "min_quantity": 5,
  "unit_price": 24000000.00
}
```

**Response:**
```json
{
  "success": true,
  "status": 200,
  "message": "Inventory updated successfully",
  "api_version": "v1",
  "data": {
    "id": 6,
    "distributor_id": 1,
    "name": "Dell XPS 13 Updated",
    "description": "Dell XPS 13 laptop with Intel i7 processor - Updated",
    "sku": "DELL-XPS13-I7-512",
    "category_id": "electronics",
    "quantity": 20,
    "min_quantity": 5,
    "unit_price": 24000000
  }
}
```

### 3. Get Inventory by SKU

**Endpoint:** `GET /api/v1/inventory/{sku}`

**cURL Example:**
```bash
curl -X GET http://localhost:8081/api/v1/inventory/DELL-XPS13-I7-512 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Headers:**
```
Authorization: Bearer <your_token>
```

**Example Request:**
```
GET /api/v1/inventory/DELL-XPS13-I7-512
```

**Response:**
```json
{
  "success": true,
  "status": 200,
  "message": "Inventory retrieved successfully",
  "api_version": "v1",
  "data": {
    "id": 6,
    "distributor_id": 1,
    "name": "Dell XPS 13",
    "description": "Dell XPS 13 laptop with Intel i7 processor",
    "sku": "DELL-XPS13-I7-512",
    "category_id": "electronics",
    "quantity": 15,
    "min_quantity": 3,
    "unit_price": 25000000,
    "distributor": {
      "id": 1,
      "name": "Tech Solutions Ltd",
      "email": "contact@techsolutions.com",
      "phone_number": "021-1234567",
      "address": "123 Tech Street, Jakarta, Indonesia"
    }
  }
}
```

### 4. Get All Inventories

**Endpoint:** `GET /api/v1/inventories`

**cURL Example:**
```bash
curl -X GET http://localhost:8081/api/v1/inventories \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Headers:**
```
Authorization: Bearer <your_token>
```

**Response:**
```json
{
  "success": true,
  "status": 200,
  "message": "Inventory list retrieved successfully",
  "api_version": "v1",
  "data": [
    {
      "id": 1,
      "distributor_id": 1,
      "name": "MacBook Pro 14-inch",
      "description": "Apple MacBook Pro with M2 chip, 16GB RAM, 512GB SSD",
      "sku": "MBP-14-M2-16-512",
      "category_id": "electronics",
      "quantity": 25,
      "min_quantity": 5,
      "unit_price": 29999000,
      "distributor": {
        "id": 1,
        "name": "Tech Solutions Ltd",
        "email": "contact@techsolutions.com",
        "phone_number": "021-1234567",
        "address": "123 Tech Street, Jakarta, Indonesia"
      }
    }
  ]
}
```

## Request/Response Field Descriptions

### Create/Update Inventory Request
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `distributor_id` | number | Optional | ID of the distributor |
| `name` | string | Required | Name of the inventory item (max 255 chars) |
| `description` | string | Optional | Description of the item (max 1000 chars) |
| `sku` | string | Required | Unique Stock Keeping Unit (max 100 chars) |
| `category_id` | string | Optional | Category enum: "electronics", "clothing", "food", "books", "furniture", "automotive", "health", "sports", "toys", "home" |
| `quantity` | number | Required | Current stock quantity (min 0) |
| `min_quantity` | number | Optional | Minimum stock threshold (min 0) |
| `unit_price` | number | Required | Price per unit (min 0) |

### Update Inventory Additional Fields
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | number | Required | ID of the inventory item to update |

### Inventory Response
| Field | Type | Description |
|-------|------|-------------|
| `id` | number | Unique identifier |
| `distributor_id` | number | ID of associated distributor |
| `name` | string | Name of the inventory item |
| `description` | string | Description of the item |
| `sku` | string | Stock Keeping Unit |
| `category_id` | string | Category enum value |
| `quantity` | number | Current stock quantity |
| `min_quantity` | number | Minimum stock threshold |
| `unit_price` | number | Price per unit |
| `distributor` | object | Distributor information (when available) |

### Distributor Object
| Field | Type | Description |
|-------|------|-------------|
| `id` | number | Unique identifier |
| `name` | string | Distributor name |
| `email` | string | Distributor email |
| `phone_number` | string | Distributor phone number |
| `address` | string | Distributor address |

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "status": 400,
  "message": "Bad request. Please check your input data.",
  "api_version": "v1"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "status": 401,
  "message": "Unauthorized access. Please provide valid authentication credentials.",
  "api_version": "v1"
}
```

### 404 Not Found
```json
{
  "success": false,
  "status": 404,
  "message": "Inventory item not found.",
  "api_version": "v1"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "status": 500,
  "message": "Internal server error. Please try again later.",
  "api_version": "v1"
}
```

## Available Categories

- `electronics` - Electronic devices and gadgets
- `clothing` - Apparel and fashion items
- `food` - Food and beverage products
- `books` - Books and publications
- `furniture` - Furniture and home decor
- `automotive` - Vehicle parts and accessories
- `health` - Health and medical products
- `sports` - Sports and fitness equipment
- `toys` - Toys and games
- `home` - Home and household items

## Seeded Data

The service comes with pre-seeded data including:
- 3 distributors (Tech Solutions Ltd, Global Electronics Corp, Premium Supplies Inc)
- 5 inventory items (MacBook Pro, iPhone, Nike shoes, Coffee beans, Programming books)
- 3 users for authentication (john@example.com, jane@example.com, bob@example.com)

## Testing

Run repository unit tests:
```bash
make inventory-repository-unit-test
```

## Service Dependencies

- **User Service** (port 8080) - For authentication
- **Database** - PostgreSQL for data persistence
- **Redis** - For caching (optional)