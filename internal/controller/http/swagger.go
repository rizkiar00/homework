package http

import (
	_ "embed"
	nethttp "net/http"

	"github.com/labstack/echo/v4"
)

//go:embed contract.yaml
var contractYAML string

func registerSwaggerRoutes(e *echo.Echo) {
	e.GET("/contract.yaml", contract)
	e.GET("/swagger", swaggerRedirect)
	e.GET("/swagger/index.html", swaggerIndex)
}

func contract(ctx echo.Context) error {
	return ctx.Blob(nethttp.StatusOK, "application/yaml; charset=utf-8", []byte(contractYAML))
}

func swaggerRedirect(ctx echo.Context) error {
	return ctx.Redirect(nethttp.StatusMovedPermanently, "/swagger/index.html")
}

func swaggerIndex(ctx echo.Context) error {
	return ctx.HTML(nethttp.StatusOK, `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Homework API Swagger</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
      body { margin: 0; background: #ffffff; }
      .swagger-ui .topbar { display: none; }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = function () {
        window.ui = SwaggerUIBundle({
          url: "/contract.yaml",
          dom_id: "#swagger-ui",
          deepLinking: true,
          presets: [SwaggerUIBundle.presets.apis],
          layout: "BaseLayout"
        });
      };
    </script>
  </body>
</html>`)
}
