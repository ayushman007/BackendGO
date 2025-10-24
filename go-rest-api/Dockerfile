## Multi-stage Dockerfile for building and running the Go REST API

# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /src

# Install git (needed for go modules that fetch from VCS) and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum first for efficient caching
COPY go.mod .
# If you have a go.sum, copy it too. (Not required here)
# COPY go.sum .

# Download modules
RUN go mod download

# Copy rest of the source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/go-rest-api ./

# Final stage: small runtime image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/go-rest-api /usr/local/bin/go-rest-api

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/go-rest-api"]
