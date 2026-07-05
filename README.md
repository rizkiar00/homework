# Homework Go Backend

A Go backend playground built to practice a lightweight version of a work-style backend architecture. The project uses Echo, a simple Clean Architecture layout, PostgreSQL through GORM, JWT authentication, Swagger UI, and a private CRUD example for `test_table`.

## Tech Stack

- Go
- Echo (`github.com/labstack/echo/v4`)
- GORM with PostgreSQL driver
- PostgreSQL
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
  - `GET /auth/me` private
- Private `test_db` CRUD:
  - `POST /test_db`
  - `GET /test_db`
  - `GET /test_db/{id_test}`
  - `PUT /test_db/{id_test}`
  - `DELETE /test_db/{id_test}`
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

DB_DRIVER=postgres
DB_HOST=auto
DB_PORT=5432
DB_NAME=postgres
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_SCHEMA=public

JWT_SECRET=change-me
JWT_EXPIRES_IN_SECONDS=3600
```

Use `DB_HOST=auto` for local development. It resolves to the Windows host when running from WSL, and falls back to `127.0.0.1` when running from Windows PowerShell.

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

## Run The App

```bash
go run ./cmd/server
```

The server will be available at:

```text
http://127.0.0.1:8081
```

If the global Go build cache has permission issues on Windows, use a local cache:

```powershell
$env:GOCACHE="E:\Work\homework\.gocache"
go run ./cmd/server
```

## Swagger

Open:

```text
http://127.0.0.1:8081/swagger/index.html
```

OpenAPI contract:

```text
http://127.0.0.1:8081/contract.yaml
```

For private endpoints, login first, then click `Authorize` in Swagger UI and paste the JWT access token.

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
- There is no server-side logout token revocation yet.
- Swagger UI loads assets from the `swagger-ui-dist` CDN, so opening Swagger UI needs internet access.
