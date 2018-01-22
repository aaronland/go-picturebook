package main

import (
	"flag"
	"fmt"
	"github.com/iand/salience"
	"github.com/straup/go-image-tools/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	width := flag.Int("width", 200, "...")
	height := flag.Int("height", 200, "...")

	flag.Parse()

	for _, path := range flag.Args() {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		im, format, err := util.DecodeImage(abs_path)

		if err != nil {
			log.Fatal(err)
		}

		cropped := salience.Crop(im, *width, *height)

		root := filepath.Dir(abs_path)
		fname := filepath.Base(abs_path)
		ext := filepath.Ext(abs_path)

		new_ext := fmt.Sprintf("-%s%s", "crop", ext)
		fname = strings.Replace(fname, ext, new_ext, -1)

		new_path := filepath.Join(root, fname)

		fh, err := os.Create(new_path)

		if err != nil {
			log.Fatal(err)
		}

		defer fh.Close()

		err = util.EncodeImage(cropped, format, fh)

		if err != nil {
			log.Fatal(err)
		}
	}

}
