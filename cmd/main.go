package main

import (
	"expvar"
	"flag"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/shingeki-no-kyojin/ymir/config"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/id"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/middleware"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/security"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/user"
	"gitlab.com/shingeki-no-kyojin/ymir/pkg/logger"
)

const version = "1.0.0"

func main() {
	var cfg *config.Config

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	logger := logger.New()

	flag.IntVar(&cfg.Port, "port", 4000, "Application port")
	flag.StringVar(&cfg.Env, "env", "dev", "Environment (dev|staging|prod)")
	flag.Parse()

	// migration
	// removed from showcase

	// cleaner
	// removed from showcase

	rootMux := http.NewServeMux()

	handler := setupRoutes(rootMux, logger, cfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// TODO: Graceful shutdown
	logger.Info(fmt.Sprintf("starting %s server on %s", cfg.Env, srv.Addr))
	err = srv.ListenAndServe()
	logger.Error(err.Error())
}

func setupRoutes(rootMux *http.ServeMux, logger logger.Logger, config *config.Config) http.Handler {
	repo := user.NewInMemoryUserRepository()
	idService := id.NewUUID()
	hashService := security.NewBcryptHash()
	tokenService := security.NewJWT(config.AccessSecret, config.RefreshSecret)

	userUsecases := user.NewUseCase(repo, idService, hashService, tokenService, time.Duration(config.AccessTokenExpire), time.Duration(config.RefreshTokenExpire))
	userController := user.NewController(logger, config, userUsecases)

	// public routes
	rootMux.Handle("GET /debug/vars", expvar.Handler())

	rootMux.HandleFunc("POST /user/registrieren", userController.CreateUser)
	rootMux.HandleFunc("POST /user/anmelden", userController.LoginUser)

	// private routes
	authMux := http.NewServeMux()
	// authMux.HandleFunc("GET /user/refreshtoken", authController.RefreshToken)
	authMux.HandleFunc("PUT /user/bearbeiten", userController.UpdateUser)
	authMux.HandleFunc("PUT /user/ausloggen", userController.LogoutUser)
	authMux.HandleFunc("DELETE /user/entfernen", userController.DeleteUser)

	authMiddleware := middleware.NewAuthorization(tokenService)
	rootMux.Handle("/", authMiddleware.Authorize(authMux))

	// middleware
	handler := middleware.Chain(
		middleware.RecoverPanic,
		middleware.NewLogger(nil).Log,
		middleware.EnableCORS,
	)(rootMux)

	return handler
}
