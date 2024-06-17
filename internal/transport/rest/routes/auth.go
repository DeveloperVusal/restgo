package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"apibgo/internal/config"
	domainAuth "apibgo/internal/domain/auth"
	"apibgo/internal/service"
	"apibgo/internal/storage"
	"apibgo/internal/storage/pgsql"
	"apibgo/internal/utils/request"
	"apibgo/internal/utils/response"
	"apibgo/pkg/logger"
	"apibgo/pkg/logger/feature/slog"
	"apibgo/pkg/utils"

	_ "apibgo/docs/swagger"

	"github.com/gorilla/mux"
)

type Auth struct {
	Config  *config.Config
	Storage *storage.Config
}

func (a *Auth) NewHandler(r *mux.Router) {

	r.HandleFunc("/auth/login/", a.AuthLogin).Methods(http.MethodPost)

	r.HandleFunc("/auth/registration/", a.AuthRegistration).Methods(http.MethodPost)

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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		response, err := authService.Logout(context.Background(), r.Header["Authorization"])

		if err != nil {
			log.Error("failed to execute Logout service", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.SetCookies(&w, log)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
	}).Methods(http.MethodPost)

	// route: /auth/refresh/
	r.HandleFunc("/auth/refresh/", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)
		cookie, err := r.Cookie("refresh_token")

		if err != nil {
			log.Error("failed to get cookie", slog.Err(err))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("starting database")

		dto := domainAuth.LoginDto{
			Ip:        utils.RealIp(r),
			UserAgent: r.UserAgent(),
		}

		authService := service.NewAuthService(pg)
		response, err := authService.Refresh(context.Background(), cookie, dto)

		if err != nil {
			log.Error("failed to execute Refresh service", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.SetCookies(&w, log)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
	}).Methods(http.MethodGet)

	// route: /auth/verify/
	r.HandleFunc("/auth/verify/", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)

		if _, ok := r.Header["Authorization"]; !ok {
			log.Error("failed to get header of the Authorization")
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		// Parse header Authorization and get token
		split := strings.Split(r.Header["Authorization"][0], " ")
		token := split[1]
		isVerify, err := authService.VerifyToken(context.Background(), token)

		if err != nil {
			log.Error("failed to execute VerifyToken service", slog.Err(err))
		}

		if !isVerify {
			w.WriteHeader(http.StatusUnauthorized)
		}

		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// route: /auth/activation/
	r.HandleFunc("/auth/activation/", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)
		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		b, _ := io.ReadAll(r.Body)
		dto := domainAuth.ActivationDto{}
		_ = json.Unmarshal(b, &dto)

		response, err := authService.Activation(context.Background(), dto)

		if err != nil {
			log.Error("failed to execute Activation service", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if response == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
	}).Methods(http.MethodPost)

	// route: /auth/forgot/
	r.HandleFunc("/auth/forgot/", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)
		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		b, _ := io.ReadAll(r.Body)
		dto := domainAuth.ForgotDto{}
		_ = json.Unmarshal(b, &dto)

		response, err := authService.Forgot(context.Background(), dto)

		if err != nil {
			log.Error("failed to execute Forgot service", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if response == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
	}).Methods(http.MethodPost)

	// route: /auth/recovery/
	r.HandleFunc("/auth/recovery/", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)
		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		b, _ := io.ReadAll(r.Body)
		dto := domainAuth.RecoveryDto{}
		_ = json.Unmarshal(b, &dto)

		response, err := authService.Recovery(context.Background(), dto)

		if err != nil {
			log.Error("failed to execute Recovery service", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if response == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
	}).Methods(http.MethodPost)

	// route: /auth/confirm_check/
	r.HandleFunc("/auth/confirm_check/", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)
		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("starting database")

		authService := service.NewAuthService(pg)
		b, _ := io.ReadAll(r.Body)
		dto := domainAuth.ConfirmCheckDto{}
		_ = json.Unmarshal(b, &dto)

		response, err := authService.ConfirmCheck(context.Background(), dto)

		if err != nil {
			log.Error("failed to execute ConfirmCheck service", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if response == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
	}).Methods(http.MethodPost)

	// route: /auth/resend/{section}
	r.HandleFunc("/auth/resend/{section}/", func(w http.ResponseWriter, r *http.Request) {
		log := logger.Setup(a.Config.Env)
		pg, err := pgsql.New(a.Storage, "master")

		if err != nil {
			log.Error("failed to init storage", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("starting database")

		vars := mux.Vars(r)
		authService := service.NewAuthService(pg)
		jsond, _ := io.ReadAll(r.Body)

		response, err := authService.Resend(context.Background(), service.SectionSend(vars["section"]), jsond)

		if err != nil {
			log.Error("failed to execute Resend service", slog.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if response == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
	}).Methods(http.MethodPost)
}

// HandleAuthLogin handles authentication login.
// @Summary Handle authentication login
// @Description Handles user authentication.
// @Tags auth
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 200 {object} response.DocResponse
// @Failure 400 {object} response.DocResponse
// @Router /auth/login [post]
func (a *Auth) AuthLogin(w http.ResponseWriter, r *http.Request) {
	log := logger.Setup(a.Config.Env)
	pg, err := pgsql.New(a.Storage, "master")

	if err != nil {
		log.Error("failed to init storage", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("starting database")

	authService := service.NewAuthService(pg)
	b, _ := io.ReadAll(r.Body)
	dto := domainAuth.LoginDto{}
	_ = json.Unmarshal(b, &dto)

	validator := request.NewValidator()
	isValid, failMessages := validator.Validate(dto)

	if !isValid {
		response := response.Response{
			Code:    response.ErrorValidation,
			Message: "validation error",
			Result:  failMessages,
			Status:  response.StatusError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
		return
	}

	dto.Ip = utils.RealIp(r)
	dto.UserAgent = r.UserAgent()

	response, err := authService.Login(context.Background(), dto)

	if err != nil {
		log.Error("failed to execute Login service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if response == nil {
			log.Error("response is empty")
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	response.SetCookies(&w, log)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response.CreateResponseData())
}

// HandleAuthLogin handles authentication login.
// @Summary Handle authentication login
// @Description Handles user authentication.
// @Tags auth
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Param confirm_password body string true "Confirm Password"
// @Param name body string true "Name"
// @Param surname body string true "Surname"
// @Success 200 {object} response.DocResponse
// @Failure 400 {object} response.DocResponse
// @Router /auth/registration [post]
func (a *Auth) AuthRegistration(w http.ResponseWriter, r *http.Request) {
	log := logger.Setup(a.Config.Env)
	pg, err := pgsql.New(a.Storage, "master")

	if err != nil {
		log.Error("failed to init storage", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("starting database")

	authService := service.NewAuthService(pg)
	b, _ := io.ReadAll(r.Body)
	dto := domainAuth.RegistrationDto{}
	_ = json.Unmarshal(b, &dto)

	validator := request.NewValidator()
	isValid, failMessages := validator.Validate(dto)

	if !isValid {
		response := response.Response{
			Code:    response.ErrorValidation,
			Message: "validation error",
			Result:  failMessages,
			Status:  response.StatusError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response.CreateResponseData())
		return
	}

	response, err := authService.Registration(context.Background(), dto)

	if err != nil {
		log.Error("failed to execute Registration service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.HttpCode == 0 {
		response.HttpCode = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response.CreateResponseData())
	w.WriteHeader(response.HttpCode)
}
