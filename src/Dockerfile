FROM golang:1.22-alpine AS base

# Development
FROM base AS dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
COPY go.mod go.sum ./
RUN go mod download
CMD ["air", "-c", ".air.toml"]

# Build for production
FROM base AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/tmp/main /app/cmd/api/main.go

# Run for production
FROM base AS runner
WORKDIR /app
COPY --from=builder /app/tmp/main /tmp/main
CMD ["./tmp/main"]
