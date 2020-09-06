package application

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/go-picturebook/caption"
	"github.com/aaronland/go-picturebook/filter"
	"github.com/aaronland/go-picturebook/process"
	"github.com/aaronland/go-picturebook/sort"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"os"
	"regexp"
	"strings"
)

type CommandLineApplication struct {
	Application
	flagset *flag.FlagSet
}

func CommandLineApplicationDefaultFlagSet() (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("picturebook")

	available_filters := strings.Join(filter.AvailableFilters(), ", ")
	available_filters = strings.ToLower(available_filters)

	available_captions := strings.Join(caption.AvailableCaptions(), ", ")
	available_captions = strings.ToLower(available_captions)

	available_processes := strings.Join(process.AvailableProcesses(), ", ")
	available_processes = strings.ToLower(available_processes)

	available_sorters := strings.Join(sort.AvailableSorters(), ", ")
	available_sorters = strings.ToLower(available_sorters)

	desc_filters := fmt.Sprintf("A valid filter.Filter URI. Valid schemes are: %s", available_filters)
	desc_captions := fmt.Sprintf("A valid caption.Caption URI. Valid schemes are: %s", available_captions)
	desc_processes := fmt.Sprintf("A valid process.Process URI. Valid schemes are: %s", available_processes)
	desc_sorters := fmt.Sprintf("A valid sort.Sorter URI. Valid schemes are: %s", available_sorters)

	fs.String("orientation", "P", "The orientation of your picturebook. Valid orientations are: [please write me]")
	fs.String("size", "letter", "A common paper size to use for the size of your picturebook. Valid sizes are: [please write me]")
	fs.Float64("width", 8.5, "A custom height to use as the size of your picturebook. Units are currently defined in inches. This fs.overrides the -size fs.")
	fs.Float64("height", 11, "A custom width to use as the size of your picturebook. Units are currently defined in inches. This fs.overrides the -size fs.")
	fs.Float64("dpi", 150, "The DPI (dots per inch) resolution for your picturebook.")
	fs.Float64("border", 0.01, "The size of the border around images.")

	fs.Bool("fill-page", false, "If necessary rotate image 90 degrees to use the most available page space.")

	fs.String("filename", "picturebook.pdf", "The filename (path) for your picturebook.")

	fs.Bool("verbose", false, "Display verbose output as the picturebook is created.")
	fs.Bool("debug", false, "DEPRECATED: Please use the -verbose fs.instead.")

	fs.String("caption", "", desc_captions)
	fs.String("sort", "", desc_sorters)

	fs.Bool("ocra-font", false, "Use an OCR-compatible font for captions.")

	var filter_uris multi.MultiString
	var process_uris multi.MultiString

	fs.Var(&filter_uris, "filter", desc_filters)
	fs.Var(&process_uris, "process", desc_processes)

	// Deprecated flags

	var preprocess_uris multi.MultiString
	var include multi.MultiRegexp
	var exclude multi.MultiRegexp

	fs.Var(&preprocess_uris, "pre-process", "DEPRECATED: Please use -process {PROCESS_NAME}:// instead.")
	fs.Var(&include, "include", "A valid regular expression to use for testing whether a file should be included in your picturebook. DEPRECATED: Please use -filter regexp://include/?pattern={REGULAR_EXPRESSION} instead.")
	fs.Var(&exclude, "exclude", "A valid regular expression to use for testing whether a file should be excluded from your picturebook. DEPRECATED: Please use -filter regexp://exclude/?pattern={REGULAR_EXPRESSION} instead.")

	fs.String("target", "", "Valid targets are: cooperhewitt; flickr; orthis. If defined this flag will set the -filter and -caption flags accordingly. DEPRECATED: Please use specific -filter and -caption flags as needed.")

	return fs, nil
}

func NewCommandLineApplication(ctx context.Context, fs *flag.FlagSet) (Application, error) {

	app := &CommandLineApplication{
		flagset: fs,
	}

	return app, nil
}

func (app *CommandLineApplication) Run(ctx context.Context) error {

	flagset.Parse(app.flagset)

	uri_re, err := regexp.Compile(`(?:[a-z0-9_]+):\/\/.*`)

	if err != nil {
		msg := fmt.Sprintf("Failed to compile URI regular expression, %v", err)
		return errors.New(msg)
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
			msg := fmt.Sprintf("Failed to assign filter '%s', %v", str_filter, err)
			return errors.New(msg)
		}

		if *caption_uri != "" {
			msg := fmt.Sprintf("Can not assign -caption using -target since -caption is already defined.")
			return errors.New(msg)
		}

		*caption_uri = str_caption
	}

	if len(preprocess_uris) > 0 {

		log.Println("WARNING The -pre-process flag is deprecated. Please use -process process://{PROCESS_NAME} flags instead.")

		for _, pr := range preprocess_uris {

			str_process := fmt.Sprintf("%s://", pr)
			err := process_uris.Set(str_process)

			if err != nil {
				msg := fmt.Sprintf("Failed to assign process '%s', %v", str_process, err)
				return errors.New(msg)
			}
		}
	}

	if len(include) > 0 {

		log.Println("WARNING The -include flag is deprecated. Please use -filter regexp://include?pattern=... flags instead.")

		for _, re := range include {

			str_filter := fmt.Sprintf("regexp://include?pattern=%s", re.String())
			err := filter_uris.Set(str_filter)

			if err != nil {
				msg := fmt.Sprintf("Failed to assign filter '%s', %v", str_filter, err)
				return errors.New(msg)
			}
		}
	}

	if len(exclude) > 0 {

		log.Println("WARNING The -exclude flag is deprecated. Please use -filter regexp://exclude?pattern=... flags instead.")

		for _, re := range exclude {

			str_filter := fmt.Sprintf("regexp://exclude?pattern=%s", re.String())
			err := filter_uris.Set(str_filter)

			if err != nil {
				msg := fmt.Sprintf("Failed to assign filter '%s', %v", str_filter, err)
				return errors.New(msg)
			}
		}
	}

	opts, err := picturebook.NewPictureBookDefaultOptions(ctx)

	if err != nil {
		msg := fmt.Sprintf("Failed to create default picturebook options, %v", err)
		return errors.New(msg)
	}

	opts.Orientation = *orientation
	opts.Size = *size
	opts.Width = *width
	opts.Height = *height
	opts.DPI = *dpi
	opts.Border = *border
	opts.FillPage = *fill_page
	opts.Verbose = *verbose
	opts.OCRAFont = *ocra_font

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
				msg := fmt.Sprintf("Failed to create filter '%s', %v", filter_uri, err)
				return errors.New(msg)
			}

			filters[idx] = f
		}

		multi, err := filter.NewMultiFilter(ctx, filters...)

		if err != nil {
			msg := fmt.Sprintf("Failed to create multi filter, %v", err)
			return errors.New(msg)
		}

		opts.Filter = multi
	}

	if len(process_uris) > 0 {

		processes := make([]process.Process, len(process_uris))

		for idx, process_uri := range process_uris {

			f, err := process.NewProcess(ctx, process_uri)

			if err != nil {
				msg := fmt.Sprintf("Failed to create process '%s', %v", process_uri, err)
				return errors.New(msg)
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

	if *sort_uri != "" {

		s, err := sort.NewSorter(ctx, *sort_uri)

		if err != nil {
			return err
		}

		opts.Sort = s
	}

	pb, err := picturebook.NewPictureBook(ctx, opts)

	if err != nil {
		msg := fmt.Sprintf("Failed to create new picturebook, %v", err)
		return errors.New(msg)
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
