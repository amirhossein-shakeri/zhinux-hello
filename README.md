# zhinux-hello

`zhinux-hello` is the HelloService implementation workspace following Hexagonal Architecture, DDD layering, and EDA-ready adapter boundaries.

## Local Run

Prerequisites:

- Go `1.25.6`
- `make`
- `golangci-lint` (only for `make lint`)
- Docker (only for `make docker-build`)

Bootstrap and verify:

```bash
go mod tidy
make test
```

Run from source:

```bash
make run
```

Build binary:

```bash
make build
./bin/zhinux-hello
```

## RPC Examples (grpcurl)

Once gRPC server wiring is active (Phase 4), these commands are the baseline smoke checks.

Unary:

```bash
grpcurl -plaintext \
  -d '{"name":"Alice"}' \
  localhost:50051 hello.v1.HelloService/SayHello
```

Server stream:

```bash
grpcurl -plaintext \
  -d '{"name":"Alice","count":3}' \
  localhost:50051 hello.v1.HelloService/StreamGreetings
```

Client stream:

```bash
grpcurl -plaintext -d @ localhost:50051 hello.v1.HelloService/CollectGreetings <<'EOF'
{"name":"Alice"}
{"name":"Bob"}
EOF
```

Bidirectional stream:

```bash
grpcurl -plaintext -d @ localhost:50051 hello.v1.HelloService/Chat <<'EOF'
{"name":"Alice","message":"hi"}
{"name":"Bob","message":"hello"}
EOF
```

## Architecture Map

Top-level boundaries:

- `cmd/`: process bootstrap and runtime entrypoint
- `internal/domain/`: entities, value objects, and business policies
- `internal/application/`: use-case orchestration and stream coordination
- `internal/ports/`: inbound and outbound interfaces for hexagonal boundaries
- `internal/adapters/grpc/`: primary adapter exposing HelloService over gRPC
- `internal/adapters/db/`: persistence adapters (in-memory first)
- `internal/adapters/mq/`: event publication adapters (noop first)
- `internal/adapters/http/`: operational HTTP endpoints (health, diagnostics)
- `tests/integration/`: black-box integration tests for transport contracts

Contracts module wiring for local development:

- Proto-generated package import path: `github.com/amirhossein-shakeri/zhinux-contracts/gen/go/hello/v1`
- `go.mod` uses a local replace: `github.com/amirhossein-shakeri/zhinux-contracts => ../zhinux-contracts`
