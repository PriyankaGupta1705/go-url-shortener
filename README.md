# Go URL Shortener

A lightweight URL shortener REST API built with Go and Docker.

## Tech Stack
- Go 1.21
- Gin HTTP framework
- Docker & Docker Compose

## Run Locally

### With Docker (recommended)
docker compose up


### Without Docker
go run main.go


## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check |
| POST | /shorten | Shorten a URL |
| GET | /:code | Redirect to original URL |

## Example Usage

### Shorten a URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com"}'

### Response
{"short_code":"abc123","short_url":"http://localhost:8080/abc123"}
