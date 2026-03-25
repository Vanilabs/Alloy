# ##############################
# # STAGE 1: Build stage (for production)
# ##############################
# FROM golang:1.25.1-alpine AS build

# ENV CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN go build -o app ./cmd

# ##############################
# # STAGE 2: Dev image with Air for hot reload
# ##############################
# FROM golang:1.25.1-alpine AS dev

# # Install Air for hot reload
# RUN go install github.com/air-verse/air@latest

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# COPY .air.toml ./

# EXPOSE 8082

# CMD ["air", "-c", ".air.toml"]

# ##############################
# # STAGE 3: Production image
# ##############################
# FROM alpine:latest AS prod

# RUN apk --no-cache add ca-certificates tzdata

# WORKDIR /app

# # Copy the binary from build stage
# COPY --from=build /app/app .

# EXPOSE 8080

# CMD ["./app"]


##############################
# STAGE 1: Build stage (for production)
##############################
FROM --platform=$BUILDPLATFORM golang:1.25.1-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o app ./cmd

##############################
# STAGE 2: Dev image with Air for hot reload
##############################
FROM golang:1.25.1-alpine AS dev

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .air.toml ./

EXPOSE 8082

CMD ["air", "-c", ".air.toml"]

##############################
# STAGE 3: Production image
##############################
FROM alpine:latest AS prod

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from build stage
COPY --from=build /app/app .

EXPOSE 8080

CMD ["./app"]
