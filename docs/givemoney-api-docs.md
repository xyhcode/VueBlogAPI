# GiveMoney API Documentation

## Overview
API endpoints for managing donation/sponsor records (打赏记录).

## Base URL
```
http://localhost:8091/api
```

## Public Endpoints

### Get All Donation Records
Retrieve all donation records with pagination support.

**Endpoint:** `GET /public/givemoney`

**Authentication:** None required

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | integer | No | 1 | Page number (must be ≥ 1) |
| pageSize | integer | No | 10 | Items per page (1-100) |

**Example Request:**
```bash
curl "http://localhost:8091/api/public/givemoney?page=1&pageSize=10"
```

**Response:**
```json
{
  "code": 200,
  "message": "获取打赏记录成功",
  "data": {
    "list": [
      {
        "id": 1,
        "nickname": "Anonymous",
        "figure": 50,
        "created_at": "2026-01-24T10:00:00Z",
        "updated_at": "2026-01-24T10:00:00Z",
        "deleted_at": null
      }
    ],
    "total": 45,
    "pageNum": 1,
    "pageSize": 10
  }
}
```

**Response Fields:**
- `list`: Array of donation records
- `total`: Total number of records
- `pageNum`: Current page number
- `pageSize`: Number of items per page

## Authenticated Endpoints

### Create Donation Record
Create a new donation record.

**Endpoint:** `POST /givemoney`

**Authentication:** Required (JWT + Admin)

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "nickname": "John Doe",
  "figure": 100
}
```

**Body Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| nickname | string | Yes | Donator's nickname (1-50 characters) |
| figure | integer | Yes | Donation amount |

**Example Request:**
```bash
curl -X POST "http://localhost:8091/api/givemoney" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{"nickname": "Alice", "figure": 75}'
```

**Success Response (200):**
```json
{
  "code": 200,
  "message": "创建打赏记录成功",
  "data": {
    "id": 1,
    "nickname": "Alice",
    "figure": 75,
    "created_at": "2026-01-24T10:00:00Z",
    "updated_at": "2026-01-24T10:00:00Z",
    "deleted_at": null
  }
}
```

### Update Donation Record
Update an existing donation record.

**Endpoint:** `PUT /givemoney/{id}`

**Authentication:** Required (JWT + Admin)

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | integer | Yes | Record ID to update |

**Request Body:**
```json
{
  "nickname": "John Smith",
  "figure": 150
}
```

**Example Request:**
```bash
curl -X PUT "http://localhost:8091/api/givemoney/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{"nickname": "John Smith", "figure": 150}'
```

**Success Response (200):**
```json
{
  "code": 200,
  "message": "更新打赏记录成功",
  "data": {
    "id": 1,
    "nickname": "John Smith",
    "figure": 150,
    "created_at": "2026-01-24T10:00:00Z",
    "updated_at": "2026-01-24T11:00:00Z",
    "deleted_at": null
  }
}
```

### Delete Donation Record
Delete a donation record (soft delete).

**Endpoint:** `DELETE /givemoney/{id}`

**Authentication:** Required (JWT + Admin)

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | integer | Yes | Record ID to delete |

**Example Request:**
```bash
curl -X DELETE "http://localhost:8091/api/givemoney/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Success Response (200):**
```json
{
  "code": 200,
  "message": "删除打赏记录成功",
  "data": null
}
```

## Error Responses

### Common Error Codes:

**400 Bad Request**
```json
{
  "code": 400,
  "message": "参数错误: nickname is required",
  "data": null
}
```

**401 Unauthorized**
```json
{
  "code": 401,
  "message": "未授权",
  "data": null
}
```

**404 Not Found**
```json
{
  "code": 404,
  "message": "记录不存在",
  "data": null
}
```

**500 Internal Server Error**
```json
{
  "code": 500,
  "message": "获取打赏记录失败: database error",
  "data": null
}
```

## Data Models

### GiveMoney Object
```json
{
  "id": 1,
  "nickname": "Donator Name",
  "figure": 100,
  "created_at": "2026-01-24T10:00:00Z",
  "updated_at": "2026-01-24T10:00:00Z",
  "deleted_at": null
}
```

**Fields:**
- `id` (integer): Unique identifier
- `nickname` (string): Donator's nickname (1-50 characters)
- `figure` (integer): Donation amount
- `created_at` (string): Creation timestamp (ISO 8601)
- `updated_at` (string): Last update timestamp (ISO 8601)
- `deleted_at` (string/null): Deletion timestamp (null if not deleted)

## Rate Limiting
No explicit rate limiting is implemented. Standard application limits may apply.

## Notes
- All timestamps are in UTC format
- Soft delete is implemented - deleted records are not physically removed
- Authentication uses JWT tokens with admin privileges for write operations
- Pagination defaults to 10 items per page with a maximum of 100