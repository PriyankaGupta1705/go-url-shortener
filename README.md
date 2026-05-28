# Go URL Shortener

A production-grade URL shortener built with Go, PostgreSQL, and Redis.

## Tech Stack
- **Go 1.25** + Gin HTTP framework
- **PostgreSQL** — persistent storage
- **Redis** — caching layer for fast redirects
- **Docker & Docker Compose** — containerised setup

## Architecture
Redirect request → Redis cache → Hit: instant return
                              → Miss: PostgreSQL → cache in Redis → return

## Quick Start
git clone https://github.com/YOUR_USERNAME/go-url-shortener
cd go-url-shortener
cp .env.example .env     # fill in your values
make local               # starts app + postgres + redis

## API

| Method | Endpoint       | Description          |
|--------|---------------|----------------------|
| GET    | /health        | Health check         |
| POST   | /shorten       | Shorten a URL        |
| GET    | /:code         | Redirect to URL      |
| GET    | /stats/:code   | View visit stats     |

## Example
# shorten
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com"}'

# stats
curl http://localhost:8080/stats/abc123
