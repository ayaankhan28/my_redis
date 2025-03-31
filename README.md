# Go Redis with Locust Load Testing
## Docker Repo
```angular2html
ayaankhan2812/go_locust_combined
```
A high-performance key-value cache server implemented in Go with Locust load testing capabilities.

## Prerequisites

- Go 1.20 or later
- Python 3.10 or later
- Docker (optional, for containerized deployment)

## Local Development Setup

### Installing Dependencies

1. Install Go dependencies:
   ```bash
   go mod tidy
   ```

2. Install Python dependencies:
   ```bash
   pip install -r requirements.txt
   ```

### Running the Server

1. Start the Go server:
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:7171`

### Running Locust Load Tests

1. Start Locust:
   ```bash
   locust -f analyser.py --host=http://localhost:7171
   ```
2. Open Locust web interface at `http://localhost:8089`
3. Enter the number of users to simulate and spawn rate
4. Click "Start swarming" to begin the load test

## Docker Deployment

### Building the Docker Image

```bash
docker build -t go-redis-locust .
```

### Running the Container

```bash
docker run -p 7171:7171 -p 8089:8089 go-redis-locust
```

This will:
- Start the Go server on port 7171
- Start Locust on port 8089

Access:
- Go server: `http://localhost:7171`
- Locust web interface: `http://localhost:8089`

## API Endpoints

### GET /get
Retrieve a value by key:
```bash
curl "http://localhost:7171/get?key=testKey"
```

### POST /put
Store a key-value pair:
```bash
curl -X POST http://localhost:7171/put \
  -H "Content-Type: application/json" \
  -d '{"key":"testKey","value":"testValue"}'
```

## Troubleshooting

1. If ports are already in use, modify the port numbers in:
   - `main.go` for the Go server
   - Locust command line arguments
   - Docker port mappings

2. If Docker build fails, ensure:
   - Docker daemon is running
   - You have necessary permissions
   - All required files are in the build context

3. For connection issues:
   - Check if the Go server is running
   - Verify the host URL in Locust configuration
   - Ensure ports are properly exposed in Docker