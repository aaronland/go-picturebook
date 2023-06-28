// package picture provides methods and data structures for internal representations of images.
package picture

import (
	"gocloud.dev/blob"
)

// type PictureBookPicture provides a struct containing details about an image to be added to a picturebook.
type PictureBookPicture struct {
	// The original (relative) path of the image being added to a picturebook
	Source string
	// The (relative) path of the final image to add to a picturebook
	Path string
	// The caption associated with the image being added to a picturebook
	Caption string
	// The long-form (or at least longer than a caption) text associated with the image being added to a picturebook
	Text string
	// The `blob.Bucket` instance where the image (Path) being added to a picturebook is stored.
	Bucket *blob.Bucket
	// The path of any temporary file that has been created in the process of adding an image to a picturebook
	TempFile string
}
