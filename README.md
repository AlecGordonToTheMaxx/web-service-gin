# Albums API - Go Gin CRUD Service

A RESTful API service built with Go, Gin framework, and pgx for managing an albums collection. Features full CRUD operations with PostgreSQL database using jackc/pgx - a high-performance pure Go PostgreSQL driver.

## Project Structure

```
web-service-gin/
  controllers/      HTTP request handlers (dependency injected)
  database/         Database connection and migrations
  models/           Data models with GORM annotations
  repository/       Data access layer (repository pattern)
  routes/           API route definitions
  main.go           Application entry point with graceful shutdown
  .env              Environment configuration (not in git)
  .env.example      Environment template
  README.md         This file
```

## Features

- RESTful API with CRUD operations
- PostgreSQL database with pgx/v5 driver
- Connection pooling for high performance
- Raw SQL queries with full control (no ORM magic)
- Environment-based configuration
- Clean architecture with repository pattern
- Dependency injection for testability
- Context-aware database operations
- Graceful shutdown handling
- Input validation and proper error handling
- Soft delete support
- JSON request/response handling

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- No CGO required (pure Go implementation)

## Setup

1. **Install PostgreSQL**

   Download and install PostgreSQL from [postgresql.org](https://www.postgresql.org/download/)

   Or using Docker:
   ```bash
   docker run --name albums-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15
   ```

2. **Create the database**
   ```bash
   # Connect to PostgreSQL
   psql -U postgres

   # Create database
   CREATE DATABASE albums;

   # Exit psql
   \q
   ```

3. **Clone the repository**
   ```bash
   cd web-service-gin
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```

   Edit `.env` to customize:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=albums
   DB_SSLMODE=disable
   SERVER_PORT=8080
   GIN_MODE=debug
   ```

   Or use a single DATABASE_URL:
   ```env
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/albums?sslmode=disable
   ```

6. **Run the application**
   ```bash
   go run main.go
   ```

   The application will:
   - Connect to PostgreSQL
   - Automatically create the `albums` table
   - Start the server on `http://localhost:8080`

## API Endpoints

### Base URL
```
http://localhost:8080
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/albums` | Get all albums |
| GET | `/albums/:id` | Get album by ID |
| POST | `/albums` | Create new album |
| PUT | `/albums/:id` | Update album |
| DELETE | `/albums/:id` | Delete album |

### Album Model

```json
{
  "id": 1,
  "title": "Album Title",
  "artist": "Artist Name",
  "price": 29.99,
  "created_at": "2025-10-15T12:00:00Z",
  "updated_at": "2025-10-15T12:00:00Z",
  "deleted_at": null
}
```

## API Usage Examples

### 1. Create an Album

**Request:**
```bash
curl -X POST http://localhost:8080/albums \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Blue Train",
    "artist": "John Coltrane",
    "price": 56.99
  }'
```

**Response:**
```json
{
  "ID": 1,
  "CreatedAt": "2025-10-15T12:01:50.019Z",
  "UpdatedAt": "2025-10-15T12:01:50.019Z",
  "DeletedAt": null,
  "title": "Blue Train",
  "artist": "John Coltrane",
  "price": 56.99
}
```

### 2. Get All Albums

**Request:**
```bash
curl http://localhost:8080/albums
```

**Response:**
```json
[
  {
    "ID": 1,
    "CreatedAt": "2025-10-15T12:01:50.019Z",
    "UpdatedAt": "2025-10-15T12:01:50.019Z",
    "DeletedAt": null,
    "title": "Blue Train",
    "artist": "John Coltrane",
    "price": 56.99
  },
  {
    "ID": 2,
    "CreatedAt": "2025-10-15T12:02:01.262Z",
    "UpdatedAt": "2025-10-15T12:02:01.262Z",
    "DeletedAt": null,
    "title": "Jeru",
    "artist": "Gerry Mulligan",
    "price": 17.99
  }
]
```

### 3. Get Album by ID

**Request:**
```bash
curl http://localhost:8080/albums/1
```

**Response:**
```json
{
  "ID": 1,
  "CreatedAt": "2025-10-15T12:01:50.019Z",
  "UpdatedAt": "2025-10-15T12:01:50.019Z",
  "DeletedAt": null,
  "title": "Blue Train",
  "artist": "John Coltrane",
  "price": 56.99
}
```

### 4. Update an Album

**Request:**
```bash
curl -X PUT http://localhost:8080/albums/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Blue Train",
    "artist": "John Coltrane",
    "price": 59.99
  }'
```

**Response:**
```json
{
  "ID": 1,
  "CreatedAt": "2025-10-15T12:01:50.019Z",
  "UpdatedAt": "2025-10-15T12:02:24.206Z",
  "DeletedAt": null,
  "title": "Blue Train",
  "artist": "John Coltrane",
  "price": 59.99
}
```

### 5. Delete an Album

**Request:**
```bash
curl -X DELETE http://localhost:8080/albums/2
```

**Response:**
```json
{
  "message": "Album deleted successfully"
}
```

### 6. Error Responses

**Invalid ID (non-numeric):**
```bash
curl http://localhost:8080/albums/invalid
```
Response (400 Bad Request):
```json
{
  "error": "Invalid album ID"
}
```

**Album Not Found:**
```bash
curl http://localhost:8080/albums/999
```
Response (404 Not Found):
```json
{
  "error": "Album not found"
}
```

**Missing Required Fields:**
```bash
curl -X POST http://localhost:8080/albums \
  -H "Content-Type: application/json" \
  -d '{"title": "Incomplete"}'
```
Response (400 Bad Request):
```json
{
  "error": "Key: 'Album.Artist' Error:Field validation for 'Artist' failed on the 'required' tag\nKey: 'Album.Price' Error:Field validation for 'Price' failed on the 'required' tag"
}
```

## Using PowerShell (Windows)

If you're on Windows and using PowerShell instead of curl:

### Create Album
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/albums" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"title":"Blue Train","artist":"John Coltrane","price":56.99}'
```

### Get All Albums
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/albums" -Method Get
```

### Get Album by ID
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/albums/1" -Method Get
```

### Update Album
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/albums/1" `
  -Method Put `
  -ContentType "application/json" `
  -Body '{"title":"Blue Train","artist":"John Coltrane","price":59.99}'
```

### Delete Album
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/albums/2" -Method Delete
```

## Building for Production

### Build binary
```bash
go build -o albums-api.exe main.go
```

### Run the binary
```bash
./albums-api.exe
```

### Build without CGO (default)
```bash
set CGO_ENABLED=0
go build -o albums-api.exe main.go
```

## Development

### Project follows clean architecture with best practices:

- **models/** - Database models with GORM annotations
- **repository/** - Data access layer (repository pattern for abstraction)
- **controllers/** - HTTP handlers with dependency injection
- **database/** - Database connection and migration logic
- **routes/** - API route definitions

### Key Design Patterns:

1. **Repository Pattern** - Abstracts data access, making it easy to test and swap implementations
2. **Dependency Injection** - Controllers receive dependencies (repo) via constructors
3. **Error Handling** - Distinguishes between not-found, validation, and server errors
4. **Input Validation** - ID parameters are validated before use
5. **Graceful Shutdown** - Handles SIGINT/SIGTERM signals properly
6. **Resource Cleanup** - Database connections are closed on shutdown

### Adding new endpoints:

1. Add model in `models/`
2. Create repository interface and implementation in `repository/`
3. Update migration in `database/database.go`
4. Create controller with injected repository in `controllers/`
5. Wire up dependencies in `main.go`
6. Register routes in `routes/routes.go`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | albums |
| DB_SSLMODE | SSL mode | disable |
| DATABASE_URL | Full connection string (optional) | - |
| SERVER_PORT | Server port number | 8080 |
| GIN_MODE | Gin mode (debug/release/test) | debug |

## Technologies Used

- **[Gin](https://github.com/gin-gonic/gin)** - HTTP web framework
- **[pgx/v5](https://github.com/jackc/pgx)** - High-performance PostgreSQL driver and toolkit
- **[PostgreSQL](https://www.postgresql.org/)** - Advanced open-source relational database
- **[godotenv](https://github.com/joho/godotenv)** - Environment variable loader

## Why pgx over GORM?

- **Performance**: pgx is significantly faster than ORMs for PostgreSQL
- **Control**: Write explicit SQL queries - no hidden behavior or N+1 queries
- **PostgreSQL-specific**: Takes full advantage of PostgreSQL features
- **Connection Pooling**: Built-in pgxpool for efficient connection management
- **Type Safety**: Strong typing with PostgreSQL-specific types
- **Lightweight**: No ORM overhead or reflection magic

## Notes

- This project uses `jackc/pgx/v5`, a pure Go PostgreSQL driver (no CGO required)
- Database uses soft deletes (records are marked as deleted, not removed)
- Timestamps are managed via PostgreSQL DEFAULT values and application logic
- Connection pooling is handled automatically by pgxpool
- All database operations are context-aware for proper cancellation/timeout handling
- The `.env` file is excluded from git for security

## License

MIT
