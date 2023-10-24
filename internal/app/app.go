package app

import (
	"apibgo/internal/config"

	"apibgo/pkg/logger"
)

func Run() {
	cfg := config.MustLoad()
	log := logger.Setup(cfg.Env)

	log.Info("starting restapi server")
	log.Debug("debug messages are enabled")
}
