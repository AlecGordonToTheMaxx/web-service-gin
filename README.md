# Album Manager - Full Stack Application

A full-stack album management application with a Go Gin + PostgreSQL backend and React frontend.

## Project Structure

```
web-service-gin/
  backend/          Go Gin API with PostgreSQL
    controllers/    HTTP request handlers
    database/       Database connection and migrations
    middleware/     CORS and other middleware
    models/         Data models
    repository/     Data access layer
    routes/         API route definitions
    main.go         Application entry point
  frontend/         React web application
    src/
      services/     API service layer
      App.js        Main React component
      App.css       Styling
```

## Technologies

### Backend
- **Go 1.19+** - Server language
- **Gin** - HTTP web framework
- **pgx/v5** - PostgreSQL driver with connection pooling
- **PostgreSQL** - Database

### Frontend
- **React** - UI library
- **JavaScript (ES6+)** - Programming language
- **Fetch API** - HTTP client

## Quick Start

### Prerequisites
- Go 1.19 or higher
- Node.js 16 or higher
- PostgreSQL 12+ or Docker

### 1. Start PostgreSQL

**Using Docker (Recommended):**
```bash
docker run --name albums-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  -d postgres:15

# Create database
docker exec -it albums-postgres psql -U postgres -c "CREATE DATABASE albums;"
```

### 2. Start the Backend

```bash
cd backend

# Install dependencies
go mod download

# Configure environment (copy and edit .env.example to .env)
cp .env.example .env

# Run the server
go run main.go
```

Backend will start on `http://localhost:8080`

### 3. Start the Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm start
```

Frontend will start on `http://localhost:3000`

### 4. Open the Application

Navigate to `http://localhost:3000` in your browser!

## Features

✅ **Create** - Add new albums with title, artist, and price
✅ **Read** - View all albums in a beautiful card layout
✅ **Update** - Edit existing album information
✅ **Delete** - Remove albums with confirmation
✅ **Real-time** - Changes reflect immediately
✅ **Responsive** - Works on desktop and mobile

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/albums` | Get all albums |
| GET | `/albums/:id` | Get album by ID |
| POST | `/albums` | Create new album |
| PUT | `/albums/:id` | Update album |
| DELETE | `/albums/:id` | Delete album (soft delete) |

## Development

### Backend Development

```bash
cd backend

# Run with hot reload (install air first: go install github.com/cosmtrek/air@latest)
air

# Run tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run tests for a specific package
go test ./repository

# Build for production
go build -o albums-api main.go
```

### Frontend Development

```bash
cd frontend

# Start development server
npm start

# Build for production
npm run build

# Run tests (interactive watch mode)
npm test

# Run tests once (for CI/CD)
npm test -- --watchAll=false

# Run tests with coverage
npm test -- --coverage --watchAll=false
```

## Testing

### Backend Tests (Go)

Currently, the backend doesn't have tests yet. To add tests, create files ending in `_test.go`:

Example test structure:
```bash
backend/
  repository/
    album_repository_test.go    # Test repository layer
  controllers/
    album_controller_test.go    # Test HTTP handlers
```

Run backend tests:
```bash
cd backend
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -cover ./...             # With coverage
go test ./repository -v          # Test specific package
```

### Frontend Tests (React)

The frontend includes Jest and React Testing Library. The default test file is `App.test.js`.

Run frontend tests:
```bash
cd frontend

# Interactive watch mode (default)
npm test

# Run once (for CI/CD)
npm test -- --watchAll=false

# With coverage report
npm test -- --coverage --watchAll=false

# Update snapshots
npm test -- -u
```

**Test Commands While Running:**
- Press `a` to run all tests
- Press `f` to run only failed tests
- Press `q` to quit watch mode
- Press `p` to filter by filename
- Press `t` to filter by test name

## Environment Variables

### Backend (.env)
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

### Frontend
The API URL is configured in `frontend/src/services/albumService.js`:
```javascript
const API_URL = 'http://localhost:8080';
```

## Architecture Highlights

### Backend
- **Repository Pattern** - Clean separation of data access
- **Dependency Injection** - Testable and maintainable code
- **Context-aware** - Proper timeout and cancellation handling
- **Connection Pooling** - High-performance database access
- **CORS Enabled** - Allows frontend to communicate with backend
- **Graceful Shutdown** - Handles termination signals properly

### Frontend
- **Component-based** - Reusable React components
- **State Management** - React hooks (useState, useEffect)
- **Service Layer** - Centralized API communication
- **Error Handling** - User-friendly error messages
- **Responsive Design** - Mobile-first CSS

## Deployment

### Backend
```bash
cd backend
go build -o albums-api main.go
./albums-api
```

### Frontend
```bash
cd frontend
npm run build
# Serve the build folder with any static server
```

## Quick Commands Reference

### Start Everything
```bash
# Terminal 1: Backend
cd backend && go run main.go

# Terminal 2: Frontend
cd frontend && npm start
```

### Run Tests
```bash
# Backend tests
cd backend && go test ./...

# Frontend tests (watch mode)
cd frontend && npm test

# Frontend tests (single run)
cd frontend && npm test -- --watchAll=false

# Frontend tests with coverage
cd frontend && npm test -- --coverage --watchAll=false
```

### Build for Production
```bash
# Backend
cd backend && go build -o albums-api main.go

# Frontend
cd frontend && npm run build
```

### Database Commands
```bash
# Start PostgreSQL (Docker)
docker run --name albums-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15

# Create database
docker exec -it albums-postgres psql -U postgres -c "CREATE DATABASE albums;"

# Stop PostgreSQL
docker stop albums-postgres

# Start PostgreSQL (if already created)
docker start albums-postgres

# Access PostgreSQL shell
docker exec -it albums-postgres psql -U postgres -d albums
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT

---

Made with ❤️ using Go, Gin, PostgreSQL, and React
