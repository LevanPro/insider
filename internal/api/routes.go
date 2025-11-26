package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *App) setupRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	router.Post("/api/v1/scheduler/start", app.StartScheduler)
	router.Post("/api/v1/scheduler/stop", app.StopScheduler)
	router.Get("/api/v1/scheduler/status", app.SchedulerStatus)
	router.Get("/api/v1/messages/sent", app.GetSentMessages)

	router.Get("/debug/liveness", app.Liveness)
	router.Get("/debug/readiness", app.Readiness)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	return router
}
