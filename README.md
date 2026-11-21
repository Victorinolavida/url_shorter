# URL Shorter

A URL shortener written in Go, with Redis caching, asynchronous link expiration, and click analytics.

## Features

- Create short links and reuse an existing active link for the same URL.
- Optional expiration for links.
- Redirect by short code.
- Redis cache-aside reads for redirects.
- Background worker that expires stale links and invalidates cache entries.
- Click events with user-agent parsing, bot detection, country, and HMAC-SHA256 hashed IP addresses.
- OpenAPI specification in [`openapi.yaml`](openapi.yaml).

## Tech Stack

| Area | Tool |
| --- | --- |
| Language | Go |
| HTTP router | chi |
| Database | PostgreSQL 17 |
| DB access | pgx + sqlc |
| Migrations | dbmate |
| Cache | Redis 7 |
| Configuration | koanf |
| Logging | log/slog |

## Requirements

- Go 1.26.5 or newer
- Docker and Docker Compose
- [just](https://github.com/casey/just)
- [dbmate](https://github.com/amacneil/dbmate) for migrations
- [sqlc](https://sqlc.dev) only when regenerating query code

## Getting Started

Create `config/config.yaml` using the example in [Configuration](#configuration), then start PostgreSQL and Redis:

```sh
docker compose up -d
```

Apply migrations:

```sh
just migrate
```

Run the API:

```sh
just run
```

Run the expiration worker in a separate terminal:

```sh
just run-worker
```

With the example configuration, the API listens on `http://localhost:5100`. Docker Compose runs only PostgreSQL and Redis; the API and worker run locally.

### Quick Check

```sh
curl http://localhost:5100/health

curl -X POST http://localhost:5100/url-shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/some/long/path","expires_in":"24h"}'

# Replace <code> with the code from the returned URL.
curl -i 'http://localhost:5100/<code>'
```

The create request returns a JSON object containing `url` and `expires_at`. Omit `expires_in` for a link without an expiration; durations use Go syntax, such as `30m` or `24h`. Following the returned URL responds with a `303 See Other` redirect.

## Configuration

Configuration is loaded from `config/config.yaml`, which is not committed. Create that file before running either process:

```yaml
env: local

server:
  host: localhost
  port: 5100
  public_domain: http://localhost:5100

database:
  user: postgres
  password: postgres
  host: localhost
  port: 5432
  name: url_shorten
  min_conns: 2
  max_conns: 2

redis:
  host: localhost
  port: "6379"

worker_duration: 1s
log_level: debug
ip_hash_key: local-dev-key
```

Environment variable names are lowercased and every underscore becomes a dot. For example, these override `server.port`, `database.host`, and `redis.host`:

```sh
SERVER_PORT=8080
DATABASE_HOST=localhost
REDIS_HOST=localhost
```

Keys containing underscores, such as `ip_hash_key` and `server.public_domain`, should be set in YAML because underscores in environment variable names are converted to dots. In local mode, the public domain is derived from the server host and port. In production mode, set `server.public_domain` without a scheme; generated URLs use HTTPS. Keep `ip_hash_key` empty to skip IP hashing, or set a secret value in production.

The migration commands use `DATABASE_URL` separately from the application's database settings. Its default matches the Compose credentials; set it explicitly if you change them.

## API

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Health check |
| `POST` | `/url-shorten` | Create or reuse a short link |
| `GET` | `/{code}` | Redirect to the original URL |

See [`openapi.yaml`](openapi.yaml) for request and response schemas.

## Common Commands

```sh
just run          # run the API
just run-worker   # run the expiration worker
just migrate      # apply pending migrations
just status       # show migration status
just down         # roll back the last migration batch
just clean        # drop and recreate the database
just new <name>   # create a new migration
just schema       # regenerate db/schema.sql from the database
just generate     # run schema + sqlc generate
go test ./...     # run the Go tests
```

Stop the local database and cache with `docker compose down`. Their data is retained in Docker volumes.

## Project Layout

```text
cmd/api/            API entrypoint
cmd/worker/         Expiration worker entrypoint
config/             Configuration loading and YAML defaults
db/migrations/      dbmate migrations
db/query.sql        sqlc queries
db/schema.sql       Schema used by sqlc
db/sqlc/            Generated sqlc code
cache/              Cache interface and Redis implementation
internal/entity/    Domain entities and validation
internal/service/   Business logic
internal/repository/ Database access
internal/handler/   HTTP handlers
internal/http_server/ Server wiring and middleware
internal/httpjson/  JSON request/response helpers
utils/logger/       slog setup
```

## Migrations and Code Generation

Create a migration:

```sh
just new add-something
```

After changing the schema, regenerate the schema dump and sqlc code:

```sh
just generate
```

Do not edit files in `db/sqlc/`; they are generated.

## License

No license has been added yet.
