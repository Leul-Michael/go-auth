package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RealIP)

	// api v1 route handler
	v1Router := chi.NewRouter()
	router.Mount("/api/v1", v1Router)

	a.router = router
}
