package main

import (
	"flag"
	"fmt"
	"github.com/straup/go-image-tools/halftone"
	"github.com/straup/go-image-tools/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	mode := flag.String("mode", "atkinson", "...")
	scale_factor := flag.Float64("scale-factor", 2.0, "...")

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

		opts := halftone.NewDefaultHalftoneOptions()
		opts.Mode = *mode
		opts.ScaleFactor = *scale_factor

		dithered, err := halftone.Halftone(im, opts)

		if err != nil {
			log.Fatal(err)
		}

		root := filepath.Dir(abs_path)
		fname := filepath.Base(abs_path)
		ext := filepath.Ext(abs_path)

		new_ext := fmt.Sprintf("-%s%s", *mode, ext)
		fname = strings.Replace(fname, ext, new_ext, -1)

		new_path := filepath.Join(root, fname)

		fh, err := os.Create(new_path)

		if err != nil {
			log.Fatal(err)
		}

		err = util.EncodeImage(dithered, format, fh)

		if err != nil {
			log.Fatal(err)
		}

	}
}
