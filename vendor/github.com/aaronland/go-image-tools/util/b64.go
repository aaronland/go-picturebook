package util

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"image"
)

func Base64EncodeImage(im image.Image, format string) (string, error) {

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	err := EncodeImage(im, format, wr)

	if err != nil {
		return "", err
	}

	wr.Flush()

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
