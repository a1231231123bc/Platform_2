# Этап 4: Пользователи и организации

## Что реализовано

Добавлены модули `users` и `organizations` с organization-scope и role-based доступом.

### Users

- `internal/service/users.go`
  - список пользователей организации
  - получение по id
  - обновление
  - удаление
- `internal/handler/users.go`
  - `GET /users`
  - `GET /users/{id}`
  - `PATCH /users/{id}`
  - `DELETE /users/{id}`
- `internal/repository/user.go`
  - добавлены scoped-методы:
    - `FindByIDAndOrganization`
    - `UpdateInOrganization`
    - `DeleteInOrganization`

### Organizations

- `internal/service/organizations.go`
  - получение и обновление организации
  - проверка доступа к organization id
- `internal/handler/organizations.go`
  - `GET /organizations/{id}`
  - `PATCH /organizations/{id}`

### Roles middleware

- `internal/middleware/roles.go`
  - `RequireRoles(...)`
- `internal/middleware/roles_test.go`
  - unauthorized / forbidden / allowed кейсы

## Доступ к маршрутам

- JWT обязателен для users/org routes.
- `PATCH /users/{id}` и `PATCH /organizations/{id}`: `ADMIN | MANAGER`.
- `DELETE /users/{id}`: `ADMIN`.

## DTO

- `internal/dto/user.go` — `UpdateUserRequest`
- `internal/dto/organization.go` — `UpdateOrganizationRequest`

## Тесты

- middleware tests (`roles_test.go`) + общий `go test ./...` — PASS.
