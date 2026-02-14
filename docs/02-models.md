# Этап 2: Модели, DTO и валидация

## Модели (internal/models/models.go)

Все Go-структуры соответствуют таблицам БД из Prisma-схемы.

### Enum-типы

| Тип | Значения |
|-----|----------|
| UserRole | ADMIN, MANAGER, OBSERVER |
| ContractorType | INDIVIDUAL, IP, LLC |
| JobStatus | DRAFT, ACTIVE, IN_PROGRESS, COMPLETED, CANCELLED |
| DispatchStatus | PENDING, SENDING, COMPLETED, PAUSED, CANCELLED |
| OfferStatus | SENT, ACCEPTED, DECLINED, EXPIRED |
| AssignmentStatus | ASSIGNED, CONFIRMED, IN_WORK, COMPLETED, CANCELLED |
| CommunicationDirection | OUTGOING, INCOMING |
| RatingAuthorType | CUSTOMER, CONTRACTOR |

### Модели

- **Organization** — организация (id, name, inn, contactEmail, contactPhone)
- **User** — пользователь (id, email, passwordHash, name, role, organizationId)
- **Contractor** — подрядчик (id, name, phone, email, type, regions[], equipmentTypes[])
- **Job** — работа (id, title, description, region, equipmentType, status, price)
- **JobDispatch** — рассылка по работе
- **JobOffer** — предложение подрядчику
- **Assignment** — назначение подрядчика на работу
- **CommunicationLog** — лог коммуникации
- **ContractorBlacklist** — чёрный список
- **Rating** — рейтинг/отзыв

Составные модели (`WithCounts`, `WithJob`, `WithContractor`) используются для join-запросов.

## DTO (internal/dto/)

### Валидация

Используется `go-playground/validator/v10`. Теги валидации:
- `required` — обязательное поле
- `email` — формат email
- `min=N` — минимальная длина строки
- `max=N` — максимальная длина строки
- `oneof=A B C` — значение из списка
- `omitempty` — пропускать валидацию если пусто

### Auth DTOs

| DTO | Поля | Валидация |
|-----|------|-----------|
| RegisterRequest | organizationName, inn?, contactEmail?, contactPhone?, adminName, adminEmail, adminPassword | required, min=2/8, email |
| LoginRequest | email, password | required, email, min=8 |

### Job DTOs

| DTO | Поля | Валидация |
|-----|------|-----------|
| CreateJobRequest | title, description?, region, equipmentType?, volume?, deadline?, price?, conditions? | required title min=3, required region |
| UpdateJobRequest | все поля optional | min=3 для title |
| QueryJobsRequest | status?, region?, page, limit | oneof для status, min=1 для page/limit |

### Contractor DTOs

| DTO | Поля | Валидация |
|-----|------|-----------|
| CreateContractorRequest | name, phone, email?, type, regions[], equipmentTypes[], priceExpectations?, isAvailable?, consentGiven | required, oneof для type, email |
| UpdateContractorRequest | все поля optional | email, oneof для type |
| QueryContractorsRequest | region?, type?, isAvailable?, search?, page, limit | oneof для type |

## Утилиты (internal/dto/common.go)

- `DecodeAndValidate(r, dst)` — декодирует JSON body и валидирует
- `WriteJSON(w, status, data)` — пишет JSON-ответ
- `WriteError(w, status, message)` — пишет ошибку
- `WriteValidationError(w, err)` — пишет ошибку валидации с деталями по полям
