# go-image-halftone

Go package implementing the `aaronland/go-image/v2/transform.Transformation` for interface halftone (dithering) processing of images.

## Documentation

Documentation is incomplete at this time.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/halftone cmd/halftone/main.go
```

### halftone

Apply a halftone (dithering) process to one or more images.

```
$> ./bin/halftone -h
Apply a halftone (dithering) process to one or more images.
Usage:
	./bin/halftone uri(N) uri(N)
  -preserve-exif
    	Copy EXIF data from source image final target image.
  -process string
    	The halftone process to use. (default "atkinson")
  -rotate
    	Automatically rotate based on EXIF orientation. This does NOT update any of the original EXIF data with one exception: If the -rotate flag is true OR the original image of type HEIC then the EXIF "Orientation" tag is re-written to be "1". (default true)
  -scale-factor int
    	The scale factor to use for the halftone process. (default 2)
  -source-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are read from. (default "file:///")
  -target-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are written to. (default "file:///")
  -transformation-uri transform.Transformation
    	Zero or more additional transform.Transformation URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).
```

#### Example

```
$> ./bin/halftone -scale-factor 2 -process atkinson ./big-fish-001.jpg
```

## See also

* https://github.com/aaronland/go-image
