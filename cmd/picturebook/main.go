package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-picturebook"
	"github.com/aaronland/go-picturebook/caption"
	"github.com/aaronland/go-picturebook/filter"
	"github.com/aaronland/go-picturebook/process"
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"os"
	"regexp"
)

func main() {

	err := Picturebook()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func Picturebook() error {

	// available_filters := filter.

	var orientation = flag.String("orientation", "P", "The orientation of your picturebook. Valid orientations are: [please write me]")
	var size = flag.String("size", "letter", "A common paper size to use for the size of your picturebook. Valid sizes are: [please write me]")
	var width = flag.Float64("width", 8.5, "A custom height to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag.")
	var height = flag.Float64("height", 11, "A custom width to use as the size of your picturebook. Units are currently defined in inches. This flag overrides the -size flag.")
	var dpi = flag.Float64("dpi", 150, "The DPI (dots per inch) resolution for your picturebook.")
	var border = flag.Float64("border", 0.01, "The size of the border around images.")

	var filename = flag.String("filename", "picturebook.pdf", "The filename (path) for your picturebook.")

	var verbose = flag.Bool("verbose", false, "Display verbose output as the picturebook is created.")
	var debug = flag.Bool("debug", false, "DEPRECATED: Please use the -verbose flag instead.")

	var caption_uri = flag.String("caption", "", "...")
	var filter_uris multi.MultiString
	var process_uris multi.MultiString

	flag.Var(&filter_uris, "filter", "...")
	flag.Var(&process_uris, "process", "...")

	// Deprecated flags

	var preprocess_uris multi.MultiString
	var include multi.MultiRegexp
	var exclude multi.MultiRegexp

	flag.Var(&preprocess_uris, "pre-process", "DEPRECATED: Please use -process process://{PROCESS_NAME} instead.")
	flag.Var(&include, "include", "A valid regular expression to use for testing whether a file should be included in your picturebook. DEPRECATED: Please use -filter regexp://include/?pattern={REGULAR_EXPRESSION} instead.")
	flag.Var(&exclude, "exclude", "A valid regular expression to use for testing whether a file should be excluded from your picturebook. DEPRECATED: Please use -filter regexp://exclude/?pattern={REGULAR_EXPRESSION} instead.")

	var target = flag.String("target", "", "Valid targets are: cooperhewitt; flickr; orthis. If defined this flag will set the -filter and -caption flags accordingly. DEPRECATED: Please use specific -filter and -caption flags as needed.")

	flag.Parse()

	ctx := context.Background()

	uri_re, err := regexp.Compile(`(?:[a-z0-9_]+):\/\/.*`)

	if err != nil {
		log.Fatalf("Failed to compile URI regular expression, %v", err)
	}

	if *debug {

		log.Println("WARNING The -debug flag is deprecated. Please use the -verbose flag instead.")
		*verbose = *debug
	}

	if *target != "" {

		log.Println("WARNING The -target flag is deprecated. Please use specific -filter and -caption flags as needed.")

		str_filter := fmt.Sprintf("%s://", *target)
		str_caption := fmt.Sprintf("%s://", *target)

		err := filter_uris.Set(str_filter)

		if err != nil {
			log.Fatalf("Failed to assign filter '%s', %v", str_filter, err)
		}

		if *caption_uri != "" {
			log.Fatalf("Can not assign -caption using -target since -caption is already defined.")
		}

		*caption_uri = str_caption
	}

	if len(preprocess_uris) > 0 {

		log.Println("WARNING The -pre-process flag is deprecated. Please use -process process://{PROCESS_NAME} flags instead.")

		for _, pr := range preprocess_uris {

			str_process := fmt.Sprintf("%s://", pr)
			err := process_uris.Set(str_process)

			if err != nil {
				log.Fatalf("Failed to assign process '%s', %v", str_process, err)
			}
		}
	}

	if len(include) > 0 {

		log.Println("WARNING The -include flag is deprecated. Please use -filter regexp://include?pattern=... flags instead.")
		for _, re := range include {

			str_filter := fmt.Sprintf("regexp://include?pattern=%s", re.String())
			err := filter_uris.Set(str_filter)

			if err != nil {
				log.Fatalf("Failed to assign filter '%s', %v", str_filter, err)
			}
		}
	}

	if len(exclude) > 0 {

		log.Println("WARNING The -exclude flag is deprecated. Please use -filter regexp://exclude?pattern=... flags instead.")

		for _, re := range exclude {

			str_filter := fmt.Sprintf("regexp://exclude?pattern=%s", re.String())
			err := filter_uris.Set(str_filter)

			if err != nil {
				log.Fatalf("Failed to assign filter '%s', %v", str_filter, err)
			}
		}
	}

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
	opts.Verbose = *verbose

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

			if !uri_re.MatchString(filter_uri) {
				filter_uri = fmt.Sprintf("%s://", filter_uri)
			}

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

		if !uri_re.MatchString(*caption_uri) {
			*caption_uri = fmt.Sprintf("%s://", *caption_uri)
		}

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
