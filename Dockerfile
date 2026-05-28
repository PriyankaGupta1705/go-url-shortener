FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:latest

WORKDIR /app

# copy binary
COPY --from=builder /app/server .

# copy migrations folder ← this is the missing line
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080

CMD ["./server"]