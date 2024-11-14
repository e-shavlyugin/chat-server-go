# Use the official Golang image
FROM golang:1.23-alpine

# Install necessary packages including delve
RUN apk add --no-cache git

# Install air
RUN go install github.com/air-verse/air@latest

# Install delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Set air as the default command for live reloading
CMD ["air"]
