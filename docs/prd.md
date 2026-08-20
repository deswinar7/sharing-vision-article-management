# PRD — Article Management Fullstack MVP

**Project Type:** Technical Test / Portfolio-ready MVP
**Backend:** Go + Gin
**Database:** MySQL
**Frontend:** React + TypeScript + Vite
**Architecture:** Single repository, modular monolith
**Target:** MVP / Take-home technical assessment

---

# 1. Product Overview

Membangun aplikasi Article Management sederhana yang memungkinkan user:

* melihat daftar artikel berdasarkan status;
* membuat artikel;
* mengubah artikel;
* mem-publish artikel;
* menyimpan artikel sebagai draft;
* memindahkan artikel ke trash;
* menghapus artikel;
* melihat preview artikel yang sudah published;
* melakukan pagination.

Aplikasi terdiri dari:

```text
React Web App
      │
      │ REST API / JSON
      ▼
Go Gin API
      │
      ▼
MySQL
```

Tidak menggunakan microservices terpisah karena scope terlalu kecil untuk membenarkan kompleksitas tersebut.

---

# 2. Goals

MVP harus:

1. Memenuhi seluruh requirement Backend Sharing Vision.
2. Memenuhi seluruh requirement Frontend Sharing Vision.
3. Memiliki code structure yang clean dan mudah dievaluasi.
4. Memiliki validation yang konsisten.
5. Memiliki migration database.
6. Memiliki Postman Collection.
7. Bisa dijalankan secara lokal dengan mudah.
8. Bisa di-deploy menggunakan free-tier hosting.
9. Memiliki README yang cukup sehingga reviewer bisa menjalankan project tanpa perlu bertanya.

---

# 3. Non-Goals

Tidak perlu:

* authentication;
* authorization;
* user management;
* image upload;
* rich text editor kompleks;
* Elasticsearch;
* Redis;
* WebSocket;
* message queue;
* microservices;
* Kubernetes;
* event-driven architecture;
* CQRS;
* GraphQL;
* SSR;
* distributed tracing.

Semua hal tersebut tidak memberikan manfaat berarti untuk scope technical test ini.

---

# 4. Recommended Technology Stack

## Backend

| Area           | Technology              |
| -------------- | ----------------------- |
| Language       | Go                      |
| HTTP Framework | Gin                     |
| Database       | MySQL                   |
| DB Driver      | go-sql-driver/mysql     |
| SQL Helper     | sqlx                    |
| Migration      | golang-migrate          |
| Validation     | go-playground/validator |
| Environment    | godotenv                |
| Testing        | Go testing + httptest   |
| Logging        | standard `log/slog`     |
| API            | REST JSON               |

### Kenapa sqlx daripada ORM?

Untuk technical assessment sederhana, `sqlx` memberikan keseimbangan yang baik:

* SQL tetap eksplisit;
* reviewer dapat melihat kemampuan SQL;
* abstraction sangat sedikit;
* scanning database tetap nyaman;
* tidak membawa kompleksitas ORM yang sebenarnya tidak dibutuhkan.

---

# 5. Frontend Stack

| Area         | Technology                     |
| ------------ | ------------------------------ |
| Language     | TypeScript                     |
| UI           | React                          |
| Build Tool   | Vite                           |
| Routing      | React Router                   |
| Server State | TanStack Query                 |
| Form         | React Hook Form                |
| Validation   | Zod                            |
| Styling      | Tailwind CSS                   |
| Icons        | Lucide React                   |
| HTTP         | native Fetch wrapper           |
| Testing      | Vitest + React Testing Library |

### Kenapa React + Vite?

Aplikasi ini merupakan dashboard SPA dan tidak memerlukan SSR.

Vite + React:

* lebih sederhana daripada Next.js;
* build menjadi static assets;
* mudah di-host gratis;
* cocok untuk REST API;
* startup dan development cepat;
* tidak memperkenalkan server runtime frontend yang sebenarnya tidak diperlukan.

---

# 6. Repository Strategy

Gunakan **single repository**, tetapi jangan menggunakan Turborepo/Nx karena tidak memberikan manfaat berarti untuk hanya satu Go backend dan satu React frontend.

```text
article-management/
│
├── backend/
├── frontend/
├── postman/
├── docs/
│
├── docker-compose.yml
├── .gitignore
├── README.md
└── Makefile
```

Keuntungan:

* reviewer cukup clone satu repository;
* frontend/backend versioning sinkron;
* deployment tetap bisa dilakukan secara terpisah;
* local development lebih mudah.

---

# 7. Backend Architecture

Gunakan architecture sederhana:

```text
HTTP Request
     ↓
Handler
     ↓
Service
     ↓
Repository
     ↓
MySQL
```

Tidak perlu Clean Architecture penuh dengan banyak abstraction.

Struktur:

```text
backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── mysql.go
│   │
│   ├── article/
│   │   ├── model.go
│   │   ├── dto.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── routes.go
│   │
│   └── middleware/
│       ├── cors.go
│       └── recovery.go
│
├── migrations/
│   ├── 000001_create_posts.up.sql
│   └── 000001_create_posts.down.sql
│
├── tests/
│
├── .env.example
├── go.mod
├── go.sum
└── Dockerfile
```

---

# 8. Database Design

Database:

```text
article
```

Table:

```text
posts
```

Schema:

```sql
CREATE TABLE posts (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,
    status VARCHAR(100) NOT NULL,

    PRIMARY KEY (id),
    INDEX idx_posts_status (status),
    INDEX idx_posts_created_date (created_date)
);
```

Status aplikasi:

```text
publish
draft
thrash
```

Walaupun `trash` secara bahasa Inggris lebih tepat, gunakan `thrash` secara internal karena technical specification menggunakan value tersebut.

UI dapat tetap menampilkan:

```text
Published
Drafts
Trashed
```

---

# 9. Data Model

```go
type Article struct {
    ID          uint64    `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    Category    string    `json:"category"`
    Status      string    `json:"status"`
    CreatedDate time.Time `json:"created_date"`
    UpdatedDate time.Time `json:"updated_date"`
}
```

---

# 10. Input DTO

```text
CreateArticleRequest
- title
- content
- category
- status

UpdateArticleRequest
- title
- content
- category
- status
```

Tidak menerima:

```text
id
created_date
updated_date
```

dari client.

Field tersebut server-controlled.

---

# 11. Validation Rules

## Title

```text
required
minimum 20 characters
maximum 200 characters
```

## Content

```text
required
minimum 200 characters
```

## Category

```text
required
minimum 3 characters
maximum 100 characters
```

## Status

Harus salah satu:

```text
publish
draft
thrash
```

Validation diterapkan pada:

```text
POST /article/
PUT/PATCH /article/:id
```

Frontend validation membantu UX, tetapi backend tetap menjadi source of truth.

---

# 12. REST API Contract

## 12.1 Create Article

```http
POST /article/
```

Request:

```json
{
  "title": "Example article title with sufficient characters",
  "content": "Article content with minimum length of two hundred characters...",
  "category": "Technology",
  "status": "draft"
}
```

Recommended success:

```http
201 Created
```

Response:

```json
{
  "data": {
    "id": 1,
    "title": "...",
    "content": "...",
    "category": "Technology",
    "status": "draft",
    "created_date": "...",
    "updated_date": "..."
  }
}
```

---

# 13. Get Articles

Required endpoint:

```http
GET /article/:limit/:offset
```

Example:

```http
GET /article/10/0
```

Optional query parameter:

```http
GET /article/10/0?status=publish
```

Supported:

```text
publish
draft
thrash
```

Query `status` merupakan extension kecil terhadap requirement asli agar tab frontend dan pagination dapat dilakukan di server secara benar.

Response:

```json
[
  {
    "id": 1,
    "title": "...",
    "content": "...",
    "category": "...",
    "status": "publish",
    "created_date": "...",
    "updated_date": "..."
  }
]
```

Gunakan:

```sql
ORDER BY created_date DESC
LIMIT ?
OFFSET ?
```

Pagination harus deterministic.

---

# 14. Get Article Detail

```http
GET /article/:id
```

Success:

```http
200 OK
```

Not found:

```http
404 Not Found
```

---

# 15. Update Article

Gunakan:

```http
PUT /article/:id
```

atau:

```http
PATCH /article/:id
```

Untuk implementation MVP, pilih **PUT** apabila form mengirim seluruh field.

Request:

```json
{
  "title": "...",
  "content": "...",
  "category": "...",
  "status": "publish"
}
```

Success:

```http
200 OK
```

---

# 16. Delete Article

```http
DELETE /article/:id
```

Endpoint ini melakukan permanent deletion.

Namun tombol trash pada halaman Published/Drafts **tidak menggunakan DELETE**.

Flow:

```text
Published / Draft
       ↓
click trash
       ↓
PATCH/PUT Article
status = thrash
       ↓
Article muncul pada Trashed
```

Hal ini diperlukan karena frontend requirement mengatakan article harus berpindah ke tab Trash.

`DELETE` dapat digunakan sebagai permanent delete dari Trashed jika diperlukan.

---

# 17. Health Check

Tambahkan endpoint:

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

Digunakan untuk:

* deployment health check;
* debugging;
* memastikan API hidup.

---

# 18. Standard HTTP Status Codes

Gunakan:

```text
200 OK
201 Created
400 Bad Request
404 Not Found
422 Unprocessable Entity
500 Internal Server Error
```

Tidak mengembalikan `200` untuk semua kondisi.

---

# 19. Standard Error Format

Gunakan satu bentuk error secara konsisten.

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "fields": {
      "title": "title must contain at least 20 characters"
    }
  }
}
```

Contoh:

```text
VALIDATION_ERROR
ARTICLE_NOT_FOUND
INVALID_STATUS
INTERNAL_ERROR
```

Jangan leak raw database error kepada frontend.

---

# 20. Frontend Information Architecture

Main navigation:

```text
Posts
 ├── All Posts
 ├── Add New
 └── Preview
```

Routes:

```text
/
→ redirect /posts

/posts
/posts/new
/posts/:id/edit
/preview
```

---

# 21. All Posts Page

Route:

```text
/posts
```

Tabs:

```text
Published
Drafts
Trashed
```

Setiap tab menampilkan:

| Title         | Category   | Actions      |
| ------------- | ---------- | ------------ |
| Article title | Technology | Edit / Trash |

Published:

```text
status=publish
```

Draft:

```text
status=draft
```

Trashed:

```text
status=thrash
```

---

# 22. Article Actions

## Edit

Icon:

```text
Pencil
```

Navigate:

```text
/posts/:id/edit
```

## Trash

Icon:

```text
Trash2
```

Behavior:

```text
update article:
status = thrash
```

Setelah berhasil:

```text
invalidate TanStack Query cache
→ table otomatis refresh
```

Tidak perlu melakukan manual global-state synchronization.

---

# 23. Add New Page

Route:

```text
/posts/new
```

Fields:

```text
Title
Content
Category
```

Actions:

```text
Save as Draft
Publish
```

### Publish

Kirim:

```json
{
  "...": "...",
  "status": "publish"
}
```

### Draft

Kirim:

```json
{
  "...": "...",
  "status": "draft"
}
```

---

# 24. Edit Page

Route:

```text
/posts/:id/edit
```

Load:

```http
GET /article/:id
```

Pre-fill:

```text
Title
Content
Category
```

Buttons:

```text
Save as Draft
Publish
```

Setelah update berhasil:

```text
invalidate article query
invalidate article-list query
navigate /posts
```

---

# 25. Preview Page

Route:

```text
/preview
```

Hanya menampilkan:

```text
status=publish
```

Setiap article minimal menampilkan:

```text
Title
Category
Content
Published/created date
```

Pagination:

```text
Previous
1
2
3
Next
```

MVP page size:

```text
5 articles
```

Backend menggunakan:

```text
limit = 5
offset = (page - 1) * 5
status = publish
```

---

# 26. Frontend Structure

```text
frontend/
├── src/
│   ├── app/
│   │   ├── router.tsx
│   │   └── providers.tsx
│   │
│   ├── components/
│   │   ├── ui/
│   │   ├── layout/
│   │   └── feedback/
│   │
│   ├── features/
│   │   └── articles/
│   │       ├── api/
│   │       │   └── article.api.ts
│   │       ├── components/
│   │       │   ├── ArticleForm.tsx
│   │       │   ├── ArticleTable.tsx
│   │       │   └── ArticlePreview.tsx
│   │       ├── hooks/
│   │       ├── schemas/
│   │       │   └── article.schema.ts
│   │       ├── types/
│   │       │   └── article.ts
│   │       └── pages/
│   │           ├── AllPostsPage.tsx
│   │           ├── AddPostPage.tsx
│   │           ├── EditPostPage.tsx
│   │           └── PreviewPage.tsx
│   │
│   ├── lib/
│   │   └── api.ts
│   │
│   ├── main.tsx
│   └── index.css
│
├── .env.example
├── package.json
├── tsconfig.json
└── vite.config.ts
```

Gunakan feature-based structure tetapi tetap sederhana.

---

# 27. State Management

Tidak perlu:

```text
Redux
Zustand
MobX
Context global state
```

Gunakan TanStack Query untuk server state:

```text
useQuery
useMutation
invalidateQueries
```

Local UI state tetap dengan:

```text
useState
```

---

# 28. Query Keys

Contoh:

```typescript
articleKeys.all

articleKeys.list({
  status,
  limit,
  offset,
})

articleKeys.detail(id)
```

Query key harus centralized agar invalidation predictable.

---

# 29. Frontend Form Validation

Gunakan:

```text
React Hook Form
+
Zod
```

Schema frontend harus mengikuti backend:

```text
title:
min 20
max 200

content:
min 200

category:
min 3
max 100

status:
publish | draft | thrash
```

Backend validation tetap wajib.

---

# 30. Loading / Error / Empty States

Semua data page wajib memiliki:

### Loading

```text
Loading...
```

atau skeleton sederhana.

### Empty

```text
No published articles found.
```

### Error

```text
Failed to load articles.
Retry
```

Tidak boleh menghasilkan blank screen.

---

# 31. UX Requirements

Form:

* disable button ketika request sedang berlangsung;
* tampilkan validation dekat field;
* hindari double submit;
* tampilkan feedback setelah berhasil;
* konfirmasi sebelum memindahkan article ke trash;
* tombol Back/Cancel pada Add/Edit.

Tidak perlu animation kompleks.

---

# 32. Responsive Design

Priority:

```text
Desktop
Tablet
Mobile
```

Dashboard tetap harus usable pada mobile.

Table dapat:

* horizontal scroll;
* atau berubah menjadi card sederhana pada layar kecil.

Jangan menghabiskan banyak waktu untuk visual polish yang tidak mempengaruhi core requirement.

---

# 33. Local Development

Gunakan Docker hanya untuk MySQL.

```text
docker compose up -d mysql
```

Backend dan frontend tetap bisa dijalankan native:

```text
Backend:
go run ./cmd/api

Frontend:
npm run dev
```

Hal ini membuat feedback loop development lebih cepat.

---

# 34. Environment Variables

Backend:

```env
APP_ENV=development
PORT=8080

DB_HOST=localhost
DB_PORT=3306
DB_NAME=article
DB_USER=article_user
DB_PASSWORD=article_password

CORS_ALLOWED_ORIGINS=http://localhost:5173
```

Frontend:

```env
VITE_API_BASE_URL=http://localhost:8080
```

Tidak commit `.env`.

Commit:

```text
.env.example
```

---

# 35. Database Migration

Wajib memiliki:

```text
up migration
down migration
```

Commands:

```text
make migrate-up
make migrate-down
```

Database schema tidak boleh hanya mengandalkan pembuatan manual.

---

# 36. Postman Collection

Folder:

```text
postman/
└── SharingVisionArticleAPI.postman_collection.json
```

Collection minimal berisi:

```text
Health Check

Article
├── Create Draft Article
├── Create Published Article
├── Get Articles
├── Get Published Articles
├── Get Article Detail
├── Update Article
├── Move Article To Trash
└── Delete Article
```

Gunakan variable:

```text
{{base_url}}
{{article_id}}
```

---

# 37. Testing Scope

Tidak perlu mengejar 100% coverage.

Prioritaskan critical behavior.

## Backend

Test:

```text
create article validation
invalid status
title < 20
content < 200
category < 3
get non-existing article
pagination validation
successful creation
successful update
```

Gunakan:

```text
testing
httptest
```

## Frontend

Minimal:

```text
ArticleForm validation
status mutation
All Posts loading/error state
```

Menggunakan:

```text
Vitest
React Testing Library
```

---

# 38. Code Quality

Backend:

```bash
gofmt
go vet
go test ./...
```

Frontend:

```bash
npm run lint
npm run typecheck
npm run test
npm run build
```

Tidak boleh ada:

```text
console.log debugging
unused code
hardcoded API URLs
committed credentials
```

---

# 39. GitHub Actions

Satu workflow:

```text
.github/workflows/ci.yml
```

Jobs:

```text
backend-check
frontend-check
```

Backend:

```text
go vet
go test
go build
```

Frontend:

```text
npm ci
npm run lint
npm run typecheck
npm run test
npm run build
```

Tidak perlu pipeline deployment kompleks.

---

# 40. Deployment Architecture

Recommended free deployment:

```text
                         ┌─────────────────────┐
                         │ Cloudflare Pages    │
User ───────────────────▶│ React + Vite        │
                         └──────────┬──────────┘
                                    │
                                    │ HTTPS REST API
                                    ▼
                         ┌─────────────────────┐
                         │ Render              │
                         │ Go + Gin            │
                         └──────────┬──────────┘
                                    │
                                    │ TLS MySQL
                                    ▼
                         ┌─────────────────────┐
                         │ Aiven               │
                         │ MySQL Free          │
                         └─────────────────────┘
```

---

# 41. Deployment Target

## Frontend

Cloudflare Pages.

Build:

```bash
npm run build
```

Output:

```text
dist
```

Environment:

```env
VITE_API_BASE_URL=https://<backend>.onrender.com
```

---

## Backend

Render Web Service.

Build command:

```bash
go build -tags netgo -ldflags="-s -w" -o app ./cmd/api
```

Start:

```bash
./app
```

Application wajib listen:

```go
0.0.0.0:$PORT
```

---

## Database

Aiven MySQL Free.

Gunakan credentials melalui environment variables.

Production connection harus menggunakan TLS sesuai connection configuration yang diberikan Aiven.

---

# 42. CORS

Development:

```text
http://localhost:5173
```

Production:

```text
https://<cloudflare-pages-domain>
```

Jangan menggunakan:

```text
Access-Control-Allow-Origin: *
```

secara permanen untuk production configuration.

---

# 43. README Requirements

README merupakan bagian penting dari technical test.

Wajib menjelaskan:

```text
Project overview
Architecture
Tech stack
Prerequisites
Installation
Environment variables
Database setup
Migration
Running backend
Running frontend
Testing
API documentation
Postman
Deployment URLs
Design decisions
Known limitations
```

Tambahkan:

```text
Frontend URL:
Backend URL:
Health URL:
```

---

# 44. Design Decisions

Dokumentasikan secara singkat:

### Why React + Vite?

Karena aplikasi merupakan SPA dashboard dan tidak membutuhkan SSR.

### Why sqlx?

Karena CRUD sederhana lebih jelas dengan explicit SQL dan dependency lebih kecil.

### Why status `thrash`?

Untuk menjaga kompatibilitas dengan requirement technical test.

### Why trash ≠ DELETE?

Karena requirement frontend meminta article berpindah ke Trashed sehingga membutuhkan soft-state transition terlebih dahulu.

### Why no authentication?

Di luar scope requirement.

### Why no microservices?

Hanya terdapat satu bounded functionality dan satu database. Pemisahan microservices akan menambah kompleksitas tanpa benefit.

---

# 45. Security Basics

Walaupun MVP:

* prepared statements / parameterized queries;
* input validation;
* CORS whitelist;
* environment variables;
* jangan expose DB errors;
* body size reasonable;
* TLS pada production DB;
* credentials tidak masuk repository.

Gin recovery middleware digunakan untuk mencegah panic membuat application process mati.

---

# 46. Performance

Tidak perlu premature optimization.

Cukup:

```text
database connection pooling
status index
created_date index
pagination
TanStack Query cache
```

Tidak perlu Redis.

---

# 47. Suggested Implementation Order

### Phase 1 — Foundation

1. Initialize repository.
2. Backend Go module.
3. React Vite project.
4. Docker Compose MySQL.
5. Environment configuration.

### Phase 2 — Database

6. Migration.
7. Database connection.
8. Article model.
9. Repository.

### Phase 3 — API

10. Validation.
11. Service.
12. Handlers.
13. Routes.
14. Error handling.
15. CORS.
16. Health check.

### Phase 4 — Backend Verification

17. Tests.
18. Postman Collection.

### Phase 5 — Frontend

19. API client.
20. TanStack Query setup.
21. Routing.
22. App layout.
23. All Posts.
24. Add New.
25. Edit.
26. Trash behavior.
27. Preview.
28. Pagination.

### Phase 6 — Quality

29. Loading/error states.
30. Responsive UI.
31. Lint/typecheck/tests.
32. README.

### Phase 7 — Deployment

33. Aiven MySQL.
34. Render backend.
35. Cloudflare Pages frontend.
36. Configure CORS.
37. Smoke test production.

---

# 48. Definition of Done

Project dianggap selesai apabila:

* [ ] MySQL schema tersedia melalui migration.
* [ ] POST article berjalan.
* [ ] GET article list + pagination berjalan.
* [ ] GET article detail berjalan.
* [ ] UPDATE article berjalan.
* [ ] DELETE article berjalan.
* [ ] Semua backend validation bekerja.
* [ ] Invalid request menghasilkan error yang benar.
* [ ] Published tab bekerja.
* [ ] Drafts tab bekerja.
* [ ] Trashed tab bekerja.
* [ ] Add New bekerja.
* [ ] Edit bekerja.
* [ ] Publish bekerja.
* [ ] Save Draft bekerja.
* [ ] Move to Trash bekerja.
* [ ] Preview hanya menampilkan published posts.
* [ ] Preview memiliki pagination.
* [ ] Loading/error/empty states tersedia.
* [ ] Postman Collection tersedia.
* [ ] Backend tests lulus.
* [ ] Frontend checks lulus.
* [ ] README lengkap.
* [ ] Frontend berhasil di-deploy.
* [ ] Backend berhasil di-deploy.
* [ ] Backend tersambung ke hosted MySQL.
* [ ] Tidak ada credential di repository.

---

# 49. Final Deliverable

```text
GitHub Repository
│
├── Complete source code
├── Database migrations
├── Docker Compose
├── Postman Collection
├── Tests
├── GitHub Actions CI
├── README
│
├── Live Frontend URL
└── Live Backend URL
```

Target akhir bukan sistem enterprise.

Target akhirnya adalah **small, clean, complete, deployable fullstack application** yang menunjukkan fundamental software engineering dengan baik tanpa over-engineering.
