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
  -border float
    	The size of the border around images. (default 0.01)
  -caption string
    	A valid caption.Caption URI. Valid schemes are: cooperhewitt, filename, flickr, none, orthis
  -debug
    	DEPRECATED: Please use the -verbose flag instead.
  -dpi float
    	The DPI (dots per inch) resolution for your picturebook. (default 150)
  -exclude value
    	A valid regular expression to use for testing whether a file should be excluded from your picturebook. DEPRECATED: Please use -filter regexp://exclude/?pattern={REGULAR_EXPRESSION} flag instead.
  -filename string
    	The filename (path) for your picturebook. (default "picturebook.pdf")
  -fill-page
    	If necessary rotate image 90 degrees to use the most available page space.
  -filter value
    	A valid filter.Filter URI. Valid schemes are: any, cooperhewitt, flickr, orthis, regexp
  -height float
    	A custom width to use as the size of your picturebook. Units are currently defined in inches. This fs.overrides the -size fs. (default 11)
  -include value
    	A valid regular expression to use for testing whether a file should be included in your picturebook. DEPRECATED: Please use -filter regexp://include/?pattern={REGULAR_EXPRESSION} flag instead.
  -ocra-font
    	Use an OCR-compatible font for captions.
  -orientation string
    	The orientation of your picturebook. Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively. (default "P")
  -pre-process value
    	DEPRECATED: Please use -process {PROCESS_NAME}:// flag instead.
  -process value
    	A valid process.Process URI. Valid schemes are: halftone, null, rotate
  -size string
    	A common paper size to use for the size of your picturebook. Valid sizes are: [please write me] (default "letter")
  -sort string
    	A valid sort.Sorter URI. Valid schemes are: modtime, orthis	
  -source-uri string
    	A valid GoCloud blob URI to specify where files should be read from. By default file:// URIs are supported.
  -target string
    	Valid targets are: cooperhewitt; flickr; orthis. If defined this flag will set the -filter and -caption flags accordingly. DEPRECATED: Please use specific -filter and -caption flags as needed.
  -target-uri string
    	A valid GoCloud blob URI to specify where your final picturebook PDF file should be written to. By default file:// URIs are supported.
  -verbose
    	Display verbose output as the picturebook is created.
  -width float
    	A custom height to use as the size of your picturebook. Units are currently defined in inches. This fs.overrides the -size fs. (default 8.5)
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
	Text(context.Context, string) (string, error)
}
```

For an example of how to create and register a custom `Caption` handler take a look at the code in [caption/filename.go](caption/filename.go).

The following schemes for caption handlers are supported by default:

#### cooperhewitt://

This handler will derive the title for a Cooper Hewitt collection object using data stored in a `index.json` file, alongside an image. The data in the file is expected to be the out of a call to the [cooperhewitt.shoebox.items.getInfo](https://collection.cooperhewitt.org/api/methods/cooperhewitt.shoebox.items.getInfo) API method.

_This handler will eventually be moved in to a separate `go-picturebook-cooperhewitt` package._

#### filename://

This handler will return the filename for a given path of an image.

#### flickr://

This handler will derive the title for a Flickr photo using data stored in a `{PHOTO_ID}_{SECRET}_i.json` file, alongside an image. The data in the file is expected to be the out of a call to the [flickr.photos.getInfo](https://www.flickr.com/services/api/flickr.photos.getInfo.html) API method.

_This handler will eventually be moved in to a separate `go-picturebook-flickr` package._

#### none://

The handler will return an empty string for all images.

#### orthis://

This is really specific to [me and only me](https://aaronland.info/orthis) so you can ignore this for the time being.

_This handler will eventually be moved in to a separate `go-picturebook-orthis` package._

### Filters

```
type Filter interface {
	Continue(context.Context, string) (bool, error)
}
```

For an example of how to create and register a custom `Filter` handler take a look at the code in [filter/regexp.go](filter/regexp.go).

The following schemes for filter handlers are supported by default:

#### any://

Allow all images to be included.

#### cooperhewitt://

This handler will ensure that only images whose filename matches `_b.jpg$` and that have a sibling `index.json` file are included.

_This handler will eventually be moved in to a separate `go-picturebook-cooperhewitt` package._

#### flickr://

This handler will ensure that only images whose filename matches `o_\.\.*$` are included.

_This handler will eventually be moved in to a separate `go-picturebook-flickr` package._

#### orthis://

This is really specific to [me and only me](https://aaronland.info/orthis) so you can ignore this for the time being.

Parameters

| Name | Value | Required |
| --- | --- | --- |
| year | A valid year | no |

_This handler will eventually be moved in to a separate `go-picturebook-orthis` package._

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
	Transform(context.Context, string) (string, error)
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
	Sort(context.Context, []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error)
}
```

For an example of how to create and register a custom `Sorter` handler take a look at the code in [sort/orthis.go](sort/orthis.go).

The following schemes for sorter handlers are supported by default:

#### modtime://

Sort images, in ascending order, by their modification times. If two or more images have the same modification they will sorted again by their file size.

#### orthis://

This is really specific to [me and only me](https://aaronland.info/orthis) so you can ignore this for the time being.

_This handler will eventually be moved in to a separate `go-picturebook-orthis` package._

## See also

* https://github.com/jung-kurt/gofpdf
* https://github.com/aaronland/go-image-tools
* https://github.com/aaronland/go-image-halftone
