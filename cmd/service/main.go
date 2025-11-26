package main

import (
	"log/slog"
	"os"

	"github.com/LevanPro/insider/internal/api"
)

func main() {
	if err := api.Run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}
