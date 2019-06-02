# go-picturebook

Create a PDF file (a "picturebook") from a folder (containing images).

## Install

You will need to have both `Go` (specifically version [1.12](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make tools
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

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

## Functions

### Caption functions

#### Cooper Hewitt (shoebox) caption functions

_Please write me_

#### Filename caption functions

_Please write me_

#### Filename and parent caption functions

_Please write me_

#### Flickr caption functions

_Please write me_

#### Or This caption functions

_Please write me_

### Filter functions

#### Cooper Hewitt (shoebox) filter functions

Only include files ending in `_b.jpg`

#### Flickr filter functions

Only include files matching `o_\.\.*$`

#### Or This filter functions

Only include files ending in `-or-this.jpg`

### Pre-process functions

#### Rotate preprocess functions

_Please write me_

#### Halftone preprocess functions

_Please write me_

## See also

* https://github.com/jung-kurt/gofpdf
* https://github.com/aaronland/go-image-tools
