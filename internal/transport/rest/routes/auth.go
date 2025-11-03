package routes

import (
	myhttp "apibgo/internal/utils/http"
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

	r.HandleFunc("/auth/logout/", a.AuthLogout).Methods(http.MethodPost)

	r.HandleFunc("/auth/refresh/", a.AuthRefresh).Methods(http.MethodGet)

	r.HandleFunc("/auth/verify/", a.AuthVerify).Methods(http.MethodGet)

	r.HandleFunc("/auth/activation/", a.AuthActivation).Methods(http.MethodPatch)

	r.HandleFunc("/auth/forgot/", a.AuthForgot).Methods(http.MethodPost)

	r.HandleFunc("/auth/recovery/", a.AuthRecovery).Methods(http.MethodPost)

	r.HandleFunc("/auth/confirm-check/", a.AuthConfirmCheck).Methods(http.MethodPost)

	r.HandleFunc("/auth/resend/{section}/", a.AuthResend).Methods(http.MethodPost)
}

// HandleAuthLogin handles authentication login.
// @Summary Handle authentication login
// @Description Handles user authentication.
// @Tags Auth
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 200 {object} response.DocSuccessResponse
// @Failure 400 {object} response.DocErrorResponse
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
		_response := response.Response{
			Code:    response.ErrorValidation,
			Message: "validation error",
			Result:  failMessages,
			Status:  response.StatusError,
		}
		w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
		w.Write(_response.CreateResponseData())
		return
	}

	dto.Ip = utils.RealIp(r)
	dto.UserAgent = r.UserAgent()

	_response, err := authService.Login(context.Background(), dto)

	if err != nil {
		log.Error("failed to execute Login service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if _response == nil {
			log.Error("response is empty")
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	_response.SetCookies(&w, log)
	w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
	w.Write(_response.CreateResponseData())
}

// HandleAuthLogin handles registration user.
// @Summary Handle registration account
// @Description Handles user account registration.
// @Tags Auth
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Param confirm_password body string true "Confirm Password"
// @Param name body string true "Name"
// @Param surname body string true "Surname"
// @Success 200 {object} response.DocSuccessResponse
// @Failure 400 {object} response.DocErrorResponse
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
		_response := response.Response{
			Code:    response.ErrorValidation,
			Message: "validation error",
			Result:  failMessages,
			Status:  response.StatusError,
		}
		w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
		w.Write(_response.CreateResponseData())
		return
	}

	_response, err := authService.Registration(context.Background(), dto)

	if err != nil {
		log.Error("failed to execute Registration service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _response.HttpCode == 0 {
		_response.HttpCode = http.StatusOK
	}

	w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
	w.Write(_response.CreateResponseData())
	w.WriteHeader(_response.HttpCode)
}

// HandleAuthLogin handles logout user.
// @Summary Handle logout account
// @Description Handles logout user account.
// @Tags Auth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.DocSuccessResponse
// @Failure 422 {object} response.DocErrorResponse
// @Router /auth/logout [post]
func (a *Auth) AuthLogout(w http.ResponseWriter, r *http.Request) {
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
	_response, err := authService.Logout(context.Background(), r.Header["Authorization"])

	if err != nil {
		log.Error("failed to execute Logout service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_response.SetCookies(&w, log)
	w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
	w.Write(_response.CreateResponseData())
}

// HandleAuthLogin handles refresh jwt tokens.
// @Summary Handle refresh jwt tokens
// @Description Handles refresh jwt tokens.
// @Tags Auth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.DocSuccessResponse
// @Failure 422 {object} response.DocErrorResponse
// @Router /auth/refresh [get]
func (a *Auth) AuthRefresh(w http.ResponseWriter, r *http.Request) {
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
	_response, err := authService.Refresh(context.Background(), cookie, dto)

	if err != nil {
		log.Error("failed to execute Refresh service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_response.SetCookies(&w, log)
	w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
	w.Write(_response.CreateResponseData())
}

// HandleAuthLogin handles verify jwt token.
// @Summary Handle verify jwt token
// @Description Handles verify jwt token.
// @Tags Auth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} nil
// @Failure 401 {object} nil
// @Router /auth/verify [get]
func (a *Auth) AuthVerify(w http.ResponseWriter, r *http.Request) {
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
}

// HandleAuthLogin handles account activation.
// @Summary Handle account activation
// @Description Handles a user account activation.
// @Tags Auth
// @Param email body string true "Email"
// @Param code body string true "Code"
// @Param key body string true "Signed key"
// @Success 200 {object} response.DocSuccessResponse
// @Failure 422 {object} response.DocErrorResponse
// @Router /auth/activation [patch]
func (a *Auth) AuthActivation(w http.ResponseWriter, r *http.Request) {
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

	_response, err := authService.Activation(context.Background(), dto)

	if err != nil {
		log.Error("failed to execute Activation service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if _response == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
	w.Write(_response.CreateResponseData())
}

// HandleAuthLogin handles account forgot password.
// @Summary Handle account forgot password
// @Description Handles a user account forgot password.
// @Tags Auth
// @Param email body string true "Email"
// @Success 200 {object} response.DocSuccessResponse
// @Failure 422 {object} response.DocErrorResponse
// @Router /auth/forgot [post]
func (a *Auth) AuthForgot(w http.ResponseWriter, r *http.Request) {
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

	_response, err := authService.Forgot(context.Background(), dto)

	if err != nil {
		log.Error("failed to execute Forgot service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if _response == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
	w.Write(_response.CreateResponseData())
}

// HandleAuthLogin handles account recovery password.
// @Summary Handle account recovery password
// @Description Handles a user account recovery password.
// @Tags Auth
// @Param email body string true "Email"
// @Param code body string true "Code"
// @Param password body string true "Password"
// @Param confirm_password body string true "Confirm Password"
// @Success 200 {object} response.DocSuccessResponse
// @Failure 422 {object} response.DocErrorResponse
// @Router /auth/recovery [post]
func (a *Auth) AuthRecovery(w http.ResponseWriter, r *http.Request) {
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

	_response, err := authService.Recovery(context.Background(), dto)

	if err != nil {
		log.Error("failed to execute Recovery service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if _response == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
	w.Write(_response.CreateResponseData())
}

// HandleAuthLogin handles checks confirm code.
// @Summary Handle checks confirm code
// @Description Handles checks confirm code.
// @Tags Auth
// @Param email body string true "Email"
// @Param action body service.SectionConfirm true "Action"
// @Param code body string true "Code"
// @Success 200 {object} response.DocSuccessResponse
// @Failure 422 {object} response.DocErrorResponse
// @Router /auth/confirm-check [post]
func (a *Auth) AuthConfirmCheck(w http.ResponseWriter, r *http.Request) {
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

	_response, err := authService.ConfirmCheck(context.Background(), dto)

	if err != nil {
		log.Error("failed to execute ConfirmCheck service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if _response == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
	w.Write(_response.CreateResponseData())
}

// HandleAuthLogin handles resend confirm code.
// @Summary Handle resend confirm code
// @Description Handles resend confirm code.
// @Tags Auth
// @Param section path service.SectionSend true "Resend for section"
// @Param email body string true "Email"
// @Success 200 {object} response.DocSuccessResponse
// @Failure 422 {object} response.DocErrorResponse
// @Router /auth/resend/{section}/ [post]
func (a *Auth) AuthResend(w http.ResponseWriter, r *http.Request) {
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

	_response, err := authService.Resend(context.Background(), service.SectionSend(vars["section"]), jsond)

	if err != nil {
		log.Error("failed to execute Resend service", slog.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if _response == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", string(myhttp.ContentType_JSON))
	w.Write(_response.CreateResponseData())
}
