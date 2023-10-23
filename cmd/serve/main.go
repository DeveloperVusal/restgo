package main

import (
	"apibgo/internal/config"

	"apibgo/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.Setup(cfg.Env)

	log.Info("starting restapi server")
	log.Debug("debug messages are enabled")
}
