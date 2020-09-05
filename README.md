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
Usage of ./bin/picturebook:
  -border float
    	The size of the border around images. (default 0.01)
  -caption string
    	A valid caption.Caption URI. Valid schemes are: cooperhewitt, filename, flickr, none, orthis
  -debug
    	DEPRECATED: Please use the -verbose flag instead.
  -dpi float
    	The DPI (dots per inch) resolution for your picturebook. (default 150)
  -exclude value
    	A valid regular expression to use for testing whether a file should be excluded from your picturebook. DEPRECATED: Please use -filter regexp://exclude/?pattern={REGULAR_EXPRESSION} instead.
  -filename string
    	The filename (path) for your picturebook. (default "picturebook.pdf")
  -fill-page
    	If necessary rotate image 90 degrees to use the most available page space.
  -filter value
    	A valid filter.Filter URI. Valid schemes are: any, cooperhewitt, flickr, orthis, regexp
  -height float
    	A custom width to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag. (default 11)
  -include value
    	A valid regular expression to use for testing whether a file should be included in your picturebook. DEPRECATED: Please use -filter regexp://include/?pattern={REGULAR_EXPRESSION} instead.
  -orientation string
    	The orientation of your picturebook. Valid orientations are: [please write me] (default "P")
  -pre-process value
    	DEPRECATED: Please use -process {PROCESS_NAME}:// instead.
  -process value
    	A valid process.Process URI. Valid schemes are: halftone, null, rotate
  -sort string
    	A valid sort.Sorter URI. Valid schemes are: orthis
   -size string
    	A common paper size to use for the size of your picturebook. Valid sizes are: [please write me] (default "letter")
  -target string
    	Valid targets are: cooperhewitt; flickr; orthis. If defined this flag will set the -filter and -caption flags accordingly. DEPRECATED: Please use specific -filter and -caption flags as needed.
  -verbose
    	Display verbose output as the picturebook is created.
  -width float
    	A custom height to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag. (default 8.5)
```

## Handlers

The `picturebook` application supports a number of "handlers" for customizing which images are included, how and whether they are transformed before inclusion and how to derive that image's caption.

### Captions

```
type Caption interface {
	Text(context.Context, string) (string, error)
}
```

The following schemes for caption handlers are supported:

#### cooperhewitt://

This handler will derive the title for a Cooper Hewitt collection object using data stored in a `index.json` file, alongside an image. The data in the file is expected to be the out of a call to the [cooperhewitt.shoebox.items.getInfo](https://collection.cooperhewitt.org/api/methods/cooperhewitt.shoebox.items.getInfo) API method.

_This handler will eventually be moved in to a separate `go-picturebook-cooperhewitt` package._

#### filename://

This handler will return the filename for a given path of an image.

_This is the default caption handler for picturebooks._

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

The following schemes for filter handlers are supported:

#### any://

Allow all images to be included.

_This is the default fitler handler for picturebooks._

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

#### halftone://

This handler will dither (halftone) an image before including it in your picturebook.

#### null://

This handler doesn't do anything to an image before including it in your picturebook.

_This is the default process handler for picturebooks._

#### rotate://

This handler will attempt to auto-rotate an image, based on any available EXIF `Orientation` data, before including it in your picturebook.

### Sorters

```
type Sorter interface {
	Sort(context.Context, []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error)
}
```

#### orthis://

This is really specific to [me and only me](https://aaronland.info/orthis) so you can ignore this for the time being.

_This handler will eventually be moved in to a separate `go-picturebook-orthis` package._

## See also

* https://github.com/jung-kurt/gofpdf
* https://github.com/aaronland/go-image-tools
