# HelloService Implementation Playbook (Hexagonal + DDD + EDA)

## 1. Scope and Goal

Implement `hello.v1.HelloService` from `zhinux-contracts/hello/v1/hello.proto` in `zhinux-hello` with production-grade structure:

- Hexagonal architecture (ports/adapters)
- DDD-friendly layering (domain/application boundaries)
- EDA integration (event publication and async processing hooks)
- Full gRPC support (unary, server stream, client stream, bidi stream)

This playbook is written for a senior team and is intentionally explicit about file creation and ownership.

## 2. Current State Assessment

### Overall rating against large-scale production standards

- Structural intent: `8/10` (directories already reflect Hexagonal separation)
- Implementation maturity: `1/10` (almost all files are placeholders)
- Production readiness: `2/10` (missing runtime, tests, quality gates, observability)
- Hexagonal DDD EDA fit today: `3/10` (shape exists, behavior does not)

### Key gaps right now

- No executable service runtime (`cmd/main.go` is empty)
- No domain model, policies, or invariants
- No application orchestration logic
- No adapter implementations (grpc/db/http/mq are mostly empty)
- No config/logging/metrics/tracing platform layer
- No test pyramid (unit/integration/e2e)
- No CI/lint/test/release automation in this repo

## 3. Target Architecture and File Plan

Create and/or populate the following files.

```text
zhinux-hello/
  go.mod
  go.sum
  README.md
  Makefile
  Dockerfile
  .golangci.yml
  cmd/
    hello/
      main.go
  internal/
    domain/
      greeting/
        aggregate.go
        policy.go
        value_objects.go
        errors.go
      session/
        session.go
    application/
      hello/
        service.go
        commands.go
        streams.go
        dto.go
    ports/
      inbound/
        hello_usecase.go
      outbound/
        greeting_repository.go
        event_publisher.go
        clock.go
        id_generator.go
    adapters/
      grpc/
        server.go
        mapper.go
        interceptors.go
        stream_hub.go
      db/
        memory/
          greeting_repository.go
      mq/
        noop/
          publisher.go
      http/
        health/
          handler.go
    platform/
      config/
        config.go
      logging/
        logger.go
      telemetry/
        metrics.go
      server/
        grpc_server.go
        http_server.go
      shutdown/
        graceful.go
  tests/
    integration/
      hello_grpc_test.go
      hello_streams_test.go
```

## 4. Execution Plan (Phased)

## Phase 0: Bootstrap and Contracts Wiring

1. Create `go.mod` with module name for `zhinux-hello`.
2. Add dependency on contracts module:
   - Import path target: `github.com/amirhossein-shakeri/zhinux-contracts/gen/go/hello/v1`
   - For local dev, use a `replace` directive pointing to `../zhinux-contracts`.
3. Add `Makefile` targets:
   - `fmt`, `lint`, `test`, `test-integration`, `run`, `build`, `docker-build`
4. Populate `README.md` with:
   - local run
   - RPC examples
   - architecture map

Deliverable: `go test ./...` can execute successfully (even with minimal tests).

## Phase 1: Domain Model

Create:

- `internal/domain/greeting/aggregate.go`
- `internal/domain/greeting/policy.go`
- `internal/domain/greeting/value_objects.go`
- `internal/domain/greeting/errors.go`

Rules to enforce:

- `name` cannot be empty, max length policy configurable
- stream `count` bounded (for backpressure safety)
- message normalization (trim, UTF-8 safe checks)
- deterministic greeting composition policy

Deliverable: domain unit tests green; no framework imports in domain package.

## Phase 2: Ports (Hexagonal Contracts)

Create:

- `internal/ports/inbound/hello_usecase.go`
- `internal/ports/outbound/greeting_repository.go`
- `internal/ports/outbound/event_publisher.go`
- `internal/ports/outbound/clock.go`
- `internal/ports/outbound/id_generator.go`

Inbound port defines use-cases for:

- `SayHello`
- `StreamGreetings`
- `CollectGreetings`
- `Chat`

Outbound ports include:

- persistence abstraction for audit/session state
- event publication (`GreetingGenerated`, `ChatMessageReceived` optional)
- deterministic time and IDs for testability

Deliverable: application layer can compile against ports only.

## Phase 3: Application Service

Create:

- `internal/application/hello/service.go`
- `internal/application/hello/commands.go`
- `internal/application/hello/streams.go`
- `internal/application/hello/dto.go`

Responsibilities:

- Orchestrate domain + outbound ports
- Apply request-level validation and error mapping boundary
- Handle stream lifecycle and cancellation
- Emit domain events to outbound publisher port

Deliverable: pure application tests cover use-case orchestration and stream behavior.

## Phase 4: gRPC Adapter (Primary Adapter)

Create:

- `internal/adapters/grpc/server.go`
- `internal/adapters/grpc/mapper.go`
- `internal/adapters/grpc/interceptors.go`
- `internal/adapters/grpc/stream_hub.go`

Implement:

- gRPC server implementing generated `HelloServiceServer`
- mapping between proto messages and application DTOs
- interceptors:
  - request ID propagation
  - structured logging
  - panic recovery
  - timeout enforcement

Error mapping:

- validation errors -> `codes.InvalidArgument`
- internal errors -> `codes.Internal`
- canceled/timeout -> `codes.Canceled` / `codes.DeadlineExceeded`

Deliverable: all 4 RPC patterns behave according to proto contract.

## Phase 5: Secondary Adapters (EDA + Infra Hooks)

Create:

- `internal/adapters/db/memory/greeting_repository.go`
- `internal/adapters/mq/noop/publisher.go`
- `internal/adapters/http/health/handler.go`

Notes:

- Start with in-memory repository for deterministic tests.
- Keep publisher adapter as noop or structured log publisher first.
- Health endpoint supports Kubernetes probes later.

Deliverable: service runs standalone with safe defaults.

## Phase 6: Platform Layer and Bootstrap

Create:

- `internal/platform/config/config.go`
- `internal/platform/logging/logger.go`
- `internal/platform/telemetry/metrics.go`
- `internal/platform/server/grpc_server.go`
- `internal/platform/server/http_server.go`
- `internal/platform/shutdown/graceful.go`
- `cmd/hello/main.go`

Bootstrap flow:

1. Load config
2. Build logger + metrics registry
3. Wire adapters and application service
4. Start gRPC + health HTTP servers
5. Handle SIGINT/SIGTERM graceful drain

Deliverable: `make run` starts both listeners and shuts down cleanly.

## Phase 7: Testing and Quality Gates

Create:

- `tests/integration/hello_grpc_test.go`
- `tests/integration/hello_streams_test.go`

Coverage targets:

- unary happy + invalid input
- server streaming bounded count and ordering
- client streaming aggregation correctness
- bidi streaming concurrent clients and cancellation

Quality gates:

- `go test ./... -race`
- `golangci-lint run`
- minimal benchmark for stream path (`go test -bench=. ./...`)

Deliverable: CI green with race + lint + integration tests.

## Phase 8: Delivery Assets

Populate:

- `Dockerfile` (multi-stage build, non-root runtime)
- `helm/` chart skeleton:
  - `Chart.yaml`
  - `values.yaml`
  - templates for deployment/service/probes

Deliverable: container image and local helm template validation.

## 5. Team Work Allocation (Suggested)

- Domain lead: Phase 1
- Application lead: Phase 3
- Transport lead: Phase 4
- Platform/DevEx lead: Phases 0, 6, 8
- QA lead: Phase 7

Run parallel where safe, but lock API contracts at phase boundaries.

## 6. Senior-Level Design Constraints (Must-Haves)

- No adapter/framework dependency in domain/application packages
- Explicit context propagation in all use-cases and streams
- No global state for chat/session handling
- Backpressure and stream limits are explicit and tested
- Deterministic tests via injected clock/id
- Structured logs with correlation/request IDs

## 7. Optional EDA Enrichment

For stronger EDA depth, add event contracts in `zhinux-contracts/events/hello/v1/`:

- `greeting_generated.proto`
- `greetings_collected.proto`
- `chat_message_received.proto`

Then implement `internal/adapters/mq/nats/publisher.go` and emit these events from application layer.

## 8. Challenge RPC Extensions (Recommended)

Additive methods (backward-compatible for `hello.v1`):

1. `ScheduleGreeting` (async command)
   - returns operation id; processing occurs async
2. `GetScheduledGreetingStatus` (query)
   - returns state machine (`PENDING|RUNNING|DONE|FAILED`)
3. `WatchScheduledGreeting` (server stream)
   - push status transitions
4. `ModeratedChat` (bidi stream)
   - includes moderation decision payload and policy violations
5. `TemplateRender` (unary)
   - applies locale/template variables with validation and fallbacks

Why these are senior-level:

- require command/query separation
- require eventual consistency handling
- require stream fan-out/concurrency safety
- require policy engine and structured error taxonomy
- require idempotency + retry design

## 9. Definition of Done

- All `HelloService` RPCs implemented and tested
- Architecture boundaries enforced by package dependencies
- Basic observability (logs + metrics + health)
- CI quality gates are mandatory and green
- Service is runnable locally and containerized
- Extension track documented with at least one implemented advanced RPC
