package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"apibgo/internal/config"
	domainAuth "apibgo/internal/domain/auth"
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
	// route: /auth/login/
	r.HandleFunc("/auth/login/", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)
		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		b, _ := io.ReadAll(r.Body)
		dto := domainAuth.LoginDto{}
		_ = json.Unmarshal(b, &dto)

		dto.Ip = utils.RealIp(r)
		dto.UserAgent = r.UserAgent()

		response, err := authService.Login(context.Background(), dto)

		if err != nil {
			log.Error("failed to execute Login service", slog.Err(err))
			return
		}

		response.SetCookies(&w, log)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
	}).Methods(http.MethodPost)

	// route: /auth/registration/
	r.HandleFunc("/auth/registration/", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)
		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		b, _ := io.ReadAll(r.Body)
		dto := domainAuth.RegistrationDto{}
		_ = json.Unmarshal(b, &dto)

		response, err := authService.Registration(context.Background(), dto)

		if err != nil {
			log.Error("failed to execute Registration service", slog.Err(err))
			return
		}

		if response.HttpCode == 0 {
			response.HttpCode = http.StatusOK
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
		w.WriteHeader(response.HttpCode)
	}).Methods(http.MethodPost)

	// route: /auth/logout/
	r.HandleFunc("/auth/logout/", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Authorization"]; !ok {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		log := logger.Setup(a.Config.Env)
		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		response, err := authService.Logout(context.Background(), r.Header["Authorization"])

		if err != nil {
			log.Error("failed to execute Logout service", slog.Err(err))
			return
		}

		response.SetCookies(&w, log)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
	}).Methods(http.MethodPost)

}
