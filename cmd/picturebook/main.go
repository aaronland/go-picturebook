// picturebook is a command-line application for creating a PDF file from a folder containing images.
package main

import (
	"context"
	"log"
	"os"

	"github.com/aaronland/go-picturebook/application/commandline"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	ctx := context.Background()

	fs, err := commandline.DefaultFlagSet(ctx)

	if err != nil {
		log.Fatalf("Failed to create default flag set, %v", err)
	}

	app, err := commandline.NewApplication(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to create new picturebook application, %v", err)
	}

	err = app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run picturebook application, %v", err)
	}

	os.Exit(0)
}
