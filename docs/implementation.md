# Реализация Этапов 1-6

Этот файл содержит сводную документацию по уже выполненным этапам миграции.

## Этап 1: Инициализация проекта

### Стек

- **Go 1.23+** — язык
- **Chi v5** — HTTP-роутер
- **pgx v5** — PostgreSQL-драйвер
- **golang-migrate** — миграции БД
- **squirrel** — SQL query builder
- **golang-jwt** — JWT-токены
- **validator/v10** — валидация данных
- **godotenv** — загрузка .env
- **testify** — тестирование

### Запуск

```bash
cp .env.example .env
make run
make build
make test
```

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| DATABASE_URL | Строка подключения к PostgreSQL | — |
| JWT_SECRET | Секретный ключ для JWT | — |
| JWT_EXPIRES_IN | Время жизни токена (Go duration) | 168h |
| PORT | Порт сервера | 3000 |

### Эндпоинт

| Метод | Путь | Описание |
|-------|------|----------|
| GET | / | Health check → `{"status": "ok"}` |

### Миграции

- 8 PostgreSQL enum-типов
- 10 таблиц: `organizations`, `users`, `contractors`, `jobs`, `job_dispatches`, `job_offers`, `assignments`, `communication_logs`, `contractor_blacklist`, `ratings`
- индексы на связи
- триггер `update_updated_at` для автообновления `updated_at`

## Этап 2: Модели, DTO и валидация

### Модели (`internal/models/models.go`)

Реализованы модели для всех таблиц БД, включая составные модели (`WithCounts`, `WithJob`, `WithContractor`) для join-запросов.

### Enum-типы

- UserRole: `ADMIN`, `MANAGER`, `OBSERVER`
- ContractorType: `INDIVIDUAL`, `IP`, `LLC`
- JobStatus: `DRAFT`, `ACTIVE`, `IN_PROGRESS`, `COMPLETED`, `CANCELLED`
- DispatchStatus: `PENDING`, `SENDING`, `COMPLETED`, `PAUSED`, `CANCELLED`
- OfferStatus: `SENT`, `ACCEPTED`, `DECLINED`, `EXPIRED`
- AssignmentStatus: `ASSIGNED`, `CONFIRMED`, `IN_WORK`, `COMPLETED`, `CANCELLED`
- CommunicationDirection: `OUTGOING`, `INCOMING`
- RatingAuthorType: `CUSTOMER`, `CONTRACTOR`

### DTO и валидация (`internal/dto/*`)

Используется `go-playground/validator/v10` с тегами:
- `required`
- `email`
- `min`, `max`
- `oneof`
- `omitempty`

Добавлены DTO для auth/jobs/contractors/users/organizations.

### Общие утилиты (`internal/dto/common.go`)

- `DecodeAndValidate`
- `WriteJSON`
- `WriteError`
- `WriteValidationError`

## Этап 3: Аутентификация (auth)

### Что реализовано

- `internal/repository/user.go`: `FindByEmail`, `FindByID`, `Create`
- `internal/repository/organization.go`: `Create`
- `internal/service/auth.go`:
  - `Register` (организация + admin в транзакции)
  - `Login`
  - генерация JWT
- `internal/handler/auth.go`
- `internal/middleware/auth.go` (JWT middleware)
- DI wiring в `cmd/server/main.go`

Пароли хешируются через `bcrypt` (cost 12).

### Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| POST | `/auth/register` | - | Регистрация организации + admin |
| POST | `/auth/login` | - | Вход |
| GET | `/auth/me` | JWT | Текущий пользователь |

### Тесты

- `internal/middleware/auth_test.go`
- `internal/router/router_test.go`
- `go test ./...` — PASS

## Этап 4: Пользователи и организации

### Users

- `internal/service/users.go`
- `internal/handler/users.go`
- в `internal/repository/user.go` добавлены scoped-методы:
  - `FindByIDAndOrganization`
  - `UpdateInOrganization`
  - `DeleteInOrganization`

Эндпоинты:
- `GET /users`
- `GET /users/{id}`
- `PATCH /users/{id}`
- `DELETE /users/{id}`

### Organizations

- `internal/service/organizations.go`
- `internal/handler/organizations.go`

Эндпоинты:
- `GET /organizations/{id}`
- `PATCH /organizations/{id}`

### Role-based доступ

- `internal/middleware/roles.go` (`RequireRoles`)
- `internal/middleware/roles_test.go`

Правила:
- JWT обязателен для users/org routes
- `PATCH /users/{id}` и `PATCH /organizations/{id}`: `ADMIN | MANAGER`
- `DELETE /users/{id}`: `ADMIN`

## Этап 5: Подрядчики (contractors)

### Что реализовано

- `internal/repository/contractor.go`
  - create/get/update
  - list с фильтрами/пагинацией
  - history (offers + assignments)
- `internal/repository/contractor_blacklist.go`
  - add/exist
- `internal/service/contractors.go`
- `internal/handler/contractors.go`

### Фильтры `GET /contractors`

- `region`
- `type` (`INDIVIDUAL|IP|LLC`)
- `isAvailable`
- `search` (`ILIKE` по имени/телефону)
- `page`, `limit`

### Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| POST | `/contractors` | JWT + role | Создание |
| GET | `/contractors` | JWT | Список с фильтрами |
| GET | `/contractors/{id}` | JWT | Карточка подрядчика |
| PATCH | `/contractors/{id}` | JWT + role | Обновление |
| POST | `/contractors/{id}/blacklist` | JWT + role | Добавление в ЧС |
| GET | `/contractors/{id}/history` | JWT | История взаимодействий |

`role` для mutate-операций: `ADMIN | MANAGER`.

### Тесты

- `internal/handler/contractors_test.go`
- `go test ./...` — PASS

## Этап 6: Работы (jobs)

### Что реализовано

- `internal/repository/job.go`
  - create draft
  - get by id (organization-scope)
  - list с filters/pagination + related counts
  - update только для `DRAFT`
  - status transitions
  - duplicate в новый `DRAFT`
- `internal/service/jobs.go`
- `internal/handler/jobs.go`

### Lifecycle

- publish: `DRAFT -> ACTIVE`
- cancel: `DRAFT|ACTIVE|IN_PROGRESS -> CANCELLED`
- complete: `ACTIVE|IN_PROGRESS -> COMPLETED`
- duplicate: копия job в новый `DRAFT`

### Эндпоинты

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

### Фильтры списка

- `status` (`DRAFT|ACTIVE|IN_PROGRESS|COMPLETED|CANCELLED`)
- `region`
- `page`, `limit`

В list-ответе есть счётчики:
- `dispatchCount`
- `offerCount`
- `assignmentCount`

### Тесты

- `internal/handler/jobs_test.go`
- `go test ./...` — PASS
