# Platform Go Backend

Go backend for platform migration from NestJS/Prisma.

## Stack

- Go 1.25
- Chi router
- PostgreSQL + pgx
- golang-migrate (embedded migrations)
- JWT auth
- validator/v10

## Quick Start

```bash
cp .env.example .env
make run
```

Server runs on `PORT` (default `3000`).

## Main Commands

```bash
make run
make build
make test
make test-cover
make lint
make migrate-up
make migrate-down
```

## Docs

- Implementation progress (stages 1-14): `docs/implementation.md`
- Swagger JSON: `GET /api/docs/swagger.json`
- Docs page: `GET /api/docs/`

## API Areas

- Auth
- Users / Organizations
- Contractors
- Jobs
- Matching
- Responses
- Assignments
- Ratings
- Compliance
- Notifications (service module)
- Analytics

## Notes

- All protected routes use Bearer JWT.
- Mutating routes are role-guarded (`ADMIN`, `MANAGER`) where applicable.
