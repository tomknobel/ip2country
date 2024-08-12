package api

import (
	"github.com/go-chi/chi/v5"
	chiMiddl "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/tomknobel/ip2country/pkg/middleware"
	"net/http"

	"go.uber.org/zap"
)

type Application struct {
	cfg    Config
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
func (app *Application) Run() {
	rateLimiterMiddleware := middleware.NewRateLimitMiddleware(
		app.cfg.rateLimiter.Rate,
		app.cfg.rateLimiter.Window,
		app.Logger,
	)

	r := chi.NewRouter()
	r.Use(chiMiddl.Logger)
	r.Use(rateLimiterMiddleware.RateLimiterByIp)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	app.Logger.Infof("Starting the API server at", app.cfg.port)
	app.Logger.Fatal(http.ListenAndServe(":3000", r))
	//adrr := fmt.Sprintf(":%s", &app.cfg)
	//return http.ListenAndServe(adrr, r)
}
