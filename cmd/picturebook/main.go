package main

import (
	"context"
	"github.com/aaronland/go-picturebook/application"
	"log"
	"os"
)

func main() {

	ctx := context.Background()

	fs, err := application.CommandLineApplicationDefaultFlagSet(ctx)

	if err != nil {
		log.Fatal(err)
	}

	app, err := application.NewCommandLineApplication(ctx, fs)

	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
