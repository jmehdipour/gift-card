# Use a specific base image with a specific tag
FROM golang:1.19.3-alpine3.16 as builder

# Set the working directory in the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the gift-card binary
RUN go build -o gift-card

# Use a specific base image with a specific tag
FROM alpine:3.16

# Set the working directory in the container
WORKDIR /app

# Copy the gift-card binary from the builder
COPY --from=builder /app/gift-card /app

# Ensure the gift-card binary is executable
RUN chmod +x /app/gift-card

# Create a non-root user with a specific UID and GID to run the application
#RUN adduser -S -u 1000 -G 1000 appuser

# Switch to the non-root user to run the application
#USER appuser

# Expose the required port for the gift-card service
EXPOSE 8080

# Run the gift-card binary with the required arguments
CMD ["./gift-card", "start", "--config", "config.yml"]
