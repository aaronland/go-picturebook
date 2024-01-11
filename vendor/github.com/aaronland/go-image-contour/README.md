# go-image-contour

Opinionated Go package for working with the `fogleman/contourmap` package.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-image-contour.svg)](https://pkg.go.dev/github.com/aaronland/go-image-contour)

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/contour cmd/contour/main.go
go build -mod vendor -ldflags="-s -w" -o bin/contour-svg cmd/contour-svg/main.go
```

### contour

Generate contour data derived from an image and draws the results to a new image.

```
$> ./bin/contour -h
  -n int
    	The number of iterations used to generate contours. (default 12)
  -scale float
    	The scale of the final output relative to the input image. (default 1)
  -source-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are read from. (default "file:///")
  -target-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are written to. (default "file:///")
  -transformation-uri transform.Transformation
    	Zero or more additional transform.Transformation URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).
```

For example:

```
$> ./bin/contour -n 3 fixtures/tokyo.jpg
```

Will transform this:

![](fixtures/tokyo.jpg)

In to this:

![](fixtures/tokyo-contour-3.jpg)

### contour-svg

Generate contour data derived from an image and write the results as SVG path elements to a new file.

```
$> ./bin/contour-svg  -h
  -n int
    	The number of iterations used to generate contours. (default 12)
  -scale float
    	The scale of the final output relative to the input image. (default 1)
  -source-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are read from. (default "file:///")
  -target-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are written to. (default "file:///")
  -transformation-uri transform.Transformation
    	Zero or more additional transform.Transformation URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).
```

## See also

* https://github.com/aaronland/go-image
* https://github.com/fogleman/contourmap
* https://github.com/fogleman/gg