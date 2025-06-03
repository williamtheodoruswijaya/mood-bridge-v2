FROM golang:1.24.2 AS builder

WORKDIR /app

# Copy go.mod dan go.sum secara eksplisit
COPY server/go.mod ./go.mod
COPY server/go.sum ./go.sum

RUN go mod download

COPY server/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
