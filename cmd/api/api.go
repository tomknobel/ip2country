package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddl "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/tomknobel/ip2country/internal/db"
	"github.com/tomknobel/ip2country/internal/routes"
	"github.com/tomknobel/ip2country/pkg/middleware"
	"net/http"

	"go.uber.org/zap"
)

type Application struct {
	cfg    config
	Logger *zap.SugaredLogger
}

func NewApplication() *Application {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	return &Application{
		cfg:    initConfig(),
		Logger: zap.NewExample().Sugar(),
	}
}

func (app *Application) newRouts(middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddl.Logger)
	r.Use(middlewares...)
	ipDb := db.DbFactory(app.cfg.dbType, app.cfg.dbCfg)

	r.Route("/v1", func(v1 chi.Router) {
		routes.InitIp2CountryRouter(v1, ipDb)
	})
	return r
}

func (app *Application) Run() {
	rateLimiterMiddleware := middleware.NewRateLimitMiddleware(
		app.cfg.limiter.rateLimit,
		app.cfg.limiter.windowSize,
		app.Logger,
	)
	r := app.newRouts(rateLimiterMiddleware.RateLimiterByIp, chiMiddl.Logger)
	app.Logger.Infof("Starting the API server at port %s", app.cfg.port)
	addr := fmt.Sprintf(":%s", app.cfg.port)
	app.Logger.Fatal(http.ListenAndServe(addr, r))
}
