This project is a small Go scaffold demonstrating how to use oapi-codegen with the Gin framework.

Quick steps to get started:

1. Install oapi-codegen and fetch Go deps

```bash
make deps
```

2. Generate server code from the OpenAPI spec

```bash
# By default Makefile generates a chi-server; change flags to generate gin server code or adapt generated handlers.
# Example (oapi-codegen doesn't include a gin-server generator by default):
# oapi-codegen -generate "types,server,spec" -o generated/api.gen.go api/openapi.yaml

make gen
```

3. Build and run the sample server

```bash
make build
./bin/server
```

Notes:
- oapi-codegen provides multiple generation modes. It doesn't come with a built-in Gin generator, but you can:
  - generate interfaces (server) and implement them with Gin, or
  - generate chi-compatible handlers and adapt to Gin with a small wrapper.
- To generate types and server interfaces, run:
  oapi-codegen -generate "types,server,spec" -o generated/api.gen.go api/openapi.yaml

- If you want a generator that directly targets Gin, there are community templates or you can write a small adapter that translates generated net/http handlers to Gin.

Examples of next steps:
- Generate types and server interfaces, then implement the server interface in `cmd/server/main.go` and register handlers.
- Add unit tests for generated handlers.
