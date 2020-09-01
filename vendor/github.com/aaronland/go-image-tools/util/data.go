package util

import (
	"fmt"
	"image"
)

// this works but not in templates for some reason - the b64 data is
// always truncated... (20181024/straup)

func ImageToDataURL(im image.Image, format string) string {

	b64, err := Base64EncodeImage(im, format)

	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	return fmt.Sprintf("data:image/%s;base64,%s", format, b64)
}
