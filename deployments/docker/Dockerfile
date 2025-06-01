# deployments/docker/Dockerfile
FROM golang:1.23.4 AS builder

WORKDIR /app

# Copy source and go mod files together
COPY . .

# Tidy modules with all source present
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy static files and configuration
COPY --from=builder /app/web ./web
COPY --from=builder /app/.env ./

# Create a symbolic link to ensure .env is found in multiple locations
RUN ln -sf /root/.env /root/web/.env

EXPOSE 8080

# Print directory contents and env file for debugging
CMD ["sh", "-c", "ls -la && cat .env && ./main"]
