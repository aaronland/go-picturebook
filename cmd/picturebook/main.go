// picturebook is a command-line application for creating a PDF file from a folder containing images.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aaronland/go-picturebook/application/commandline"
	"github.com/aaronland/go-slog/attr"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	ctx := context.Background()

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: commandline.LogLevel,
		ReplaceAttr: attr.EmojiLevelFunc(),
	})
	
	logger := slog.New(h)

	err := commandline.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to run picturebook application", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}
