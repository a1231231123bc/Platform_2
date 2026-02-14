# Этап 5: Подрядчики (contractors)

## Что реализовано

Добавлен полноценный модуль подрядчиков:

- `internal/repository/contractor.go`
  - create/get/update
  - список с фильтрами и пагинацией
  - история взаимодействий (offers + assignments)
- `internal/repository/contractor_blacklist.go`
  - добавление в blacklist
  - проверка существования
- `internal/service/contractors.go`
  - бизнес-валидации и ошибки (`conflict/not found/bad request`)
- `internal/handler/contractors.go`
  - маршруты CRUD + blacklist + history

## Фильтры

`GET /contractors` поддерживает:

- `region` (по массиву `regions`)
- `type` (`INDIVIDUAL|IP|LLC`)
- `isAvailable` (bool)
- `search` (`ILIKE` по имени/телефону)
- `page`, `limit`

## Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| POST | `/contractors` | JWT + role | Создание |
| GET | `/contractors` | JWT | Список с фильтрами |
| GET | `/contractors/{id}` | JWT | Карточка подрядчика |
| PATCH | `/contractors/{id}` | JWT + role | Обновление |
| POST | `/contractors/{id}/blacklist` | JWT + role | Добавление в ЧС |
| GET | `/contractors/{id}/history` | JWT | История взаимодействий |

`role` для mutate-операций: `ADMIN | MANAGER`.

## Тесты

- `internal/handler/contractors_test.go`
  - parse pagination
  - parse query filters
  - invalid bool parsing
- `go test ./...` — PASS.
