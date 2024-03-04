# Description: Dockerfile for building the go application
FROM golang:1.19-alpine3.16 AS builder

# Set the current working directory inside the container
ARG MYSQL_DATABASE
ARG MYSQL_USER
ARG MYSQL_PASSWORD
ARG MYSQL_HOST
ARG MYSQL_PORT
ARG RUN_MIGRATIONS

WORKDIR /app

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o main .

# Run the application
FROM alpine:3.16

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

EXPOSE 8080

# Command to run the application
CMD ["./main"]