package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Leul-Michael/go-auth/application"
	"github.com/Leul-Michael/go-auth/global"
	"github.com/Leul-Michael/go-auth/utils"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	global.Validate = validator.New()
	global.Validate.RegisterValidation("phone", utils.ValidatePhoneNumber)
}

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatal("failed to initialize the app: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = app.Start(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
