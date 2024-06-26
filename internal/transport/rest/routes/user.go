package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"apibgo/internal/config"
	domainUser "apibgo/internal/domain/user"
	"apibgo/internal/service"
	"apibgo/internal/storage"
	"apibgo/internal/storage/pgsql"
	"apibgo/internal/transport/rest"
	"apibgo/pkg/logger"
	"apibgo/pkg/logger/feature/slog"

	"github.com/gorilla/mux"
)

type User struct {
	Config      *config.Config
	Storage     *storage.Config
	Middlewares []mux.MiddlewareFunc
}

func (u *User) NewHandler(r *mux.Router) {
	// route: user sessions
	r.HandleFunc("/users/sessions/", rest.Adapt(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.Setup(u.Config.Env)
			pg, err := pgsql.New(u.Storage, "master")

			if err != nil {
				log.Error("failed to init storage", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Info("starting database for sessions")

			userService := service.NewUserService(pg)
			response, err := userService.Sessions(context.Background(), r.Header["Authorization"])

			if err != nil {
				log.Error("failed to execute Sessions service", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if response == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(response.CreateResponseData())
		}),
		u.Middlewares...,
	).ServeHTTP).Methods(http.MethodGet)

	// route: destroy user session
	r.HandleFunc("/users/sessions/{id}/", rest.Adapt(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.Setup(u.Config.Env)
			pg, err := pgsql.New(u.Storage, "master")

			if err != nil {
				log.Error("failed to init storage", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Info("starting database for sessions")

			vars := mux.Vars(r)
			userService := service.NewUserService(pg)
			paramId, _ := strconv.Atoi(vars["id"])

			response, err := userService.DestroySession(context.Background(), r.Header["Authorization"], paramId)

			if err != nil {
				log.Error("failed to execute Sessions service", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if response == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}

			if response.HttpCode == 0 {
				response.HttpCode = http.StatusOK
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.HttpCode)
			w.Write(response.CreateResponseData())
		}),
		u.Middlewares...,
	).ServeHTTP).Methods(http.MethodDelete)

	// route: get a user
	r.HandleFunc("/users/{id}/", rest.Adapt(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.Setup(u.Config.Env)
			pg, err := pgsql.New(u.Storage, "master")

			if err != nil {
				log.Error("failed to init storage", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Info("starting database")

			vars := mux.Vars(r)
			userService := service.NewUserService(pg)
			paramId, _ := strconv.Atoi(vars["id"])
			dto := domainUser.UserDto{
				Id: paramId,
			}

			response, err := userService.GetUser(context.Background(), dto)

			if err != nil {
				log.Error("failed to execute GetUser service", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if response == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(response.CreateResponseData())
		}),
		u.Middlewares...,
	).ServeHTTP).Methods(http.MethodGet)

	// route: get users
	r.HandleFunc("/users/", rest.Adapt(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.Setup(u.Config.Env)
			pg, err := pgsql.New(u.Storage, "master")

			if err != nil {
				log.Error("failed to init storage", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Info("starting database")

			userService := service.NewUserService(pg)
			response, err := userService.GetUsers(context.Background())

			if err != nil {
				log.Error("failed to execute GetUsers service", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if response == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(response.CreateResponseData())
		}),
		u.Middlewares...,
	).ServeHTTP).Methods(http.MethodGet)

	// route: create a user
	r.HandleFunc("/users/", rest.Adapt(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.Setup(u.Config.Env)
			pg, err := pgsql.New(u.Storage, "master")

			if err != nil {
				log.Error("failed to init storage", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Info("starting database")

			userService := service.NewUserService(pg)
			b, _ := io.ReadAll(r.Body)
			dto := domainUser.CreateUserDto{}
			_ = json.Unmarshal(b, &dto)

			response, err := userService.CreateUser(context.Background(), dto)

			if err != nil {
				log.Error("failed to execute CreateUser service", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if response == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}

			if response.HttpCode == 0 {
				response.HttpCode = http.StatusOK
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.HttpCode)
			w.Write(response.CreateResponseData())
		}),
		u.Middlewares...,
	).ServeHTTP).Methods(http.MethodPost)

	// route: update a user
	r.HandleFunc("/users/{id}/", rest.Adapt(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.Setup(u.Config.Env)
			pg, err := pgsql.New(u.Storage, "master")

			if err != nil {
				log.Error("failed to init storage", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Info("starting database")

			vars := mux.Vars(r)
			userService := service.NewUserService(pg)
			paramId, _ := strconv.Atoi(vars["id"])

			b, _ := io.ReadAll(r.Body)
			dto := domainUser.UpdateUserDto{}
			_ = json.Unmarshal(b, &dto)

			dto.Id = paramId

			response, err := userService.UpdateUser(context.Background(), dto)

			if err != nil {
				log.Error("failed to execute UpdateUser service", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if response == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(response.CreateResponseData())
		}),
		u.Middlewares...,
	).ServeHTTP).Methods(http.MethodPatch)

	// route: delete a user
	r.HandleFunc("/users/{id}/", rest.Adapt(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.Setup(u.Config.Env)
			pg, err := pgsql.New(u.Storage, "master")

			if err != nil {
				log.Error("failed to init storage", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Info("starting database")

			vars := mux.Vars(r)
			userService := service.NewUserService(pg)
			paramId, _ := strconv.Atoi(vars["id"])

			response, err := userService.DeleteUser(context.Background(), paramId)

			if err != nil {
				log.Error("failed to execute DeleteUser service", slog.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if response == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(response.CreateResponseData())
		}),
		u.Middlewares...,
	).ServeHTTP).Methods(http.MethodDelete)
}
