package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-image-tools/flags"
	"github.com/aaronland/go-image-tools/pixel"
	"github.com/aaronland/go-image-tools/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var rgb flags.RGBAColor
	flag.Var(&rgb, "color", "...")

	var format = flag.String("format", "png", "...")

	flag.Parse()

	switch strings.ToUpper(*format) {
	case "PNG":
		// pass

		// this doesn't work yet...
	// case "GIF":
	// 	// pass

	default:
		log.Fatal("Invalid format for transparencies")
	}

	cb, err := pixel.MakeTransparentPixelFunc(rgb...)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		abs_path, err := filepath.Abs(path)

		im, err := pixel.ProcessPath(abs_path, cb)

		if err != nil {
			log.Fatal(err)
		}

		root := filepath.Dir(abs_path)
		fname := filepath.Base(abs_path)

		ext := filepath.Ext(fname)

		fname = strings.Replace(fname, ext, "", -1)

		fname = fmt.Sprintf("%s-tr.%s", fname, *format)
		new_path := filepath.Join(root, fname)

		fh, err := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatal(err)
		}

		err = util.EncodeImage(im, *format, fh)

		if err != nil {
			log.Fatal(err)
		}

		err = fh.Close()

		if err != nil {
			log.Fatal(err)
		}

		log.Println(new_path)
	}
}
