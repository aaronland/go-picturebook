// package tempfile provides methods for working with temporary files used to create a picturebook.
package tempfile

import (
	"context"
	"fmt"
	"image"
	"uuid"
	
	"github.com/aaronland/go-image/v2/encode"
	"github.com/aaronland/go-picturebook/bucket"
)

// TempFileWithImage will write a new JPEG file in 'bucket' derived from 'im'. The return values are the
// filename of the temporary file, its image format and any errors produced during writing.
func TempFileWithImage(ctx context.Context, bucket bucket.Bucket, im image.Image) (string, string, error) {

	id := uuid.New()

	fname := fmt.Sprintf("picturebook-%s.jpg", id.String())

	wr, err := bucket.NewWriter(ctx, fname, nil)

	if err != nil {
		return "", "", fmt.Errorf("Failed to create new writer for temp file, %w", err)
	}

	err = encode.EncodeJPEG(ctx, wr, im, nil, nil)

	if err != nil {
		return "", "", fmt.Errorf("Failed to encode temp file as JPEG, %w", err)
	}

	err = wr.Close()

	if err != nil {
		return "", "", fmt.Errorf("Failed to close writer for temp file, %w", err)
	}

	return fname, "jpeg", nil
}
