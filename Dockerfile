# ------------------------------
# STAGE 1: Build the Go binary
# ------------------------------
FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /workspace

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code
COPY . .

# Build the controller binary for Linux (because it’ll run in a Linux container)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o controller cmd/controller/main.go

# ------------------------------
# STAGE 2: Create the final image
# ------------------------------
FROM alpine:3.18

# Add trusted root certs (needed for HTTPS calls, e.g., Prometheus)
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy the compiled binary from builder stage
COPY --from=builder /workspace/controller .

# Set binary as entrypoint
ENTRYPOINT ["./controller"]
