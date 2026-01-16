##############################
# STAGE 1: Build stage (for production)
##############################
FROM golang:1.25.1-alpine AS build

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary
RUN go build -o app ./cmd

##############################
# STAGE 2: Dev image with Air for hot reload
##############################
FROM golang:1.25.1-alpine AS dev

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and Air config
COPY . .
COPY .air.toml ./

# Expose the application port (actual port controlled by PORT env var)
EXPOSE 8080

# Run Air for hot reload
CMD ["air", "-c", ".air.toml"]

##############################
# STAGE 3: Production image
##############################
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from build stage
COPY --from=build /app/app .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./app"]
