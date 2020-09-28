package tempfile

import (
	"context"
	"fmt"
	"github.com/aaronland/go-image-tools/util"
	"github.com/google/uuid"
	"gocloud.dev/blob"
	"image"
)

func TempFileWithImage(ctx context.Context, bucket *blob.Bucket, im image.Image) (string, string, error) {

	id, err := uuid.NewUUID()

	if err != nil {
		return "", "", err
	}

	fname := fmt.Sprintf("picturebook-%s.jpg", id.String())

	wr, err := bucket.NewWriter(ctx, fname, nil)

	if err != nil {
		return "", "", nil
	}

	err = util.EncodeImage(im, "jpeg", wr)

	if err != nil {
		return "", "", err
	}

	err = wr.Close()

	if err != nil {
		return "", "", err
	}

	return fname, "jpeg", nil
}
