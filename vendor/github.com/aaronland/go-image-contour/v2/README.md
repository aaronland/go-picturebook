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

## WASM (WebAssembly)

### Building

This package comes with a pre-compile WASM binary which is stored in [www/wasm/contour.wasm](www/wasm). If you need or want to rebuild the binary the easiest way is to use the handy `wasmjs` Makefile target. For example:

```
$> make wasmjs
GOOS=js GOARCH=wasm \
		go build -mod vendor -ldflags="-s -w" -tags wasmjs \
		-o www/wasm/contour.wasm \
		cmd/contour-wasm-js/main.go
```

### Usage

The `contour.wasm` binary exposes two functions: `contour_image` and `contour_svg`. Each take the same arguments: A base64-encoded image to contour and the number of iterations used to generate that contour. Both functions return a JavaScript `Promise`. The `contour_image` function returns a base64-encoded image and the `contour_svg` returns a string containing an SVG document. For example:

```
const im_b64 = "base64-encoded image here...";
const iterations = 6;

contour_svg(im_b64, iterations).then((rsp) => {
	// Display SVG here
}).catch((err) => {
	console.error("Failed to contour image", err);
});	
```

Loading the `contour.wasm` binary in a JavaScript application is expected to be handled by the application itself. I like to use the [sfomuseum/js-sfomuseum-golang-wasm](https://github.com/sfomuseum/js-sfomuseum-golang-wasm) but there is no requirement to do so.

### Example

There is an example web application demonstrating the use of the `contour.wasm` binary in the [www](www) directory. To run the example you'll see to serve the application from a web browser. I like to the use the [aaronland/go-http-fileserver](https://github.com/aaronland/go-http-fileserver?tab=readme-ov-file#fileserver) tool, mostly because I wrote it, but any old web server will do. For example:

```
$> fileserver -root www/
2025/07/13 09:25:25 Serving www/ and listening for requests on http://localhost:8080
```

And then when you open your web browser to `http://localhost:8080` you'll see this:

![](docs/images/go-image-contour-wasm-launch.png)

You can upload individual images and have them converted in to SVG contours. For example:

![](docs/images/go-image-contour-wasm-image.png)

You can also start a live video feed and generate SVG contours from individual video stills. For example:

![](docs/images/go-image-contour-wasm-video.png)

That's really all the example application does. Saving or otherwise working with contour-ed images is left as an exercise to the reader (or a pull request).

## See also

* https://github.com/aaronland/go-image
* https://github.com/fogleman/contourmap
* https://github.com/fogleman/gg