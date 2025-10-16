# Album Manager - Full Stack Application

A modern full-stack album management application with a Go Gin + PostgreSQL backend and Next.js + TypeScript frontend.

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
  frontend/         Next.js web application
    app/            Next.js app directory
      page.tsx      Main page component
      layout.tsx    Root layout
    types/          TypeScript type definitions
    services/       API service layer
    tests/e2e/      Playwright E2E tests
```

## Technologies

### Backend
- **Go 1.19+** - Server language
- **Gin** - HTTP web framework
- **pgx/v5** - PostgreSQL driver with connection pooling
- **PostgreSQL** - Database

### Frontend
- **Next.js 15** - React framework with SSR/SSG
- **TypeScript** - Type-safe programming language
- **React 19** - UI library
- **Tailwind CSS** - Utility-first CSS framework
- **Biome** - Fast linter and formatter
- **Playwright** - E2E testing framework

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

# Configure environment (copy and edit .env.example to .env.local)
cp .env.example .env.local

# Start development server
npm run dev
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

# Start development server with Turbopack
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Lint code with Biome
npm run lint

# Lint and fix code
npm run lint:fix

# Format code with Biome
npm run format

# Run E2E tests with Playwright
npm test
```

## Testing

### Backend Tests (Go)

The backend includes comprehensive tests for the repository layer.

**Setup Test Database:**
```bash
# Create test database (only needed once)
docker exec albums-postgres psql -U postgres -c "CREATE DATABASE albums_test;"
```

**Test Structure:**
```bash
backend/
  repository/
    album_repository_test.go    # Test repository layer (86% coverage)
```

**Test Coverage:**
- Create album
- Find all albums
- Find album by ID
- Update album
- Delete album (soft delete)
- Error handling (not found cases)
- Context cancellation
- Concurrent operations

**Run backend tests:**
```bash
cd backend
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -cover ./...             # With coverage (86%)
go test ./repository -v          # Test specific package
```

### Frontend Tests (Playwright E2E)

The frontend includes comprehensive E2E tests using Playwright. The test file `tests/e2e/albums.spec.ts` includes:

**Test Coverage:**
- ✅ Display Album Manager title
- ✅ Create new album
- ✅ Edit existing album
- ✅ Delete album with confirmation
- ✅ Cancel editing
- ✅ Form validation
- ✅ Album count updates

**Prerequisites:**
Both backend and frontend must be running for E2E tests:
```bash
# Terminal 1: Start backend
cd backend && go run main.go

# Terminal 2: Start frontend
cd frontend && npm run dev

# Terminal 3: Run tests
cd frontend && npm test
```

**Run frontend tests:**
```bash
cd frontend

# Run all tests
npm test

# Run in headed mode (see browser)
npx playwright test --headed

# Run specific test file
npx playwright test albums.spec.ts

# Run tests in UI mode
npx playwright test --ui

# Show test report
npx playwright show-report
```

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

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
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
- **Next.js App Router** - Modern file-based routing with SSR/SSG
- **TypeScript** - Type safety and better developer experience
- **Component-based** - Reusable React components
- **State Management** - React hooks (useState, useEffect)
- **Service Layer** - Centralized API communication with type safety
- **Error Handling** - User-friendly error messages
- **Responsive Design** - Tailwind CSS utility classes
- **Biome** - Fast linting and formatting
- **Playwright** - Comprehensive E2E testing

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
cd frontend && npm run dev
```

### Run Tests
```bash
# Backend tests
cd backend && go test ./...

# Frontend E2E tests (requires backend + frontend running)
cd frontend && npm test
```

### Lint and Format
```bash
# Backend (Go fmt)
cd backend && go fmt ./...

# Frontend (Biome)
cd frontend && npm run lint
cd frontend && npm run format
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

Made with ❤️ using Go, Gin, PostgreSQL, Next.js, TypeScript, and React
