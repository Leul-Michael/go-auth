package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Leul-Michael/go-auth/application"
)

func main() {
	app := application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
