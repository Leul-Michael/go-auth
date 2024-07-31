package application

import (
	"github.com/Leul-Michael/go-auth/handler"
	repository "github.com/Leul-Michael/go-auth/repository/user"
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
	v1Router.Route("/auth", a.loadUserRoutes)

	a.router = router
}

func (a *App) loadUserRoutes(router chi.Router) {
	userHandler := &handler.User{
		Repo: &repository.PostgresUserRepo{
			DB: a.DB,
		},
	}

	router.Post("/signup", userHandler.Signup)
	router.Post("/login", userHandler.Login)
}
