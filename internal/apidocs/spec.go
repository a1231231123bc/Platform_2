package apidocs

// SwaggerJSON is a minimal OpenAPI document for runtime docs endpoint.
const SwaggerJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Platform API",
    "version": "1.0.0",
    "description": "Go backend API for platform migration project"
  },
  "servers": [
    { "url": "/" }
  ],
  "paths": {
    "/": {
      "get": {
        "summary": "Health check",
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    },
    "/auth/login": {
      "post": {
        "summary": "Login",
        "responses": {
          "200": {
            "description": "JWT token and user payload"
          }
        }
      }
    },
    "/analytics/dashboard": {
      "get": {
        "summary": "Dashboard metrics",
        "security": [
          { "bearerAuth": [] }
        ],
        "responses": {
          "200": {
            "description": "Dashboard metrics payload"
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    }
  }
}`

// SwaggerHTML is a minimal UI page that renders the JSON URL in-browser.
const SwaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Platform API Docs</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; margin: 2rem; line-height: 1.5; }
    code { background: #f1f3f5; padding: 0.1rem 0.35rem; border-radius: 4px; }
  </style>
</head>
<body>
  <h1>Platform API Docs</h1>
  <p>OpenAPI JSON: <a href="/api/docs/swagger.json"><code>/api/docs/swagger.json</code></a></p>
  <p>You can import this URL into Swagger UI / Postman / Insomnia.</p>
</body>
</html>`
