# ----------- BUILD STAGE -------------
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go.mod dan go.sum
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy semua source code Go
COPY server/ ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./cmd

# ----------- RUN STAGE --------------
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary
COPY --from=builder /main .

# Copy migration files
COPY server/cmd/migrate/migrations ./cmd/migrate/migrations

# Load env vars at runtime via container env, not via .env file

EXPOSE 8080
CMD ["./main"]
