# Этап 3: Аутентификация (auth)

## Что реализовано

Добавлен полный auth-модуль на Go:

- `internal/repository/user.go`
  - `FindByEmail`, `FindByID`, `Create`
- `internal/repository/organization.go`
  - `Create`
- `internal/service/auth.go`
  - `Register` — создание организации и admin-пользователя в одной транзакции
  - `Login` — проверка пары email/password
  - генерация JWT
- `internal/handler/auth.go`
  - thin handlers: decode/validate -> service -> json response
- `internal/middleware/auth.go`
  - JWT middleware (`Authorization: Bearer ...`)
  - загрузка claims в request context
- `cmd/server/main.go`
  - wiring репозиториев/сервисов/хендлеров через DI

Пароль хешируется через `bcrypt` (cost 12).

## Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| POST | `/auth/register` | - | Регистрация организации + admin |
| POST | `/auth/login` | - | Вход |
| GET | `/auth/me` | JWT | Текущий пользователь |

## Ответы и ошибки

- Ошибки бизнес-слоя нормализуются через `ServiceError` (`internal/service/errors.go`).
- Ошибки валидации возвращаются в едином формате `ErrorResponse`.

## Тесты

- `internal/middleware/auth_test.go`
  - отклонение без header
  - пропуск с валидным токеном
  - проверка claims в context
- `internal/router/router_test.go`
  - smoke по health/CORS после подключения auth deps

Проверка: `go test ./...` — PASS.
