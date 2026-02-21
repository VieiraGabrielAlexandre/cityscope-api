package httpserver

import (
	"encoding/json"
	"net/http"
)

// DocsUIHandler serves a minimal Swagger UI that loads the OpenAPI spec from /openapi.json
func DocsUIHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>CityScope API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: '/openapi.json',
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis],
      layout: 'BaseLayout'
    });
  </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

// openAPISpec returns a minimal OpenAPI 3.0 spec for the current endpoints.
func openAPISpec() map[string]any {
	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "CityScope API",
			"version":     "1.0.0",
			"description": "API para consulta de localidades e cidades com dados do IBGE.\n\nEndpoints protegidos requerem Bearer token.",
		},
		"servers": []any{
			map[string]any{"url": "/"},
		},
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
			"schemas": map[string]any{
				"Error": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"error": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"code":       map[string]any{"type": "string", "example": "BAD_REQUEST"},
								"message":    map[string]any{"type": "string", "example": "Invalid parameter"},
								"request_id": map[string]any{"type": "string", "example": "f1a2b3c4d5e6f7g8"},
							},
						},
					},
				},
			},
		},
		"security": []any{
			map[string]any{"bearerAuth": []any{}},
		},
		"paths": map[string]any{
			"/health": map[string]any{
				"get": map[string]any{
					"summary":     "Health check",
					"description": "Retorna status OK para verificação de saúde do serviço.",
					"security":    []any{},
					"responses": map[string]any{
						"200": map[string]any{
							"description": "OK",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"type": "object"},
								},
							},
						},
					},
				},
			},
			"/v1/locations/states": map[string]any{
				"get": map[string]any{
					"summary":     "Lista estados",
					"description": "Lista as UFs do Brasil.",
					"responses": map[string]any{
						"200": map[string]any{
							"description": "Lista de estados",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"type": "object"},
								},
							},
						},
						"401": map[string]any{
							"description": "Unauthorized",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"$ref": "#/components/schemas/Error"},
								},
							},
						},
					},
				},
			},
			"/v1/locations/municipalities": map[string]any{
				"get": map[string]any{
					"summary": "Lista municípios por UF",
					"parameters": []any{
						map[string]any{
							"in":          "query",
							"name":        "state",
							"required":    true,
							"description": "Sigla da UF (ex: SP).",
							"schema":      map[string]any{"type": "string"},
						},
						map[string]any{
							"in":          "query",
							"name":        "q",
							"required":    false,
							"description": "Filtro por nome do município (contém).",
							"schema":      map[string]any{"type": "string"},
						},
					},
					"responses": map[string]any{
						"200": map[string]any{
							"description": "Lista de municípios",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"type": "object"},
								},
							},
						},
						"400": map[string]any{
							"description": "Parâmetro ausente",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"$ref": "#/components/schemas/Error"},
								},
							},
						},
						"401": map[string]any{
							"description": "Unauthorized",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"$ref": "#/components/schemas/Error"},
								},
							},
						},
					},
				},
			},
			"/v1/cities/{ibge_id}": map[string]any{
				"get": map[string]any{
					"summary": "Snapshot de cidade",
					"parameters": []any{
						map[string]any{
							"in":          "path",
							"name":        "ibge_id",
							"required":    true,
							"description": "Código IBGE do município (N6).",
							"schema":      map[string]any{"type": "string"},
						},
					},
					"responses": map[string]any{
						"200": map[string]any{
							"description": "Dados do município",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"type": "object"},
								},
							},
						},
						"400": map[string]any{
							"description": "Requisição inválida",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"$ref": "#/components/schemas/Error"},
								},
							},
						},
						"401": map[string]any{
							"description": "Unauthorized",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"$ref": "#/components/schemas/Error"},
								},
							},
						},
						"502": map[string]any{
							"description": "Falha ao consultar IBGE",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"$ref": "#/components/schemas/Error"},
								},
							},
						},
					},
				},
			},
		},
	}
}

// OpenAPIJSONHandler returns the OpenAPI JSON document.
func OpenAPIJSONHandler(w http.ResponseWriter, r *http.Request) {
	spec := openAPISpec()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(spec)
}
