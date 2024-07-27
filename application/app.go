package application

import (
	"context"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	router http.Handler
	db     *gorm.DB
}

func New() *App {
	app := &App{}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":5000",
		Handler: a.router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
