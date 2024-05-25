package instance

import (
	"log/slog"

	"apibgo/internal/config"
	"apibgo/internal/storage"
	"apibgo/pkg/logger"
	"apibgo/pkg/univenv"
)

type Instance struct {
	Config  *config.Config
	Storage *storage.Config
	Log     *slog.Logger
}

func GetInstance() *Instance {
	// Load .env files
	univenv.Load()

	cfg := config.MustLoad()
	log := logger.Setup(cfg.Env)
	dbcfgs := storage.MustLoad()

	return &Instance{
		Config:  cfg,
		Storage: dbcfgs,
		Log:     log,
	}
}
