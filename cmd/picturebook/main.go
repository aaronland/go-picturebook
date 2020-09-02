package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-picturebook"
	"github.com/aaronland/go-picturebook/caption"
	"github.com/aaronland/go-picturebook/filter"
	"github.com/aaronland/go-picturebook/process"
	"github.com/sfomuseum/go-flags/multi"
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

	var filename = flag.String("filename", "picturebook.pdf", "The filename (path) for your picturebook.")

	var debug = flag.Bool("debug", false, "...")

	var caption_uri = flag.String("caption", "", "...")
	var filter_uris multi.MultiString
	var process_uris multi.MultiString

	flag.Var(&filter_uris, "filter", "...")
	flag.Var(&process_uris, "process", "...")

	/*
		var include multi.MultiRegexp
		var exclude multi.MultiRegexp

		flag.Var(&include, "include", "A valid regular expression to use for testing whether a file should be included in your picturebook.")
		flag.Var(&exclude, "exclude", "A valid regular expression to use for testing whether a file should be excluded from your picturebook.")
	*/

	flag.Parse()

	ctx := context.Background()

	opts, err := picturebook.NewPictureBookDefaultOptions(ctx)

	if err != nil {
		log.Fatalf("Failed to create default picturebook options, %v", err)
	}

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

	if len(filter_uris) > 0 {

		filters := make([]filter.Filter, len(filter_uris))

		for idx, filter_uri := range filter_uris {

			f, err := filter.NewFilter(ctx, filter_uri)

			if err != nil {
				log.Fatalf("Failed to create filter '%s', %v", filter_uri, err)
			}

			filters[idx] = f
		}

		multi, err := filter.NewMultiFilter(ctx, filters...)

		if err != nil {
			log.Fatalf("Failed to create multi filter, %v", err)
		}

		opts.Filter = multi
	}

	if len(process_uris) > 0 {

		processes := make([]process.Process, len(process_uris))

		for idx, process_uri := range process_uris {

			f, err := process.NewProcess(ctx, process_uri)

			if err != nil {
				log.Fatalf("Failed to create process '%s', %v", process_uri, err)
			}

			processes[idx] = f
		}

		multi, err := process.NewMultiProcess(ctx, processes...)

		if err != nil {
			log.Fatalf("Failed to create multi process, %v", err)
		}

		opts.PreProcess = multi
	}

	if *caption_uri != "" {

		c, err := caption.NewCaption(ctx, *caption_uri)

		if err != nil {
			log.Fatal(err)
		}

		opts.Caption = c
	}

	pb, err := picturebook.NewPictureBook(ctx, opts)

	if err != nil {
		log.Fatalf("Failed to create new picturebook, %v", err)
	}

	sources := flag.Args()

	err = pb.AddPictures(ctx, sources)

	if err != nil {
		return err
	}

	err = pb.Save(ctx, *filename)

	if err != nil {
		return err
	}

	return nil
}
