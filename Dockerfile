# ==========================================
# STAGE 1: BUILDER
# ==========================================
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the statically linked Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-api .

# ==========================================
# STAGE 2: RUNNER
# ==========================================
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/inventory-api .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./inventory-api"]