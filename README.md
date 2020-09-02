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
./bin/picturebook -h
Usage of ./bin/picturebook:
  -border float
    	The size of the border around images. (default 0.01)
  -caption string
    	Valid filters are: cooperhewitt; default; flickr; orthis (default "default")
  -debug
    	...
  -dpi float
    	The DPI (dots per inch) resolution for your picturebook. (default 150)
  -exclude value
    	A valid regular expression to use for testing whether a file should be excluded from your picturebook.
  -filename string
    	The filename (path) for your picturebook. (default "picturebook.pdf")
  -filter string
    	Valid filters are: cooperhewitt; flickr; orthis
  -height float
    	A custom width to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag. (default 11)
  -include value
    	A valid regular expression to use for testing whether a file should be included in your picturebook.
  -orientation string
    	The orientation of your picturebook. Valid orientations are: [please write me] (default "P")
  -pre-process value
    	Valid processes are: rotate; halftone
  -size string
    	A common paper size to use for the size of your picturebook. Valid sizes are: [please write me] (default "letter")
  -target string
    	Valid targets are: cooperhewitt; flickr; orthis. If defined this flag will set the -filter and -caption flags accordingly.
  -width float
    	A custom height to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag. (default 8.5)
```

## Captions

### cooperhewitt://

### filename://

### flickr://

### orthis://

## Filters

### any://

### cooperhewitt://

### flickr://

### orthis://

Optional parameters are:

* `year=YYYY`

### regexp://exclude

### regexp://include

## Processes

### halftone://

### rotate://

## See also

* https://github.com/jung-kurt/gofpdf
* https://github.com/aaronland/go-image-tools
