package main

import (
	"context"
	"flag"
	"log"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/app"
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()

		return
	}

	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err = a.Run(ctx); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
