# Этап 6: Работы (jobs)

## Что реализовано

Добавлен модуль `jobs` с жизненным циклом статусов и дублированием.

- `internal/repository/job.go`
  - create draft
  - get by id (organization-scope)
  - list с filters + pagination + related counts
  - update только для `DRAFT`
  - статусные переходы
  - duplicate в новый `DRAFT`
- `internal/service/jobs.go`
  - бизнес-правила lifecycle
  - валидация UUID + not found handling
- `internal/handler/jobs.go`
  - endpoints для CRUD/list/lifecycle

## Lifecycle

Поддержаны переходы:

- publish: `DRAFT -> ACTIVE`
- cancel: `DRAFT|ACTIVE|IN_PROGRESS -> CANCELLED`
- complete: `ACTIVE|IN_PROGRESS -> COMPLETED`
- duplicate: копия любой job в новый `DRAFT`

## Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| POST | `/jobs` | JWT + role | Создать draft |
| GET | `/jobs` | JWT | Список с фильтрами |
| GET | `/jobs/{id}` | JWT | Работа по id |
| PATCH | `/jobs/{id}` | JWT + role | Обновить (только draft) |
| POST | `/jobs/{id}/publish` | JWT + role | Публикация |
| POST | `/jobs/{id}/cancel` | JWT + role | Отмена |
| POST | `/jobs/{id}/complete` | JWT + role | Завершение |
| POST | `/jobs/{id}/duplicate` | JWT + role | Дублирование |

`role` для mutate-операций: `ADMIN | MANAGER`.

## Фильтры списка

`GET /jobs`:

- `status` (`DRAFT|ACTIVE|IN_PROGRESS|COMPLETED|CANCELLED`)
- `region`
- `page`, `limit`

В выдаче списка возвращаются счётчики связей:
- `dispatchCount`
- `offerCount`
- `assignmentCount`

## Тесты

- `internal/handler/jobs_test.go`
  - parse query/status/region
  - pagination defaults
  - invalid pagination
- `go test ./...` — PASS.
