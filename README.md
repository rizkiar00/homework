# Homework Go Backend

A Go backend playground built to practice a lightweight version of a work-style backend architecture. The project uses Echo, a simple Clean Architecture layout, PostgreSQL through GORM, JWT authentication, Swagger UI, and a private CRUD example for `test_table`.

## Live Deployment

The project is deployed on Google Cloud Platform using Cloud Run.

```text
https://rizkistarter.com/swagger
```

The deployed database uses PostgreSQL on Neon Free Tier. The production deployment uses `DATABASE_URL` for the hosted PostgreSQL connection, while local development can still use a local PostgreSQL instance through `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USERNAME`, and `DB_PASSWORD`.

## Tech Stack

- Go
- Echo (`github.com/labstack/echo/v4`)
- GORM with PostgreSQL driver
- PostgreSQL
- Google Cloud Run for deployment
- Neon PostgreSQL Free Tier for the deployed database
- JWT (`github.com/golang-jwt/jwt/v5`)
- bcrypt for password hashing
- `dig` for dependency injection
- `godotenv` and `caarlos0/env` for configuration
- Swagger UI and server route generation from an OpenAPI contract

## Features

- Health endpoints:
  - `GET /health`
  - `GET /readiness`
- Authentication:
  - `POST /auth/register`
  - `POST /auth/login`
  - `POST /auth/logout` private
  - `GET /auth/me` private
- Private `test_db` CRUD:
  - `POST /test_db`
  - `GET /test_db`
  - `GET /test_db/{test_id}`
  - `PUT /test_db/{test_id}`
  - `DELETE /test_db/{test_id}`
- Action based access management:
  - `GET /actions`
  - `POST /roles`
  - `PUT /roles/{role_id}`
  - `PUT /roles/{role_id}/actions`
  - `PUT /users/{user_id}/role`
- Pagination and ordering for `GET /test_db`
- Swagger UI:
  - `/swagger/index.html`
  - `/contract.yaml`

## Project Structure

```text
cmd/server                  Application entry point
pkg/config                  Environment config loader
pkg/di                      Dependency injection wiring
pkg/resource                Resource bootstrap, such as DB and logger
pkg/token                   JWT generation and parsing
internal/controller/http    Generated API router, controllers, middleware, Swagger, OpenAPI contract
internal/usecase            Business logic
internal/repository/db      PostgreSQL repositories
internal/entity             GORM entities
internal/model              Request and response models
migrations/db               SQL migration files
files/codegen               oapi-codegen config
scripts/http.sh             oapi-codegen generation script
```

The active server is `cmd/server/main.go`. HTTP routes are registered from `internal/controller/http`.

## Prerequisites

- Go installed
- PostgreSQL running
- A PostgreSQL database, for example `postgres`
- A PostgreSQL user, for example `postgres`

## Environment Setup

Create `.env` from `.env.example`, then adjust it for your machine.

```env
APP_NAME=homework-api
APP_ENV=local
APP_HOST=127.0.0.1
APP_PORT=8081
APP_SHUTDOWN_TIMEOUT_SECONDS=10

HTTP_CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
HTTP_TIMEOUT_SECONDS=30
HTTP_BODY_LIMIT=1M
HTTP_RATE_LIMIT_REQUESTS_PER_MINUTE=60
HTTP_RATE_LIMIT_BURST=60
HTTP_AUTH_RATE_LIMIT_REQUESTS_PER_MINUTE=5
HTTP_AUTH_RATE_LIMIT_BURST=5

REDIS_ENABLED=false
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_USERNAME=
REDIS_PASSWORD=
REDIS_DB=0
REDIS_DIAL_TIMEOUT_SECONDS=5

DB_DRIVER=postgres
DATABASE_URL=
DB_HOST=auto
DB_PORT=5432
DB_NAME=postgres
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_SCHEMA=public
DB_SSLMODE=disable

JWT_SECRET=change-me
JWT_EXPIRES_IN_SECONDS=3600
```

Use `DB_HOST=auto` for local development. It resolves to the Windows host when running from WSL, and falls back to `127.0.0.1` when running from Windows PowerShell.

Use `DATABASE_URL` for hosted PostgreSQL providers such as Neon, Supabase, or GCP. When `DATABASE_URL` is set, the app ignores `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USERNAME`, `DB_PASSWORD`, `DB_SCHEMA`, and `DB_SSLMODE`.

Example hosted PostgreSQL format:

```env
DATABASE_URL=postgresql://user:password@host/database?sslmode=require
```

Redis is optional. Keep `REDIS_ENABLED=false` when Redis is not needed. Set `REDIS_ENABLED=true` to include Redis in the readiness check.

## Database Setup

Run all pending migrations:

```bash
go run ./cmd/migrate up
```

Run one migration step up:

```bash
go run ./cmd/migrate up -steps=1
```

Rollback one migration step:

```bash
go run ./cmd/migrate down
```

Check current migration version:

```bash
go run ./cmd/migrate version
```

Force a migration version after fixing a dirty migration state:

```bash
go run ./cmd/migrate force -version=1
```

Create a local admin user after migrations:

```bash
go run ./cmd/admin create -username=admin -password=your_password
```

Do not commit real admin credentials. The command hashes the password with bcrypt before saving it to the local database.

## Run The App

```bash
go run ./cmd/server
```

The server will be available at:

```text
http://127.0.0.1:8081
```

Stop the server with `Ctrl+C`. The app will stop accepting new requests, wait for active requests to finish within `APP_SHUTDOWN_TIMEOUT_SECONDS`, and close the database connection.

If the global Go build cache has permission issues on Windows, use a local cache:

```powershell
$env:GOCACHE="E:\Work\homework\.gocache"
go run ./cmd/server
```

## Run With Docker

Create Docker env file:

```bash
cp .env.docker.example .env.docker
```

Build and start the API with PostgreSQL and Redis:

```bash
docker compose up --build app
```

Run migrations:

```bash
docker compose --profile tools run --rm migrate
```

Create or update a local admin user:

```bash
docker compose --profile tools run --rm admin /app/admin create -username=admin -password=your_password
```

Open Swagger:

```text
http://127.0.0.1:8081/swagger/index.html
```

The Compose setup uses:

- API container on internal port `8080`, exposed locally as `8081`.
- PostgreSQL service hostname `postgres`.
- Redis service hostname `redis`.
- Redis AOF persistence enabled for local Docker data.

For Cloud Run or other GCP container platforms, set the app to listen on `0.0.0.0` and use environment variables for PostgreSQL, Redis, and JWT secrets. The app supports the platform-provided `PORT` variable when `APP_PORT` is not set.

## Swagger

Open:

```text
http://127.0.0.1:8081/swagger/index.html
```

Live Swagger:

```text
https://rizkistarter.com/swagger
```

OpenAPI contract:

```text
http://127.0.0.1:8081/contract.yaml
```

For private endpoints, login first, then click `Authorize` in Swagger UI and paste the JWT access token.

## Public API Protection

Swagger stays public so other people can try the API:

```text
/swagger
/swagger/index.html
/contract.yaml
```

The server applies these request guards:

- Global rate limit: `HTTP_RATE_LIMIT_REQUESTS_PER_MINUTE`, default `60` requests per minute per IP.
- Auth rate limit: `HTTP_AUTH_RATE_LIMIT_REQUESTS_PER_MINUTE`, default `5` requests per minute per IP for `POST /auth/register` and `POST /auth/login`.
- Body limit: `HTTP_BODY_LIMIT`, default `1M`.
- Request timeout: `HTTP_TIMEOUT_SECONDS`, default `30` seconds.

When the rate limit is exceeded, the API returns `429 Too Many Requests`.
When the request body is too large, the API returns `413 Payload Too Large`.

## API Contract

The OpenAPI contract lives at:

```text
internal/controller/http/contract.yaml
```

This file is used by Swagger UI to display the available APIs. It is also used by oapi-codegen to generate the Echo route registration file:

```text
internal/controller/http/openapi_server.gen.go
```

Regenerate the route file after changing `contract.yaml`:

```bash
./scripts/http.sh
```

On Windows PowerShell, run the generator directly:

```powershell
$env:GOCACHE="E:\Work\homework\.gocache"
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config files\codegen\server.cfg.yaml -package http -o internal\controller\http\openapi_server.gen.go internal\controller\http\contract.yaml
```

## Runtime Notes

- Auth uses stateless JWT access tokens.
- There is no refresh token flow yet.
- Logout stores the current JWT `jti` in Redis blacklist until the token expires when Redis is enabled.
- Swagger UI loads assets from the `swagger-ui-dist` CDN, so opening Swagger UI needs internet access.
- Redis is wired as an optional app resource and is checked by `/readiness` when `REDIS_ENABLED=true`.

## Redis Usage

Redis is optional. When `REDIS_ENABLED=false`, the app falls back to local behavior and still works without Redis.

Current Redis-backed behavior:

- Shared rate limiting across app instances. If Redis is enabled, global and auth rate limit counters are stored in Redis instead of local memory.
- JWT logout blacklist. `POST /auth/logout` stores the current token `jti` in Redis with TTL until the token expires.
- Action access cache. Private endpoint permission checks are cached in Redis for a short time and invalidated when role access or user role assignment changes.

Redis persistence can be enabled at the Redis server level. Common options are:

- RDB snapshot: Redis writes point-in-time snapshots to disk. It is simple and fast, but recent writes can be lost if Redis crashes before the next snapshot.
- AOF append-only file: Redis appends every write command to disk. It has better durability than RDB, but uses more disk I/O and the AOF file can grow until it is rewritten.
- RDB + AOF: safer than using only memory, but still adds disk usage and operational work.

Persistence is useful for token blacklist data, but it is usually not required for pure cache data. For this project, losing cached action access or rate limit counters is acceptable. Losing token blacklist data means a logged-out JWT could become usable again until it expires, so keep JWT expiry reasonably short or enable Redis persistence in production.

Good next use cases:

- Cache for master data that is read often but rarely changes.
- Temporary values such as OTP, reset password tokens, invite tokens, and short-lived verification data.
- Lightweight queues for non-critical background tasks before introducing a broker such as RabbitMQ.
