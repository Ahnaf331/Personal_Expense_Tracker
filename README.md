# Personal Expense Tracker API

A REST API built with Go, Gin, and SQLite for tracking personal expenses.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- Git (optional)

## Setup

1. Clone or download the project:
   ```bash
   git clone <repo-url>
   cd "Personal Expense Tracker"
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   
   ```

The server starts on `http://localhost:8080`. A `expenses.db` file is auto-created on first run.

## API Endpoints

### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Create account |
| POST | `/api/auth/login` | Login and get token |

### Expenses (requires token)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/expenses` | Add expense |
| GET | `/api/expenses` | List expenses |
| GET | `/api/expenses/:id` | Get one expense |
| PUT | `/api/expenses/:id` | Update expense |
| DELETE | `/api/expenses/:id` | Delete expense |

### Categories (requires token)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/categories` | Create category |
| GET | `/api/categories` | List categories |
| PUT | `/api/categories/:id` | Update category |
| DELETE | `/api/categories/:id` | Delete category |

### Analytics (requires token)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/analytics/monthly` | Monthly summary |
| GET | `/api/analytics/categories` | Spending by category |

## Authentication

All endpoints except `/register` and `/login` require a Bearer token in the header:

```
Authorization: Bearer <token>
```

## Quick Start Example

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'

# Add an expense (use token from login response)
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"category_id":1,"amount":12.50,"description":"Lunch","date":"2026-04-26"}'
```

Or run the included test script (requires bash and curl):

```bash
bash test-api.sh
```

## Query Filters

**GET /api/expenses** supports optional query params:
- `?category_id=1`
- `?start_date=2026-01-01&end_date=2026-04-30`

**GET /api/analytics/monthly** supports:
- `?year=2026&month=4`

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | `default-secret-change-in-production` | JWT signing key |

Set a custom secret before running in production:
```bash
JWT_SECRET=your-strong-secret go run cmd/main.go
```
