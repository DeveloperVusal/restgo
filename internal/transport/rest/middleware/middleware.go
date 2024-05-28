package middleware

import (
	"context"
	"net/http"
	"strings"

	"apibgo/internal/app/instance"
	"apibgo/internal/service"
	"apibgo/internal/storage/pgsql"
	"apibgo/pkg/logger/feature/slog"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Authorization"]; !ok {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		instance := instance.GetInstance()
		pg, err := pgsql.New(instance.Storage, "master")

		if err != nil {
			instance.Log.Error("failed to init storage", slog.Err(err))
			return
		}

		instance.Log.Info("starting database middleware")

		authService := service.NewAuthService(pg)
		split := strings.Split(r.Header["Authorization"][0], " ")
		token := split[1]
		isVerify, err := authService.VerifyToken(context.Background(), token)

		if err != nil {
			instance.Log.Error("failed to execute VerifyToken service", slog.Err(err))
		}

		if !isVerify {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
