# go-picturebook

Too soon. Move along.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

### picturebook

```
./bin/picturebook -h
Usage of ./bin/picturebook:
  -border float
    	... (default 0.01)
  -caption string
    	... (default "default")
  -debug
    	...
  -dpi float
    	... (default 150)
  -exclude value
    	...
  -filename string
    	... (default "picturebook.pdf")
  -filter string
    	...
  -height float
    	... (default 11)
  -include value
    	...
  -orientation string
    	... (default "P")
  -pre-process value
    	...
  -size string
    	... (default "letter")
  -target string
    	...
  -width float
    	... (default 8.5)
```

## See also

* https://github.com/aaronland/go-image-tools
