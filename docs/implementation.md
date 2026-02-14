# Реализация Этапов 1-14

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

## Этап 7: Подбор исполнителей (matching)

### Что реализовано

- `internal/repository/matching.go`
  - выбор доступных подрядчиков (`is_available = true`)
  - матч по региону (`region = ANY(regions)`)
  - матч по типу оборудования (`equipment_types @> ARRAY[equipmentType]`), если у job задан `equipmentType`
- `internal/service/matching.go`
  - поиск работы по `jobId` в scope организации
  - валидация UUID
  - возврат списка подходящих подрядчиков
- `internal/handler/matching.go`
  - endpoint подбора по job id

### Эндпоинт

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/matching/jobs/{jobId}/contractors` | JWT | Подходящие подрядчики для работы |

### Тесты

- `internal/repository/matching_test.go`
  - build фильтров с `equipmentType`
  - build фильтров без `equipmentType`
- `go test ./...` — PASS

## Этап 8: Отклики (responses)

### Что реализовано

- `internal/repository/job_offer.go`
  - поиск оффера по токену (`FindByToken`)
  - смена статуса оффера по токену (`UpdateStatusByToken`)
- `internal/service/responses.go`
  - `GetByToken`
  - `Accept`
  - `Decline`
  - проверка формата token (UUID)
  - ограничение на ответ только для офферов в статусе `SENT`
- `internal/handler/responses.go`
  - публичные обработчики без JWT

### Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/responses/{token}` | - | Детали предложения |
| POST | `/responses/{token}/accept` | - | Принять предложение |
| POST | `/responses/{token}/decline` | - | Отклонить предложение |

### Тесты

- `internal/service/responses_test.go`
  - проверка бизнес-правила respondable status (`SENT` only)
- `go test ./...` — PASS

## Этап 9: Назначения (assignments)

### Что реализовано

- `internal/dto/assignment.go`
  - `CreateAssignmentRequest`
  - `UpdateAssignmentStatusRequest`
- `internal/repository/assignment.go`
  - создание назначения
  - получение назначений по `jobId` в scope организации
  - обновление статуса в scope организации
- `internal/service/assignments.go`
  - проверки существования job и contractor
  - обработка конфликтов (`job_id + contractor_id` unique)
  - обновление статуса назначения
- `internal/handler/assignments.go`
  - endpoints для create/list/updateStatus

### Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| POST | `/assignments` | JWT + role | Создать назначение |
| GET | `/assignments/job/{jobId}` | JWT | Список назначений по работе |
| PATCH | `/assignments/{id}/status` | JWT + role | Обновить статус назначения |

`role` для mutate-операций: `ADMIN | MANAGER`.

### Тесты

- `internal/dto/assignment_test.go`
  - валидация DTO create/updateStatus
- `go test ./...` — PASS

## Этап 10: Рейтинги (ratings)

### Что реализовано

- `internal/dto/rating.go`
  - `CreateRatingRequest`
- `internal/repository/rating.go`
  - создание рейтинга
  - список рейтингов по подрядчику (scope организации)
  - список рейтингов по работе (scope организации)
- `internal/service/ratings.go`
  - проверка существования job/contractor
  - create/list по contractor/list по job
- `internal/handler/ratings.go`
  - create/list endpoints

### Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| POST | `/ratings` | JWT + role | Создать рейтинг |
| GET | `/ratings/contractor/{contractorId}` | JWT | Рейтинги подрядчика |
| GET | `/ratings/job/{jobId}` | JWT | Рейтинги по работе |

`role` для mutate-операций: `ADMIN | MANAGER`.

### Тесты

- `internal/dto/rating_test.go`
  - валидация DTO create
  - валидация диапазона score
- `go test ./...` — PASS

## Этап 11: Комплаенс (compliance)

### Что реализовано

- `internal/repository/contractor_blacklist.go`
  - добавлены методы:
    - `FindAllByOrganization`
    - `DeleteInOrganization`
- `internal/service/compliance.go`
  - `ListBlacklist`
  - `DeleteBlacklistEntry`
  - валидация UUID для `blacklist id`
- `internal/handler/compliance.go`
  - endpoints для получения и удаления записей blacklist

### Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/compliance/blacklist` | JWT | Список blacklist организации |
| DELETE | `/compliance/blacklist/{id}` | JWT + role | Удалить запись из blacklist |

`role` для delete: `ADMIN | MANAGER`.

### Тесты

- `go test ./...` — PASS

## Этап 12: Уведомления (notifications)

### Что реализовано

- `internal/repository/communication_log.go`
  - создание записи `communication_logs`
  - получение логов по подрядчику с пагинацией
- `internal/service/notifications.go`
  - `CreateLog`
  - `ListByContractor`
  - валидация contractor/job UUID
  - валидация `direction` (`OUTGOING`/`INCOMING`)
  - проверка существования подрядчика

### Примечание по API

На этом этапе отдельный HTTP-контроллер не добавлялся.
Модуль сделан как сервисный слой, который вызывается другими модулями.

### Тесты

- `internal/service/notifications_test.go`
  - валидация parser-а направления коммуникации
- `go test ./...` — PASS

## Этап 13: Аналитика (analytics)

### Что реализовано

- `internal/repository/analytics.go`
  - агрегированные метрики дашборда по организации:
    - `totalJobs`
    - `activeJobs` (`ACTIVE` + `IN_PROGRESS`)
    - `totalContractors` (distinct contractors по назначениям организации)
    - `totalAssignments`
- `internal/service/analytics.go`
  - `Dashboard`
- `internal/handler/analytics.go`
  - endpoint дашборда

### Эндпоинт

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/analytics/dashboard` | JWT | Метрики дашборда организации |

### Тесты

- `internal/repository/analytics_test.go`
  - проверка JSON-контракта ответа метрик
- `go test ./...` — PASS

## Этап 14: Swagger и финальная документация

### Что реализовано

- `internal/apidocs/spec.go`
  - минимальный OpenAPI JSON (`SwaggerJSON`)
  - HTML-страница документации (`SwaggerHTML`)
- `internal/handler/swagger.go`
  - `JSON`
  - `UI`
- интеграция в `internal/router/router.go`
  - `GET /api/docs/`
  - `GET /api/docs/swagger.json`
- интеграция в `cmd/server/main.go` (DI wiring)
- финальный `README.md` в корне проекта
- контейнерный контур проверки:
  - `Dockerfile`
  - `Dockerfile.verifier`
  - `docker-compose.yml`
  - `scripts/verify_all.sh`
  - `scripts/smoke_api.sh`
  - `make docker-verify`

### Эндпоинты

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/api/docs/` | - | HTML-страница документации |
| GET | `/api/docs/swagger.json` | - | OpenAPI JSON |

### Тесты

- `internal/router/router_test.go`
  - `TestSwaggerJSONEndpoint` (валидный JSON + поле `openapi`)
- `go test ./...` — PASS

### Автопроверка перед стартом работы

Одна команда:

```bash
make docker-verify
```

Выполняет:
- подъем PostgreSQL и API в Docker
- unit-тесты (`go test ./...`)
- smoke-прогон HTTP-ручек (public + protected + role-protected + responses)
- остановку и очистку контейнеров после завершения
