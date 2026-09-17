# Calculator

Full-stack calculator: a React (TypeScript) UI talks to a Go REST API for every evaluation. Arithmetic is not computed in the browser.

## Requirements

- Go 1.22+
- Node.js 22+
- Docker (optional, for Compose)

## Run locally

Terminal 1 — API:

```bash
cd backend
go test ./...
go run ./cmd/server
```

The server listens on `http://localhost:8080`.

Terminal 2 — UI:

```bash
cd frontend
npm install
npm test
npm run dev
```

Open `http://localhost:5173`. Vite proxies `/api` to the Go process, so the browser stays same-origin.

## Run with Docker

From the repo root:

```bash
docker compose up --build
```

- UI: `http://localhost:3000`
- API (direct): `http://localhost:8080`
- Nginx on the frontend container proxies `/api` to the backend service (no CORS in that path).

## API

`GET /health`

```json
{"status":"ok"}
```

`POST /api/v1/calculate`

```json
{"operation":"divide","a":10,"b":2}
```

Success (`200`):

```json
{"result":5}
```

Error (`400` invalid input, `422` domain rule):

```json
{"error":"division by zero"}
```

| operation   | formula              | notes |
|-------------|----------------------|--------|
| `add`       | a + b                | |
| `subtract`  | a − b                | |
| `multiply`  | a × b                | |
| `divide`    | a ÷ b                | `b == 0` → 422 |
| `power`     | a^b                  | non-finite (overflow) → 422 |
| `sqrt`      | √a                   | `b` omitted; `a < 0` → 422 |
| `percentage`| (a / 100) × b        | “a percent of b” |

### curl examples

```bash
curl -s http://localhost:8080/health

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":2,"b":3}'

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"subtract","a":5,"b":8}'

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"multiply","a":4,"b":2.5}'

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":10,"b":4}'

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"power","a":2,"b":10}'

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","a":9}'

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"percentage","a":25,"b":200}'

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":1,"b":0}'

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","a":-1}'
```

## Tests and coverage

```bash
cd backend && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

cd frontend && npm test -- --coverage
```

Example from this repo:

```
# backend
ok  calculator/internal/calc     coverage: 100.0% of statements
ok  calculator/internal/httpapi  coverage: 92.9% of statements

# frontend (vitest --coverage)
Test Files  2 passed (2)
Tests       9 passed (9)
```

`cmd/server` is process wiring (listen/shutdown) and is not unit-tested. Keep `coverage.out` locally if you want the HTML report (`go tool cover -html=coverage.out`).

## Design decisions

- **One calculate endpoint** instead of seven routes. Operations share validation, JSON shape, and tests; new ops are a enum + switch case.
- **Pure `calc` package** with no HTTP. The handler only decodes JSON and maps errors. That is the Go equivalent of a Spring service without the container.
- **stdlib `net/http`** (Go 1.22 method-aware mux). No Gin/Echo: fewer moving parts for a single resource.
- **`float64` / JS `number`**. This is not money math. IEEE-754 rounding (e.g. `0.1 + 0.2`) is accepted. The UI rounds for display; the JSON result is the raw float.
- **Non-finite results are errors**. Go cannot JSON-encode `Inf`/`NaN`. Overflow from `power` returns 422 rather than breaking the encoder.
- **Percentage = a% of b**. Documented so the UI and API do not drift.
- **Proxy over CORS for the UI**. Dev: Vite proxy. Prod Compose: Nginx `/api` → backend. CORS on the API remains so `curl` and non-browser clients still work.
- **Vite + React**, not Next.js. No SSR or routing to justify a framework.
- **No eval, no expression strings**. Operation is a closed set. No auth or persistence — the assignment is a calculator, not a platform.

## Layout

```
backend/cmd/server     HTTP process
backend/internal/calc  operations
backend/internal/httpapi  handlers + CORS
frontend/src/api        fetch client
frontend/src/calculator UI
```

## Assumptions

- No login, database, or history.
- A reviewer can run either native toolchains or Compose.

AI prompts used while building this repo are listed in [prompts.md](prompts.md).
