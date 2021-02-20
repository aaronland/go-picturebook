package picture

import (
	"gocloud.dev/blob"
)

type PictureBookPicture struct {
	Source  string
	Path    string
	Caption string
	Bucket *blob.Bucket
}
