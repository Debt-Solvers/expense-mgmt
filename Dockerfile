# Stage 1: Build the application
FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o expense-service ./cmd/expense-service/main.go

# Stage 2: Minimal runtime with `glibc` support
FROM frolvlad/alpine-glibc:latest

WORKDIR /app

# Copy application binary and config
COPY --from=builder /app/expense-service .
COPY configs/config.yaml /app/configs/config.yaml

EXPOSE 8081
ENV PORT=8081
CMD ["./expense-service"]
