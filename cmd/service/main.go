package main

import (
	"log/slog"
	"os"

	_ "github.com/LevanPro/insider/docs"
	"github.com/LevanPro/insider/internal/api"
)

// @title           UseInsder Message Sender API
// @version         1.0
// @description     Automatic 2-minute message sending service.
// @BasePath        /

// @host      localhost:8080
// @schemes   http
func main() {
	if err := api.Run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}
