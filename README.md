# Employee Management System

A REST API built with only the Go standard library (`net/http`, `encoding/json`,
`database/sql`, `context`, `sync`) plus a MySQL driver, following a
Handler → Service → Repository layering.

## Project layout

```
cmd/app/main.go             # entrypoint: wiring, server, graceful shutdown
internal/
  handlers/     HTTP handlers (net/http only, no framework)
  services/     business logic + validation
  repositories/ MySQL storage layer behind interfaces
  models/       Employee, Department, request/response, AppError
  middleware/   logging, panic recovery, basic auth
  utils/        validation, pagination, error normalization
  config/       environment-based configuration
migrations/     schema + seed SQL, auto-run by the mysql container
```

## Running with Docker

```
docker compose up --build
```

This starts MySQL (seeded from `migrations/001_init.sql`) and the API on
`http://localhost:8080`. Exported files land in `./exports` on the host.

Default Basic Auth credentials (override with `AUTH_USERS` env var, format
`user:pass,user2:pass2`): `admin` / `admin123`.

## Running locally without Docker

Requires Go 1.22+ and a reachable MySQL instance.

```
go mod tidy
cp .env.example .env   # edit as needed, then export the vars
go run ./cmd/app
```

> Note: this environment has no local Go toolchain, so the code was verified
> by careful review rather than a local `go build` run. Please run
> `go build ./...` (or `docker compose up --build`) to confirm before
> relying on it.

## API

All endpoints except `GET /healthz` require HTTP Basic Auth.

### Employees

| Method | Path | Description |
|---|---|---|
| POST | `/employees` | Create an employee |
| GET | `/employees?limit=&offset=&departmentId=` | Paginated list |
| GET | `/employees/search?keyword=` | Search by name or position |
| GET | `/employees/{id}` | Get one employee |
| PUT | `/employees/{id}` | Partial update |
| DELETE | `/employees/{id}` | Soft delete |
| POST | `/employees/export` | Concurrently export all employees to JSON + CSV |

```
curl -u admin:admin123 -X POST http://localhost:8080/employees \
  -H "Content-Type: application/json" \
  -d '{"name":"Nguyen Van A","age":28,"position":"Backend Engineer","departmentId":1,"salary":2000}'

curl -u admin:admin123 "http://localhost:8080/employees?limit=10&offset=0"

curl -u admin:admin123 -X PUT http://localhost:8080/employees/1 \
  -H "Content-Type: application/json" \
  -d '{"salary":2500}'

curl -u admin:admin123 -X DELETE http://localhost:8080/employees/1

curl -u admin:admin123 "http://localhost:8080/employees/search?keyword=an"

curl -u admin:admin123 -X POST http://localhost:8080/employees/export
```

### Departments

| Method | Path | Description |
|---|---|---|
| POST | `/departments` | Create a department |
| GET | `/departments` | List departments |
| GET | `/departments/{id}/employees` | Employees in a department |

```
curl -u admin:admin123 -X POST http://localhost:8080/departments \
  -H "Content-Type: application/json" -d '{"name":"Engineering"}'

curl -u admin:admin123 http://localhost:8080/departments/1/employees
```

## Error format

All errors share one JSON shape:

```json
{ "error": "Employee not found", "code": 404 }
```

Produced for: invalid IDs/body (400), missing resources (404), duplicate
departments (409), unauthorized requests (401), and internal/DB/timeout
failures (500/504).

## Design notes

- **Layering**: handlers only parse/validate the HTTP transport and call a
  service; services own business rules and per-call `context.WithTimeout`
  guards around DB calls; repositories are the only place SQL is written,
  behind `EmployeeRepository`/`DepartmentRepository` interfaces so services
  are unit-testable with in-memory fakes instead of a real database.
- **Partial update**: `UpdateEmployeeRequest` uses pointer fields; only
  non-nil fields are included in the generated `UPDATE ... SET` clause.
- **Soft delete**: `employees.deleted_at`; all reads filter it out.
- **Concurrency exercise**: `ExportService.ExportAll` fetches employees once,
  then writes the JSON and CSV files in two goroutines coordinated by a
  `sync.WaitGroup`; a `sync.Mutex` guards the shared error slice the
  goroutines may both write to.
- **Routing**: Go 1.22's `net/http.ServeMux` method + `{id}` wildcard
  patterns are used directly — no third-party router. Go's mux picks the more
  specific literal pattern (`/employees/search`) over the wildcard one
  (`/employees/{id}`) automatically.
