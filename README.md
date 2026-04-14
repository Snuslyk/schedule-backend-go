# Schedule Backend

Backend for an educational scheduling system built with Go.

Frontend repository:  
https://github.com/Eltarchik/schedule

---

## Tech Stack

- Go
- Chi (`github.com/go-chi/chi/v5`)
- PostgreSQL + sqlx
- Redis
- JWT (access & refresh)
- S3 (Yandex Cloud Object Storage via AWS SDK)

---

## Installation

```bash
go mod download
```

---

## Environment

Create a `.env` file in the project root:

```env
PORT=4242
ENV=development
COOKIE_DOMAIN=localhost

DATABASE_URL=postgresql://postgres:password@localhost:5432/schedule?sslmode=disable
REDIS_URL=localhost:6379
REDIS_PASSWORD=

JWT_SECRET=your_jwt_secret
JWT_ACCESS_TOKEN_TTL=2h
JWT_REFRESH_TOKEN_TTL=168h

S3_BUCKET=schedule-test-bucket
S3_ENDPOINT=https://storage.yandexcloud.net
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_S3_REGION=ru-central1-a
S3_USE_PATH_STYLE=false

CORS_ORIGIN=http://localhost:3000
```

---

## Running

### Docker Compose

```bash
docker compose up -d
```

### Locally

```bash
go run ./cmd/api
```

---

## API Documentation

Interactive Scalar API reference is available at:
```
http://localhost:4242/reference
```
