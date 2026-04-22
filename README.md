# hello-server

A minimal Go HTTP server that responds with `Hello, World!`.

## Run

```bash
go run .
```

Set `PORT` to override the default `8080`.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Returns `Hello, World!` |
| GET | `/health` | Returns JSON with status and total request count |

## Test

```bash
# hello world
curl http://localhost:8080/

# health check
curl http://localhost:8080/health
```

## Stats

A background goroutine logs request counts to stdout every 30 seconds.
