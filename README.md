# Albums API - Go Gin CRUD Service

A RESTful API service built with Go, Gin framework, and GORM for managing an albums collection. Features full CRUD operations with SQLite database (using pure Go - no CGO required).

## Project Structure

```
web-service-gin/
  controllers/      Request handlers for albums
  database/         Database connection and migrations
  models/           Data models with GORM annotations
  routes/           API route definitions
  main.go           Application entry point
  .env              Environment configuration (not in git)
  .env.example      Environment template
  README.md         This file
```

## Features

- RESTful API with CRUD operations
- GORM ORM with SQLite database
- Pure Go SQLite driver (no CGO dependency)
- Environment-based configuration
- Clean architecture (controllers/models/routes pattern)
- Soft delete support
- JSON request/response handling

## Prerequisites

- Go 1.19 or higher
- No CGO required (works on Windows without gcc)

## Setup

1. **Clone the repository**
   ```bash
   cd web-service-gin
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```

   Edit `.env` to customize:
   ```env
   DB_NAME=albums.db
   SERVER_PORT=8080
   GIN_MODE=debug
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080` (or your configured port).

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
  "ID": 1,
  "CreatedAt": "2025-10-15T12:00:00Z",
  "UpdatedAt": "2025-10-15T12:00:00Z",
  "DeletedAt": null,
  "title": "Album Title",
  "artist": "Artist Name",
  "price": 29.99
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

### Project follows clean architecture:

- **models/** - Database models with GORM annotations
- **database/** - Database connection and migration logic
- **controllers/** - Business logic and request handlers
- **routes/** - API route definitions

### Adding new endpoints:

1. Add model in `models/`
2. Update migration in `database/database.go`
3. Create controller in `controllers/`
4. Register routes in `routes/routes.go`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_NAME | SQLite database file name | albums.db |
| SERVER_PORT | Server port number | 8080 |
| GIN_MODE | Gin mode (debug/release/test) | debug |

## Technologies Used

- **[Gin](https://github.com/gin-gonic/gin)** - HTTP web framework
- **[GORM](https://gorm.io/)** - ORM library
- **[modernc.org/sqlite](https://gitlab.com/cznic/sqlite)** - Pure Go SQLite driver
- **[godotenv](https://github.com/joho/godotenv)** - Environment variable loader

## Notes

- This project uses `modernc.org/sqlite`, a pure Go SQLite driver that doesn't require CGO
- Database uses soft deletes (records are marked as deleted, not removed)
- All timestamps are managed automatically by GORM
- The `.env` file is excluded from git for security

## License

MIT
