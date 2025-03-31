# === Stage 1: Build Go Application ===
FROM golang:1.20 AS builder

# Set working directory
WORKDIR /app

# Copy Go files and dependencies
COPY main.go go.mod go.sum ./
RUN go mod tidy

# Build the Go binary
RUN go build -o main

# === Stage 2: Prepare Final Image ===
FROM python:3.10

# Set working directory
WORKDIR /app

# Install dependencies for Locust
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy built Go server from builder stage
COPY --from=builder /app/main .

# Copy Locust script
COPY analyser.py .

# Expose ports (Go Server & Locust Web UI)
EXPOSE 7171 8089

# Start both services (Go Server & Locust) when the container runs
CMD ["sh", "-c", "./main & locust -f analyser.py --host=http://localhost:7171"]
