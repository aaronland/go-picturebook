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
	logger := log.Default()

	err := commandline.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run picturebook application, %v", err)
	}

	os.Exit(0)
}
