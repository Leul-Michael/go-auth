package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type App struct {
	router http.Handler
	DB     *gorm.DB
}

func New() (*App, error) {
	app := &App{}

	err := app.connectToDB()
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	// // making migrations
	// app.DB.AutoMigrate(&model.User{})

	app.loadRoutes()

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":5000",
		Handler: a.router,
	}

	ch := make(chan error, 1) // buffered channel to avoid blocking background context
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := server.Shutdown(timeout); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
	}
	return nil
}
