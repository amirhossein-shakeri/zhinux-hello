package main

import (
	"log"
	"os"
	"strings"

	platformconfig "github.com/amirhossein-shakeri/zhinux-platform/config"
	platformlogging "github.com/amirhossein-shakeri/zhinux-platform/logging"
)

const defaultServiceName = "zhinux-hello"

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	if strings.TrimSpace(os.Getenv("SERVICE_NAME")) == "" {
		if err := os.Setenv("SERVICE_NAME", defaultServiceName); err != nil {
			log.Fatalf("set default service name: %v", err)
		}
	}

	cfg, err := platformconfig.LoadBaseFromEnv()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger, err := platformlogging.NewLogger(platformlogging.LoggerOptions{
		Level:       cfg.LogLevel,
		Service:     cfg.ServiceName,
		Backend:     platformlogging.BackendSlog,
		Development: cfg.Environment != "production",
	})
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}

	logger.Info(
		"zhinux-hello bootstrap initialized",
		platformlogging.KV("version", version),
		platformlogging.KV("commit", commit),
		platformlogging.KV("build_date", buildDate),
		platformlogging.KV("environment", cfg.Environment),
		platformlogging.KV("grpc_listen_addr", cfg.GRPCListenAddr),
		platformlogging.KV("http_listen_addr", cfg.HTTPListenAddr),
	)
}
