package util

import (
	"bufio"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
)

// https://golang.org/src/image/decode_test.go

func DecodeImage(path string) (image.Image, string, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, "", err
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		return nil, "", err
	}

	defer fh.Close()

	return DecodeImageFromReader(fh)
}

func DecodeImageFromReader(fh io.Reader) (image.Image, string, error) {

	return image.Decode(bufio.NewReader(fh))
}
