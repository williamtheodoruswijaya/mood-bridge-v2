# ----------- BUILD STAGE -------------
FROM golang:1.24.2 AS builder

WORKDIR /app

# Copy go.mod dan go.sum dari server
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy semua source code dari server
COPY server/ ./

# Build binary dari folder cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd

# ----------- RUN STAGE --------------
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy hasil build binary
COPY --from=builder /app/main .

# Copy migration files
COPY server/cmd/migrate/migrations ./cmd/migrate/migrations

EXPOSE 8080

CMD ["./main"]
