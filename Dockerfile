# syntax=docker/dockerfile:1.7

##
## zhinux-hello production image
##
## This Dockerfile assumes a monorepo build context rooted at `zhinux/`
## so that both `zhinux-hello/` and `zhinux-contracts/` are available.
## Build with:
##   docker build -f zhinux-hello/Dockerfile -t zhinux-hello:dev .
##

ARG GO_VERSION=1.25.6
ARG ALPINE_VERSION=3.22
ARG RUNTIME_IMAGE=gcr.io/distroless/static-debian12:nonroot

##
## Builder base
## - pinned Go toolchain
## - certs + tzdata for parity with production network/time behavior
##
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build-base
RUN apk add --no-cache ca-certificates tzdata git && update-ca-certificates
WORKDIR /workspace
ENV CGO_ENABLED=0
ENV GOFLAGS="-trimpath"

##
## Dependency layer
## - copies only module manifests first to maximize build cache reuse
##
FROM build-base AS deps
COPY zhinux-hello/go.mod /workspace/zhinux-hello/go.mod
COPY zhinux-contracts/go.mod /workspace/zhinux-contracts/go.mod
COPY zhinux-platform/go.mod /workspace/zhinux-platform/go.mod
WORKDIR /workspace/zhinux-hello
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

##
## Build and test layer
## - pulls in full sources after dependencies are cached
## - runs tests as an explicit quality gate in the container pipeline
## - emits one statically-linked artifact
##
FROM build-base AS build
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_DATE=unknown
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY zhinux-hello /workspace/zhinux-hello
COPY zhinux-contracts /workspace/zhinux-contracts
COPY zhinux-platform /workspace/zhinux-platform
WORKDIR /workspace/zhinux-hello
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go test ./...
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
      -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT_SHA} -X main.buildDate=${BUILD_DATE}" \
      -o /out/zhinux-hello ./cmd

##
## Runtime layer
## - distroless non-root image minimizes attack surface
## - binary-only deployment for immutable runtime behavior
##
FROM ${RUNTIME_IMAGE} AS runtime
LABEL org.opencontainers.image.title="zhinux-hello"
LABEL org.opencontainers.image.description="HelloService gRPC implementation"
LABEL org.opencontainers.image.source="https://github.com/amirhossein-shakeri/zhinux"
LABEL org.opencontainers.image.vendor="zhinux"
WORKDIR /app
COPY --from=build /out/zhinux-hello /app/zhinux-hello
USER nonroot:nonroot
EXPOSE 50051 8080
STOPSIGNAL SIGTERM
ENTRYPOINT ["/app/zhinux-hello"]
