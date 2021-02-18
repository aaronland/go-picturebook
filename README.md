# go-picturebook

Create a PDF file (a "picturebook") from a folder (containing images).

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
  -caption string
    	A valid caption.Caption URI. Valid schemes are: exif, filename, none
  -debug
    	DEPRECATED: Please use the -verbose flag instead.
  -dpi float
    	The DPI (dots per inch) resolution for your picturebook. (default 150)
  -even-only
    	Only include images on even-numbered pages.
  -exclude value
    	A valid regular expression to use for testing whether a file should be excluded from your picturebook. DEPRECATED: Please use -filter regexp://exclude/?pattern={REGULAR_EXPRESSION} flag instead.
  -filename string
    	The filename (path) for your picturebook. (default "picturebook.pdf")
  -fill-page
    	If necessary rotate image 90 degrees to use the most available page space.
  -filter value
    	A valid filter.Filter URI. Valid schemes are: any, regexp
  -height float
    	A custom width to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag when used in combination with the -width flag.
  -include value
    	A valid regular expression to use for testing whether a file should be included in your picturebook. DEPRECATED: Please use -filter regexp://include/?pattern={REGULAR_EXPRESSION} flag instead.
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
  -ocra-font
    	Use an OCR-compatible font for captions.
  -odd-only
    	Only include images on odd-numbered pages.
  -orientation string
    	The orientation of your picturebook. Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively. (default "P")
  -pre-process value
    	DEPRECATED: Please use -process {PROCESS_NAME}:// flag instead.
  -process value
    	A valid process.Process URI. Valid schemes are: halftone, null, rotate
  -size string
    	A common paper size to use for the size of your picturebook. Valid sizes are: "A3", "A4", "A5", "Letter", "Legal", or "Tabloid". (default "letter")
  -sort string
    	A valid sort.Sorter URI. Valid schemes are: exif, modtime
  -source-uri string
    	A valid GoCloud blob URI to specify where files should be read from. By default file:// URIs are supported.
  -target string
    	Valid targets are: cooperhewitt; flickr; orthis. If defined this flag will set the -filter and -caption flags accordingly. DEPRECATED: Please use specific -filter and -caption flags as needed.
  -target-uri string
    	A valid GoCloud blob URI to specify where your final picturebook PDF file should be written to. By default file:// URIs are supported.
  -verbose
    	Display verbose output as the picturebook is created.
  -width float
    	A custom height to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag when used in combination with the -height flag.
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

Under the hood `picturebook` is using the [Go Cloud `blob` abstraction layer](https://gocloud.dev/howto/blob/) for files and file storage. By default only [local files](https://gocloud.dev/howto/blob/#local) (or `file://` URIs) are supported. If you need to support other sources or targets you will need to create your own custom `picturebook` tool.

In order to facilitate this all of the logic of the `picturebook` tool has been moved in to the [go-picturebook/application/commandline](application/commandline) package. For example here is how you would write your own custom tool with support for reading and writing files to an S3 bucket as well as the local filesystem.

```
package main

import (
	"context"
	"github.com/aaronland/go-picturebook/application/commandline"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"	
)

func main() {

	ctx := context.Background()

	fs, _ := commandline.DefaultFlagSet(ctx)
	app, _ := commandline.NewApplication(ctx, fs)

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

* https://github.com/aaronland/go-picturebook-cooperhewitt
* https://github.com/aaronland/go-picturebook-flickr
* https://github.com/jung-kurt/gofpdf
* https://github.com/aaronland/go-image-tools
* https://github.com/aaronland/go-image-halftone
* https://gocloud.dev/howto/blob/