# Use the official Golang image as the build stage.
FROM golang:1.24-alpine as builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code.
COPY . .

# Build the Go application.
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a minimal image for the final container.
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/main .

# Expose port 8080 (the port your app listens on).
EXPOSE 8080

# Command to run the application.
CMD ["./main"]
