# go-picturebook

Create a PDF file (a "picturebook") from a folder (containing images).

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-picturebook.svg)](https://pkg.go.dev/github.com/aaronland/go-picturebook)

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> make cli
go build -mod vendor -o bin/picturebook cmd/picturebook/main.go
```

### picturebook

Create a PDF file (a "picturebook") from a folder (containing images).

```
$> ./bin/picturebook -h
  -bleed float
    	An additional bleed area to add (on all four sides) to the size of your picturebook.
  -border float
    	The size of the border around images. (default 0.01)
  -caption value
    	Zero or more valid caption.Caption URIs. Valid schemes are: exif://, filename://, json://, modtime://, multi://, none://.
  -dpi float
    	The DPI (dots per inch) resolution for your picturebook. (default 150)
  -even-only
    	Only include images on even-numbered pages.
  -filename string
    	The filename (path) for your picturebook. (default "picturebook.pdf")
  -fill-page
    	If necessary rotate image 90 degrees to use the most available page space. Note that any '-process' flags involving colour space manipulation will automatically be applied to images after they have been rotated.
  -filter value
    	A valid filter.Filter URI. Valid schemes are: any://, regexp://.
  -height float
    	A custom width to use as the size of your picturebook. Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -width flag.
  -margin float
    	The margin around all sides of a page. If non-zero this value will be used to populate all the other -margin-(N) flags.
  -margin-bottom float
    	The margin around the bottom of each page. (default 1)
  -margin-left float
    	The margin around the left-hand side of each page. (default 1)
  -margin-right float
    	The margin around the right-hand side of each page. (default 1)
  -margin-top float
    	The margin around the top of each page. (default 1)
  -max-pages int
    	An optional value to indicate that a picturebook should not exceed this number of pages
  -ocra-font
    	Use an OCR-compatible font for captions.
  -odd-only
    	Only include images on odd-numbered pages.
  -orientation string
    	The orientation of your picturebook. Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively. (default "P")
  -process value
    	A valid process.Process URI. Valid schemes are: colorspace://, colourspace://, contour://, halftone://, null://, rotate://.
  -size string
    	A common paper size to use for the size of your picturebook. Valid sizes are: "a3", "a4", "a5", "letter", "legal", or "tabloid". (default "letter")
  -sort string
    	A valid sort.Sorter URI. Valid schemes are: exif://, modtime://.
  -source-uri string
    	A valid GoCloud blob URI to specify where files should be read from. Available schemes are: file://. If no URI scheme is included then the file:// scheme is assumed.
  -target-uri string
    	A valid GoCloud blob URI to specify where files should be read from. Available schemes are: file://. If no URI scheme is included then the file:// scheme is assumed. If empty then the code will try to use the operating system's 'current working directory' where applicable.
  -text string
    	A valid text.Text URI. Valid schemes are: json://.
  -tmpfile-uri string
    	A valid GoCloud blob URI to specify where files should be read from. Available schemes are: file://. If no URI scheme is included then the file:// scheme is assumed. If empty the operating system's temporary directory will be used.
  -units string
    	The unit of measurement to apply to the -height and -width flags. Valid options are inches, millimeters, centimeters (default "inches")
  -verbose
    	Display verbose output as the picturebook is created.
  -width float
    	A custom height to use as the size of your picturebook. Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -height flag.
```

For example:

```
$> ./bin/picturebook \
	-source-uri file:///PATH/TO/go-picturebook/example \
	-target-uri file:///PATH/TO/go-picturebook/example \
	images
```

As a convenience if no paths (to folders containing images) are passed to the `picturebook` tool it will be assumed that images are found in the folder defined by the `-source-uri` flag. For example this command is functionally equivalent to the command above:

```
$> ./bin/picturebook \
	-source-uri file:///PATH/TO/go-picturebook/example/images \
	-target-uri file:///PATH/TO/go-picturebook/example \
```

For a complete example, including a PDF file produced by the `picturebook` tool, have a look in the [example](example) folder.

### Source and target URIs

Under the hood `picturebook` is using its own [Bucket](bucket/bucket.go) abstraction layer for files and file storage.

```
type Bucket interface {
	// GatherPictures returns an iterator for listing Picturebook images URIs that can passed to the (bucket implementation's) `NewReader` method.
	GatherPictures(context.Context, ...string) iter.Seq2[string, error]
	// NewReader returns an `io.ReadSeekCloser` instance for a record in the bucket.
	NewReader(context.Context, string, any) (io.ReadSeekCloser, error)
	// NewWriter returns an `io.WriterCloser` instance for writing a record to the bucket.
	NewWriter(context.Context, string, any) (io.WriteCloser, error)
	// Delete removed a record from the bucket.
	Delete(context.Context, string) error
	// Attributes returns an `Attributes` struct for a record in the bucket.
	Attributes(context.Context, string) (*Attributes, error)
	// Close signals the implementation to wrap things up (internally).
	Close() error
}
```

`Bucket` implements a simple interface for reading and writing Picturebook images to and from different storage implementations. It is modeled after the [gocloud.dev/blob.Bucket](https://gocloud.dev/howto/blob/) interface which is what this package used to use. This simplified interface reflects the limited methods from the original interface that were used. The goal is to make it easier to implement a variety of Picturebook "sources" (or buckets) without having to implement the entirety of the `gocloud.dev/blob.Bucket`

By default only [local files](https://gocloud.dev/howto/blob/#local) (or `file://` URIs) are supported. If you need to support other sources or targets you will need to create your own custom `picturebook` tool and add the relevant `import` statements.

In order to facilitate this all of the logic of the `picturebook` tool has been moved in to the [go-picturebook/application/commandline](application/commandline) package. For example here is how you would write your own custom tool with support for reading and writing files to an S3 bucket as well as the local filesystem.

```
package main

import (
	"context"
	
	"github.com/aaronland/go-picturebook/app/picturebook"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"	
)

func main() {
	ctx := context.Background()
	app.Run(ctx)
}
```

_Error handling has been omitted for the sake of brevity._

## Handlers

The `picturebook` application supports a number of "handlers" for customizing which images are included, how and whether they are transformed before inclusion and how to derive that image's caption.

### Captions

```
type Caption interface {
	Text(context.Context, *blob.Bucket, string) (string, error)
}
```

For an example of how to create and register a custom `Caption` handler take a look at the code in [caption/filename.go](caption/filename.go).

The following schemes for caption handlers are supported by default:

#### exif://?property={PROPERTY}

The handler will eventually return the value of a given EXIF property. As of this writing it will only return the value of the EXIF `DateTime` property.

If EXIF data is not present or can be loaded the handler will return an empty string.

Parameters

| Name | Value | Required |
| --- | --- | --- |
| property | "datetime" | yes |

#### filename://

This handler will return the filename for a given path of an image.

#### json://{PATH/TO/CAPTIONS.json}

This handler will assign captions derived from a JSON file that can be read from the local disk. The JSON file is expected to be a dictionary whose keys are the filenames of the images being included in the picturebook and whose values are a list of strings to use as caption text.

#### modtime://

This handler will assign captions derived from an image's modification time.

#### none://

The handler will return an empty string for all images.

### Filters

```
type Filter interface {
	Continue(context.Context, *blob.Bucket, string) (bool, error)
}
```

For an example of how to create and register a custom `Filter` handler take a look at the code in [filter/regexp.go](filter/regexp.go).

The following schemes for filter handlers are supported by default:

#### any://

Allow all images to be included.

#### regexp://exclude

This handler will exclude any images whose path matches the specified regular expression.

Parameters

| Name | Value | Required |
| --- | --- | --- |
| pattern | A valid Go language regular expresssion | yes |

#### regexp://include

This handler will include only those images whose path matches the specified regular expression.

Parameters

| Name | Value | Required |
| --- | --- | --- |
| pattern | A valid Go language regular expresssion | yes |

### Processes

```
type Process interface {
	Transform(context.Context, *blob.Bucket, string) (string, error)
}
```

For an example of how to create and register a custom `Process` handler take a look at the code in [process/halftone.go](process/halftone.go).

The following schemes for process handlers are supported by default:

#### colourspace://

This handler will map all the pixels in an image to a given colour space (Apple's Display P3, Adobe RGB) before including it in your picturebook. URIs should take the form of `colourspace://{{LABEL}.

##### Labels

| Name | Description |
| --- | --- |
| adobergb | Convert all pixels in to the Adobe RGB colour space |
| displayp3 | Convert all pixels in to Apple's Display P3 colour space |

#### colorspace://

This handler will map all the pixels in an image to a given colour space (Apple's Display P3, Adobe RGB) before including it in your picturebook. URIs should take the form of `colorspace://{{LABEL}.

##### Labels

| Name | Description |
| --- | --- |
| adobergb | Convert all pixels in to the Adobe RGB colour space |
| displayp3 | Convert all pixels in to Apple's Display P3 colour space |

#### contour://

This handler will convert an image into a series of black and white "contour" lines using the [fogleman/contourmap](https://github.com/fogleman/contourmap) package. URIs should take the form of `contour://?{PARAMETERS}`.

##### Parameters

| Name | Value | Required | Default |
| --- | --- | --- | --- |
| iterations | The number of iterations to perform during the contour process | no | 12 |
| scale | The scale of the final contoured image | no | 1.0 |

#### halftone://

This handler will dither (halftone) an image before including it in your picturebook.

#### null://

This handler doesn't do anything to an image before including it in your picturebook.

#### rotate://

This handler will attempt to auto-rotate an image, based on any available EXIF `Orientation` data, before including it in your picturebook.

### Sorters

```
type Sorter interface {
	Sort(context.Context, *blob.Bucket, []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error)
}
```

For an example of how to create and register a custom `Sorter` handler take a look at the code in [sort/orthis.go](sort/orthis.go).

The following schemes for sorter handlers are supported by default:

#### exif://

Sort images, in ascending order, by their EXIF `DateTime` property. If EXIF data is not present or can not be loaded the image's modification time will be used. If two or more images have the same modification they will sorted again by their file size.

#### modtime://

Sort images, in ascending order, by their modification times. If two or more images have the same modification they will sorted again by their file size.

## See also

* https://github.com/aaronland/go-picturebook-flickr
* https://github.com/go-pdf/fpdf
* https://github.com/aaronland/go-image
* https://github.com/aaronland/go-image-halftone
* https://github.com/aaronland/go-image-contour
* https://gocloud.dev/howto/blob/