# Этап 1: Инициализация проекта

## Стек

- **Go 1.23+** — язык
- **Chi v5** — HTTP-роутер
- **pgx v5** — PostgreSQL-драйвер
- **golang-migrate** — миграции БД
- **squirrel** — SQL query builder
- **golang-jwt** — JWT-токены
- **validator/v10** — валидация данных
- **godotenv** — загрузка .env
- **testify** — тестирование

## Запуск

```bash
# Скопировать переменные окружения
cp .env.example .env

# Запуск (миграции применяются автоматически)
make run

# Сборка бинарника
make build

# Тесты
make test
```

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| DATABASE_URL | Строка подключения к PostgreSQL | — |
| JWT_SECRET | Секретный ключ для JWT | — |
| JWT_EXPIRES_IN | Время жизни токена (Go duration) | 168h |
| PORT | Порт сервера | 3000 |

## Эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | / | Health check → `{"status": "ok"}` |

## Миграции

Миграции встроены в бинарник через `embed.FS` и применяются автоматически при запуске.

Схема БД включает:
- 8 PostgreSQL enum типов
- 10 таблиц: organizations, users, contractors, jobs, job_dispatches, job_offers, assignments, communication_logs, contractor_blacklist, ratings
- Индексы на foreign keys
- Триггер `update_updated_at` для автообновления `updated_at`

## Структура

```
go-backend/
├── cmd/server/main.go          # Точка входа
├── internal/
│   ├── config/config.go        # Загрузка конфигурации
│   ├── database/
│   │   ├── db.go               # Пул подключений pgx
│   │   ├── migrate.go          # Запуск миграций
│   │   └── migrations/         # SQL-файлы миграций
│   └── router/router.go        # Chi-роутер, CORS, health-check
├── .env / .env.example
├── Makefile
└── go.mod / go.sum
```
