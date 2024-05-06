# Description: Dockerfile for building the go application
FROM golang:1.19-alpine3.16 AS builder

WORKDIR /app

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

# Set the current working directory inside the container
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run the application
FROM alpine:3.16

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

EXPOSE 8080

# Command to run the application
CMD ["./main"]