# Stage 1: Build stage
FROM golang:1.22.5 AS builder

WORKDIR /app

# Copy the rest of the application
COPY . .
RUN go mod download
RUN go mod tidy

# Optionally copy the .env file if needed
COPY .env .

# Build the application and output the binary to /app/myapp
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp ./cmd

RUN chmod +x ./myapp

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .

# Copy configuration files if needed
COPY --from=builder /app/config/model.conf ./config/model.conf
COPY --from=builder /app/config/policy.csv ./config/policy.csv

# Optionally copy the .env file
COPY --from=builder /app/.env .
# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
