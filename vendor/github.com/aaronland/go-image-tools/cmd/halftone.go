package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-image-tools/halftone"
	"github.com/aaronland/go-image-tools/pixel"
	"github.com/aaronland/go-image-tools/util"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	opts := halftone.NewDefaultHalftoneOptions()

	mode := flag.String("mode", opts.Mode, "...")
	scale_factor := flag.Float64("scale-factor", opts.ScaleFactor, "...")

	flip := flag.Bool("flip", false, "...")

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

		opts.Mode = *mode
		opts.ScaleFactor = *scale_factor

		dithered, err := halftone.Halftone(im, opts)

		if err != nil {
			log.Fatal(err)
		}

		if *flip {

			wh := color.RGBA{
				R: uint8(255),
				G: uint8(255),
				B: uint8(255),
				A: uint8(255),
			}

			bl := color.RGBA{
				R: uint8(0),
				G: uint8(0),
				B: uint8(0),
				A: uint8(255),
			}

			wh2bl := pixel.ReplacePixelKey{
				Candidates:  []color.Color{wh},
				Replacement: bl,
			}

			bl2wh := pixel.ReplacePixelKey{
				Candidates:  []color.Color{bl},
				Replacement: wh,
			}

			f, err := pixel.MakeReplacePixelFunc(wh2bl, bl2wh)

			if err != nil {
				log.Fatal(err)
			}

			dithered, err = pixel.ProcessImage(dithered, f)

			if err != nil {
				log.Fatal(err)
			}
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
