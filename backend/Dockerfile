# Build stage
FROM golang:1.25.4-alpine AS build
WORKDIR /app

# Copy modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary from the cmd directory
RUN go build -o app ./cmd

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/app .
EXPOSE 8080
CMD ["./app"]
