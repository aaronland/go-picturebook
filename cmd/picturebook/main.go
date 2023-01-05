// picturebook is a command-line application for creating a PDF file from a folder containing images.
package main

import (
	"context"
	"github.com/aaronland/go-picturebook/app/picturebook"
	_ "gocloud.dev/blob/fileblob"
	"log"
	"os"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := picturebook.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run picturebook application, %v", err)
	}

	os.Exit(0)
}
