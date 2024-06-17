package app

import (
	"log/slog"
	"net/http"

	"apibgo/internal/app/instance"
	"apibgo/internal/config"
	"apibgo/internal/storage"
	"apibgo/internal/transport/rest"
	"apibgo/internal/transport/rest/middleware"
	"apibgo/internal/transport/rest/routes"

	"apibgo/docs/swagger"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Instance struct {
	Config  *config.Config
	Storage *storage.Config
	Log     *slog.Logger
}

func Run() {
	instance := instance.GetInstance()

	_routes := []rest.Handler{
		&routes.Auth{Config: instance.Config, Storage: instance.Storage},
		&routes.User{
			Config:  instance.Config,
			Storage: instance.Storage,
			Middlewares: []mux.MiddlewareFunc{
				middleware.LoggingMiddleware,
			},
		},
	}

	r := mux.NewRouter()

	rest.NewRouter(r, _routes...)
	r.Use(mux.CORSMethodMiddleware(r))

	// Swagger UI
	swagger.SwaggerInfo.Title = "Swagger RESTGO API"
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	instance.Log.Info("starting restapi server at http://" + instance.Config.Address)
	instance.Log.Info("Swagger URL: http://" + instance.Config.Address + "/swagger/")

	if err := http.ListenAndServe(instance.Config.HTTPServer.Address, r); err != nil {
		instance.Log.Error("failed to start server", err)
	}

}
