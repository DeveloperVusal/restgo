package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"apibgo/internal/config"
	"apibgo/internal/domain"
	"apibgo/internal/service"
	"apibgo/internal/storage"
	"apibgo/internal/storage/pgsql"
	"apibgo/pkg/logger"
	"apibgo/pkg/logger/feature/slog"
	"apibgo/pkg/utils"

	"github.com/gorilla/mux"
)

type Auth struct {
	Config  *config.Config
	Storage *storage.Config
}

func (a *Auth) NewHandler(r *mux.Router) {
	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)
		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		b, _ := io.ReadAll(r.Body)
		dto := domain.AuthDto{}
		_ = json.Unmarshal(b, &dto)

		dto.Ip = utils.RealIp(r)
		dto.UserAgent = r.UserAgent()

		authService.Login(context.Background(), dto)

		w.Write([]byte("Welcome"))
	}).Methods(http.MethodPost)
}
