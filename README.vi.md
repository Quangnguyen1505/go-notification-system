**Ngôn ngữ:** [English](README.md) | Tiếng Việt

# go-notification-system

Hệ thống Notification theo hướng **enterprise**: hỗ trợ đa kênh (Email/SMS/Push/In-app), ưu tiên **scalability**, **high availability**, và **low latency**. Repo này hiện đang ở giai đoạn **khởi tạo kiến trúc + khung code** (scaffolding) theo mô hình microservices.

![Architecture](docs/Architechture-design.png)

## Mục tiêu

- Xử lý khối lượng lớn thông báo (đỉnh tải có thể tăng đột biến theo campaign/sự kiện).
- Tách rời producer/consumer bằng message queue để giảm coupling, tăng khả năng scale theo chiều ngang.
- Có retry policy rõ ràng và **Dead Letter Queue (DLQ)** cho các thông báo thất bại vĩnh viễn.
- Theo dõi trạng thái & analytics để đo hiệu quả gửi thông báo.

## Kiến trúc tổng quan (theo sơ đồ trong `docs/`)

### Các service chính

- **API Gateway (REST)**: entrypoint cho client; chịu trách nhiệm auth/rate limit/routing.
- **gRPC Gateway / Reverse Proxy**: chuyển REST → gRPC và route đến các backend service.
- **Notification Service**: trung tâm điều phối (format nội dung, gọi template/user prefs, enqueue job vào queue).
- **User Preferences Service**: quản lý lựa chọn kênh, quiet hours, opt-in/opt-out theo user/tenant.
- **Scheduler**: quản lý thông báo schedule/cron/recurring; enqueue đúng thời điểm.
- **Analytics/Reporting Service**: thu thập event (sent/delivered/opened/clicked) và thống kê.

### Queue + Worker + Retry

- **Notification Queue**: hàng đợi nhận các job gửi thông báo.
- **Notification Workers**: subscribe queue, gửi ra các provider (SMTP/SMS Gateway/FCM...).
- **Retry Service + Retry Queue**: nhận failure từ worker, áp dụng retry strategy (ví dụ exponential backoff).
- **Dead Letter Queue (DLQ)**: chứa job đã vượt retry limit để điều tra/khắc phục.

### Storage/Cache (định hướng)

- **DB (SQL)**: lưu notification metadata, trạng thái, user preferences, lịch schedule.
- **Cache (Redis)**: cache preference, template, status gần đây; giảm DB read.

## Functional requirements (tóm tắt)

- Multi-channel notifications: Email/SMS/Push/In-app.
- User preferences: kênh, tần suất, khung giờ yên lặng (quiet hours), theo user/tenant.
- Scheduling & prioritization: hẹn giờ, ưu tiên theo mức độ quan trọng.
- Template management: template động + placeholders + versioning.
- Multi-tenancy: tách dữ liệu/cấu hình theo tenant.
- Batch: gửi hàng loạt cho campaign.
- Retry mechanism: retry có policy cấu hình được + DLQ.
- Analytics/reporting: thống kê delivered/opened/clicked, báo cáo.

## Non-functional requirements (tóm tắt)

- Scalability: scale ngang theo traffic (worker/service).
- High availability: tránh single point of failure, mục tiêu uptime cao.
- Low latency: ưu tiên xử lý nhanh cho message ưu tiên cao.
- Fault tolerance: chịu lỗi từng thành phần (provider down, network issue...).
- Security & compliance: mã hóa in-transit/at-rest, audit log, hướng GDPR.
- Rate limiting: theo user/tenant/global.

## Capacity planning (tham chiếu giả định)

Các con số trong bài phân tích là **giả định** để làm capacity planning (ví dụ 200M/day, peak 10M/min). Khi triển khai thực tế cần hiệu chỉnh theo sản phẩm và SLA.

## Cấu trúc thư mục

Repo được tổ chức theo hướng tách **entrypoints** (`cmd/`), **core domain** (`internal/`), và **shared packages** (`pkg/`).

```
.
├── build/
├── cmd/
│   ├── cli/
│   │   └── makefile/
│   │       └── notification/
│   │           └── main.mk     # Hỗ trợ migration DB cho notification
│   ├── notification/
│   │   ├── config/
│   │   │   └── config.go       # Loader config cho service notification
│   │   ├── config.yml          # File config mẫu
│   │   └── main.go             # Entry gRPC
│   ├── printdsn/
│   │   └── main.go
│   ├── proxy/
│   │   ├── config/
│   │   │   └── config.go       # Loader config cho proxy
│   │   ├── config.yml          # File config mẫu
│   │   └── main.go             # HTTP reverse proxy + grpc-gateway
│   └── user/
├── db/
│   └── migrations/             # Migrations cho PostgreSQL
├── docker/
│   └── Dockerfile-proxy
├── docs/
│   └── cli/
│       └── cli-gen.md
├── global/
│   └── noti/
│       └── global.go           # State bootstrap dùng chung cho process
├── internal/
│   └── notification/
│       ├── app/
│       │   ├── app.go          # Composition root state
│       │   ├── wire.go         # Khai báo Wire injector
│       │   ├── wire_gen.go     # Wire generated code
│       │   └── router/
│       │       └── notification_grpc_server.go
│       ├── domain/
│       │   ├── interfaces.go   # Domain ports
│       │   └── models.go       # Domain model(s)
│       ├── infras/
│       │   ├── postgresql/
│       │   │   ├── gen/        # sqlc generated code
│       │   │   └── query/      # file SQL đầu vào cho sqlc
│       │   └── repo/
│       │       └── notification_postgres.go
│       └── usecases/
│           └── notification/
│               ├── interfaces.go
│               └── service.go
├── pkg/
│   ├── config/
│   ├── logger/
│   ├── postgres/
│   ├── rabbitmq/
│   └── utils/
├── proto/
│   ├── buf.yaml
│   ├── common.proto
│   ├── notification.proto
│   └── gen/
├── rests/
│   └── client.http
├── scripts/
├── storages/
│   └── logger/
│       ├── notification/
│       └── proxy/
├── third_party/
│   └── OpenAPI/
│       ├── common.swagger.json
│       └── notification.swagger.json
├── tools/
│   └── tools.go
├── buf.gen.yaml
├── buf.work.yaml
├── docker-compose-core.yml
├── go.mod
├── Makefile
├── README.md
├── README.vi.md
└── sqlc.yaml
```

## Tổ chức theo DDD / bounded-context

Repo theo hướng DDD “thực dụng” theo từng bounded context (hiện có **notification**):

- **`internal/<context>/domain`**: entity/value object và các *port* (interface) mà use-case phụ thuộc.
- **`internal/<context>/usecases`**: application service (điều phối domain + port).
- **`internal/<context>/infras`**: adapter/integration (Postgres/sqlc, queue, provider...).
- **`internal/<context>/app`**: composition root (Wire injector, router/server bootstrap).

Nguyên tắc: domain/usecase **không import** infrastructure; infrastructure mới implement interface trong domain.

## Chạy local (phần đã làm được)

Một số service còn đang làm, nhưng bounded context **notification** đã có các phần chạy được:

- Kết nối PostgreSQL qua `pkg/postgres`
- Migration bằng Goose trong `db/migrations`
- Sinh code `sqlc` trong `internal/notification/infras/postgresql/gen`
- Repo Postgres trong `internal/notification/infras/repo`

### Chạy dependency

- Chạy PostgreSQL: `docker compose -f docker-compose-core.yml up -d postgres`

### Chạy migration (notification DB)

Repo có Makefile theo bounded-context, tự đọc `cmd/notification/config.yml` để tạo `GOOSE_DBSTRING` (không bắt buộc cài `yq`).

- Up migrations: `make -f cmd/cli/makefile/notification/main.mk upGoose`
- Reset DB (chỉ dev): `make -f cmd/cli/makefile/notification/main.mk resetGoose`

### Sinh code sqlc

- `sqlc generate`

### Chạy service

- Notification service: `go run ./cmd/notification`
- Proxy (grpc-gateway): `go run ./cmd/proxy`

## Ghi chú

- Hình kiến trúc nằm ở `docs/Architechture-design.png`.
- `cmd/notification` và `internal/notification/*` đang phát triển; API/contracts có thể thay đổi.
