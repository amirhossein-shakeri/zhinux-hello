package main

import "log"

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	log.Printf(
		"zhinux-hello bootstrap: runtime wiring not implemented yet (version=%s commit=%s buildDate=%s)",
		version,
		commit,
		buildDate,
	)
}
