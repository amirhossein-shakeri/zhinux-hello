# Platform Adoption Sequencing (zhinux-hello)

## Decision

`zhinux-hello` adopts `zhinux-platform` immediately for cross-cutting baselines and adopts transport/runtime packages when service runtime wiring is implemented.

## Why this sequence

- Phase 1 domain work should stay business-focused but can use neutral shared validation helpers.
- Runtime concerns (grpc lifecycle, interceptors, telemetry, health, shutdown) become high-value once listeners are introduced.
- Early adoption of config/logging removes one-off bootstrapping patterns and aligns all services on the same startup contract.

## Adopted now

- `config`: shared typed base configuration
- `logging`: shared structured logging baseline
- `validation`: shared input normalization and scalar guardrails

## Adopt next (Phase 4 and Phase 6)

- `grpc` + `grpc/interceptors`: server runtime and middleware chain
- `telemetry`: metrics interceptors and HTTP metrics endpoint
- `health`: health probe registry/handlers
- `shutdown`: coordinated graceful-stop hooks
- `grpc.ToStatusError`: centralized app-error to gRPC-code mapping

## Reusable logic moved out of hello

Moved to `zhinux-platform/validation` so future services can reuse the same primitives:

- UTF-8-safe normalized text transformation
- Rune-length guard helper
- Integer range guard helper
