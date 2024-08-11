# Akasia 365

## Pre Requisite

1. Minimum Go Version 1.21.1
2. Minimum Postgre Version 14

## Concept

1. Monorepo with Distributed monolith
2. Domain is divided into `services` as "microservices" but connects to a single database
3. `DTO`s are shared for all domain but `Models` are for specific domain only
   1. `DTO` handles REST contract between domain or client
   2. `Models` handles database model
   3. `Models` should **NEVER** be leaked to REST / client / transport layer
   4. Responding to transport layer should always translate the `Models` to appropriate `DTO`
4. `Repository` layer handles communication between service and database and map the data into corresponding `Models`
   1. `pkg/common` contains a generic repository that can be used to quickly create 1:1 repository between model database table/view

## Environment Variables

All environment configuration can be stored using go structures located in

[Config File](./internal/config/config.go)

---

Common environment variables for all services

```bash
DB_SSL_MODE=disable
DB_HOST=localhost
DB_NAME=akasia
DB_PORT=5432
DB_USER=test
DB_PASS=test12345678
FIREBASE_CONFIG=/{workspace}/firebase.json
```

## Authentication

## Unit Testing
1. Install mockgen `go install github.com/golang/mock/mockgen@v1.6.0`
2. Install gomock `go get github.com/golang/mock/gomock`
3. Generate mocked func `mockgen -source=service/abc.service.go -destination=service/mock/abc.service.mock.go`