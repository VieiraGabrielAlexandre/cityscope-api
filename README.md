# `CityScope API 🇧🇷`

A **Go** API that provides an *urban snapshot* of Brazilian cities using official IBGE (Brazilian Institute of Geography and Statistics) data.

The project's goal is to transform the massive volume of public statistical data from IBGE into a simple-to-consume API for applications, dashboards, chatbots, and automation.

> Instead of dealing with the complex SIDRA/IBGE models, CityScope delivers ready-to-use information about a city.

---

## ✨ Features

Currently, CityScope provides:

* **Official Localities**: List Brazilian states and municipalities by state.
* **City Snapshot**: Get a standardized snapshot of a city by its IBGE ID.
* **Population Data**: Official estimated population (IBGE – SIDRA Agregados).
* **Urban Indicators**: Access to various urban metrics (Topic 4714).
* **OpenAPI/Swagger**: Built-in interactive documentation.
* **Security**: Protected endpoints using Bearer Token authentication.
* **Observability**: Structured logging and Request ID tracking.

---

## 🚀 API Documentation

The API includes automatic documentation via Swagger UI.

* **Swagger UI**: `http://localhost:<PORT>/docs`
* **OpenAPI Spec**: `http://localhost:<PORT>/openapi.json`

---

## 🔐 Authentication

All `/v1/*` endpoints require a **Bearer Token**.

Header:
```http
Authorization: Bearer YOUR_TOKEN
```

The token is defined in the `.env` file via the `API_TOKEN` variable.

---

## 📊 Sample Response

### `GET /v1/cities/3550308` (São Paulo)

```json
{
  "data": {
    "ibge_id": 3550308,
    "name": "São Paulo",
    "state": {
      "id": 35,
      "name": "São Paulo",
      "sigla": "SP"
    },
    "population_estimate": {
      "year": 2024,
      "value": 12345678
    },
    "indicators": {
      "area_km2": 1521.11,
      "demographic_density": 7638.12
    }
  }
}
```

---

## ⚠️ Error Handling

The API returns standardized JSON errors for all failed requests, including a `request_id` for troubleshooting.

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid token",
    "request_id": "f1a2b3c4d5e6f7g8"
  }
}
```

Common error codes: `BAD_REQUEST`, `UNAUTHORIZED`, `NOT_FOUND`, `BAD_GATEWAY`.

---

## ⚙️ Configuration

Create a `.env` file in the root directory:

```env
PORT=8080
API_TOKEN=changeme-super-secret
IBGE_BASE_URL=https://servicodados.ibge.gov.br/api
IBGE_TIMEOUT_SECONDS=12
IBGE_CACHE_TTL_SECONDS=3600
```

---

## ▶️ Running Locally

1. Install Go 1.26+
2. Clone the repository
3. Run the application:
   ```bash
   go run ./cmd/api
   ```

The server will start at `http://localhost:8080`.

---

## 🔎 Testing

**Healthcheck (Public):**
```bash
curl http://localhost:8080/health
```

**List States (Protected):**
```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/v1/locations/states
```

**Search Municipalities (Protected):**
```bash
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/v1/locations/municipalities?state=SP&q=camp"
```

---

## 🏗️ Project Structure

* `cmd/api` → Entry point (main).
* `internal/config` → Environment variable configuration.
* `internal/contextutil` → Request ID and Context helpers.
* `internal/handlers` → HTTP handlers and JSON contracts.
* `internal/httpserver` → Router, Middlewares (Auth, Logging, RequestID), and OpenAPI docs.
* `internal/ibge` → IBGE API Client and data mapping.
* `internal/cache` → Simple in-memory TTL cache.

---

## 📊 Observability

The project uses structured logging (`log/slog`) and Request ID tracking to monitor API health and IBGE integration performance.

**Example Log:**
```json
2026-02-21T12:05:00Z INFO request method=GET path=/v1/cities/3550308 status=200 latency=150ms request_id=f1a2b3c4
2026-02-21T12:05:00Z INFO ibge call url=https://... status=200 latency=120ms request_id=f1a2b3c4
```

---

## License

MIT
