# Go REST API

A minimal REST API in Go.

## Features / Endpoints

- GET /health — returns 200 OK and JSON {"status":"ok"}
- GET /customers — returns the list of customers (JSON array)
- GET /customers/{id} — returns the customer with the numeric id or 404
- POST /customers — create a customer by sending JSON {"name":"...", "role":"...", "email":"...", "phone":"...", "contacted":false}
- PUT /customers/{id} — update an existing customer (same JSON shape as POST)
- DELETE /customers/{id} — delete a customer; returns 204 on success

These are the endpoints implemented in this project.

## Requirements

- Go 1.20+ installed (https://go.dev/dl/)
- (optional) Docker (https://docs.docker.com/get-docker/)

## Quick start (PowerShell)

From the project root:

```powershell
cd "C:\Users\<you>\Downloads\Code\go-rest-api"


# build and run
.\go-rest-api.exe

# or run directly (for development)
go run .
```

The server listens on port 8080 by default. To change it:

```powershell
$env:PORT = "9090"
.\go-rest-api.exe
```

## Run with Docker

Build and run the production image:

```powershell
docker compose build app
docker compose up -d app
```

Run the dev container (bind mounts local source so changes are live):

```powershell
docker compose up app-dev
```

Check logs:

```powershell
docker compose logs -f app
```

Stop and remove:

```powershell
docker compose down
```

## Example requests (PowerShell)


Create a customer:

```powershell
c -H "Content-Type: application/json" -d '{"name":"Jane Roe","role":"Designer","email":"jane@example.com","phone":"555-0404","contacted":false}'
```

List customers:

```powershell
curl.exe http://localhost:8080/customers
```

Get a customer by id:

```powershell
curl.exe http://localhost:8080/customers/1
```

Update a customer:

```powershell
curl.exe -X PUT http://localhost:8080/customers/1 -H "Content-Type: application/json" -d '{"name":"Alice Johnson","role":"Manager","email":"alice@newexample.com","phone":"555-0101","contacted":true}'
```

Delete a customer:

```powershell
curl.exe -X DELETE http://localhost:8080/customers/1
```
