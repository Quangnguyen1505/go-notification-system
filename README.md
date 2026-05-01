**Language:** English | [Tiếng Việt](README.vi.md)

# go-notification-system

An **enterprise-style** notification system supporting multiple channels (Email/SMS/Push/In-app), prioritizing **scalability**, **high availability**, and **low latency**. This repository is currently in the **architecture + scaffolding** phase for a microservices-based design.

![Architecture](docs/Architechture-design.png)

## Goals

- Handle large notification volume (spiky peak loads during campaigns/incidents).
- Decouple producers/consumers via message queues for horizontal scaling.
- Provide a clear retry policy and a **Dead Letter Queue (DLQ)** for permanent failures.
- Track delivery status and analytics for operational insight.

## High-level architecture (based on the diagram in `docs/`)

### Core services

- **API Gateway (REST)**: client entrypoint; auth/rate limiting/routing.
- **gRPC Gateway / Reverse Proxy**: bridges REST → gRPC and routes to backend services.
- **Notification Service**: orchestration core (formatting, template/user prefs, enqueue jobs).
- **User Preferences Service**: per-user/per-tenant settings (channels, quiet hours, opt-in/out).
- **Scheduler**: scheduled/recurring notifications; enqueues at the correct time.
- **Analytics/Reporting Service**: collects events (sent/delivered/opened/clicked) and metrics.

### Queue + Worker + Retry

- **Notification Queue**: holds send jobs.
- **Notification Workers**: consume jobs and deliver via providers (SMTP/SMS gateway/FCM...).
- **Retry Service + Retry Queue**: applies retry strategies (e.g., exponential backoff).
- **Dead Letter Queue (DLQ)**: stores jobs that exceeded retry limits for investigation.

### Storage/Cache (direction)

- **DB (SQL)**: notification metadata/status, user preferences, schedules.
- **Cache (Redis)**: cache preferences/templates/recent statuses to reduce DB load.

## Functional requirements (summary)

- Multi-channel notifications: Email/SMS/Push/In-app.
- User preferences: channels, frequency, quiet hours; per user/tenant.
- Scheduling & prioritization: schedule delivery; prioritize by importance.
- Template management: dynamic templates with placeholders + versioning.
- Multi-tenancy: isolate data/config per tenant.
- Batch sending: bulk campaigns.
- Retry mechanism: configurable policy + DLQ.
- Analytics/reporting: delivery + engagement metrics.

## Non-functional requirements (summary)

- Scalability: horizontal scaling (services/workers).
- High availability: avoid single points of failure.
- Low latency: fast handling for high-priority messages.
- Fault tolerance: resilient to provider/network failures.
- Security & compliance: in-transit/at-rest encryption, audit logging, GDPR direction.
- Rate limiting: per user/tenant/global.

## Capacity planning (assumption-based)

The numbers referenced in the analysis are **assumptions** for capacity planning (e.g., 200M/day, peak 10M/min). In production, calibrate them against real traffic and SLA.

## Repository structure

The codebase separates **entrypoints** (`cmd/`), **core domain** (`internal/`), and **shared packages** (`pkg/`).

```
.
├── cmd/
│   ├── proxy/                 # Reverse proxy (grpc-gateway) - REST → gRPC entrypoint
│   │   ├── config/            # Proxy config loader
│   │   ├── config.yml         # Example config
│   │   └── main.go            # HTTP server + grpc-gateway mux
│   ├── notification/          # Notification service entrypoint (skeleton)
│   └── user/                  # User Preferences service entrypoint (skeleton)
│
├── internal/
│   └── notification/          # Notification bounded context (skeleton)
│       ├── app/               # App bootstrap / wiring
│       ├── domain/            # Entities/VOs/domain services
│       ├── usecases/          # Application use-cases
│       └── infras/            # DB/queue/provider integrations
│
├── pkg/
│   ├── config/                # Shared config structs
│   ├── logger/                # Logging adapters (logrus ↔ slog)
│   ├── postgres/              # PostgreSQL helpers (direction)
│   ├── rabbitmq/              # RabbitMQ helpers (direction)
│   └── utils/                 # Shared utilities
│
├── proto/
│   └── gen/                   # Proto and/or generated code location (placeholder)
│
├── db/                        # SQL/schema/migrations (direction)
├── docker/                    # Dockerfiles
├── docs/                      # Documents & diagrams
├── rests/                     # HTTP client files for dev
├── third_party/               # OpenAPI and external assets
├── tools/                     # Tooling (protoc generators, sqlc, ...)
├── go.mod
└── go.sum
```

## Local run (scaffolding phase)

This repository focuses on the architecture scaffolding. Some services/IDLs are not fully implemented yet, so end-to-end flow may not run out of the box.

- Reverse proxy entrypoint: `cmd/proxy/`
- Docker build for proxy: `docker/Dockerfile-proxy`

## Notes

- Architecture diagram: `docs/Architechture-design.png`.
- `cmd/notification`, `cmd/user`, and `internal/notification/*` are currently skeletons to be implemented following the design.
