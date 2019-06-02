package main

import (
	"errors"
	"flag"
	"github.com/aaronland/picturebook"
	"github.com/aaronland/picturebook/flags"
	"github.com/aaronland/picturebook/functions"
	"log"
	"os"
)

func main() {

	err := Picturebook()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func Picturebook() error {

	var orientation = flag.String("orientation", "P", "The orientation of your picturebook. Valid orientations are: [please write me]")
	var size = flag.String("size", "letter", "A common paper size to use for the size of your picturebook. Valid sizes are: [please write me]")
	var width = flag.Float64("width", 8.5, "A custom height to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag.")
	var height = flag.Float64("height", 11, "A custom width to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag.")
	var dpi = flag.Float64("dpi", 150, "The DPI (dots per inch) resolution for your picturebook.")
	var border = flag.Float64("border", 0.01, "The size of the border around images.")
	var caption = flag.String("caption", "default", "Valid filters are: cooperhewitt; default; flickr; orthis")
	var filter = flag.String("filter", "", "Valid filters are: cooperhewitt; flickr; orthis")
	var filename = flag.String("filename", "picturebook.pdf", "The filename (path) for your picturebook.")
	var target = flag.String("target", "", "Valid targets are: cooperhewitt; flickr; orthis. If defined this flag will set the -filter and -caption flags accordingly.")
	var debug = flag.Bool("debug", false, "...")

	var include flags.RegexpFlag
	var exclude flags.RegexpFlag
	var preprocess flags.PreProcessFlag

	flag.Var(&include, "include", "A valid regular expression to use for testing whether a file should be included in your picturebook.")
	flag.Var(&exclude, "exclude", "A valid regular expression to use for testing whether a file should be excluded from your picturebook.")
	flag.Var(&preprocess, "pre-process", "Valid processes are: rotate; halftone")

	flag.Parse()

	switch *target {
	case "":
		// pass
	case "cooperhewitt":
		*caption = *target
		*filter = *target
	case "flickr":
		*caption = *target
		*filter = *target
	case "orthis":
		*caption = *target
		*filter = *target
	default:
		log.Fatal("Unknown or invalid target")
	}

	opts := picturebook.NewPictureBookDefaultOptions()
	opts.Orientation = *orientation
	opts.Size = *size
	opts.Width = *width
	opts.Height = *height
	opts.DPI = *dpi
	opts.Border = *border
	opts.Debug = *debug

	processed := make([]string, 0)

	defer func() {
		for _, p := range processed {

			go func(p string) {

				_, err := os.Stat(p)

				if os.IsNotExist(err) {
					return
				}

				log.Println("WOULD REMOVE", p)
				// os.Remove(p)
			}(p)
		}
	}()

	filter_func := func(path string) (bool, error) {

		if *filter != "" {

			f, err := functions.PictureBookFilterFuncFromString(*filter)

			if err != nil {
				return false, err
			}

			return f(path)
		} else {

			for _, pat := range include {

				if !pat.MatchString(path) {
					return false, nil
				}
			}

			for _, pat := range exclude {

				if pat.MatchString(path) {
					return false, nil
				}
			}

			return true, nil
		}
	}

	prep := func(path string) (string, error) {

		final := path

		for _, proc := range preprocess {

			switch proc {

			case "rotate":

				processed_path, err := functions.RotatePreProcessFunc(final)

				if err != nil {
					return "", err
				}

				if processed_path == "" {
					continue
				}

				processed = append(processed, processed_path)
				final = processed_path

			case "halftone":

				processed_path, err := functions.HalftonePreProcessFunc(final)

				if err != nil {
					return "", err
				}

				if processed_path == "" {
					continue
				}

				processed = append(processed, processed_path)
				final = processed_path

			default:
				return "", errors.New("Invalid or unsupported process")
			}
		}

		return final, nil
	}

	capt, err := functions.PictureBookCaptionFuncFromString(*caption)

	if err != nil {
		log.Fatal(err)
	}

	opts.Filter = filter_func
	opts.PreProcess = prep
	opts.Caption = capt

	pb, err := picturebook.NewPictureBook(opts)

	if err != nil {
		return err
	}

	sources := flag.Args()

	err = pb.AddPictures(sources)

	if err != nil {
		return err
	}

	err = pb.Save(*filename)

	if err != nil {
		return err
	}

	return nil
}
