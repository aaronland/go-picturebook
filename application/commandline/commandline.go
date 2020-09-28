package commandline

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/go-picturebook"
	"github.com/aaronland/go-picturebook/application"
	"github.com/aaronland/go-picturebook/caption"
	"github.com/aaronland/go-picturebook/filter"
	"github.com/aaronland/go-picturebook/process"
	"github.com/aaronland/go-picturebook/sort"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"gocloud.dev/blob"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var uri_re *regexp.Regexp

var orientation string
var size string
var width float64
var height float64
var dpi float64
var border float64

var source_uri string
var target_uri string

var fill_page bool

var filename string

var verbose bool
var debug bool

var caption_uri string
var sort_uri string

var filter_uris multi.MultiString
var process_uris multi.MultiString

var ocra_font bool

// Deprecated flags

var target string
var preprocess_uris multi.MultiString
var include multi.MultiRegexp
var exclude multi.MultiRegexp

func init() {
	uri_re = regexp.MustCompile(`(?:[a-z0-9_]+):\/\/.*`)
}

type CommandLineApplication struct {
	application.Application
	flagset *flag.FlagSet
}

func DefaultFlagSet(ctx context.Context) (*flag.FlagSet, error) {

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

	fs.StringVar(&orientation, "orientation", "P", "The orientation of your picturebook. Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively.")
	fs.StringVar(&size, "size", "letter", "A common paper size to use for the size of your picturebook. Valid sizes are: [please write me]")
	fs.Float64Var(&width, "width", 8.5, "A custom height to use as the size of your picturebook. Units are currently defined in inches. This fs.overrides the -size fs.")
	fs.Float64Var(&height, "height", 11, "A custom width to use as the size of your picturebook. Units are currently defined in inches. This fs.overrides the -size fs.")
	fs.Float64Var(&dpi, "dpi", 150, "The DPI (dots per inch) resolution for your picturebook.")
	fs.Float64Var(&border, "border", 0.01, "The size of the border around images.")

	fs.BoolVar(&fill_page, "fill-page", false, "If necessary rotate image 90 degrees to use the most available page space.")

	fs.StringVar(&filename, "filename", "picturebook.pdf", "The filename (path) for your picturebook.")

	fs.BoolVar(&verbose, "verbose", false, "Display verbose output as the picturebook is created.")
	fs.BoolVar(&debug, "debug", false, "DEPRECATED: Please use the -verbose flag instead.")

	fs.StringVar(&caption_uri, "caption", "", desc_captions)
	fs.StringVar(&sort_uri, "sort", "", desc_sorters)

	fs.BoolVar(&ocra_font, "ocra-font", false, "Use an OCR-compatible font for captions.")

	fs.Var(&filter_uris, "filter", desc_filters)
	fs.Var(&process_uris, "process", desc_processes)

	fs.StringVar(&source_uri, "source-uri", "", "A valid GoCloud blob URI to specify where files should be read from. By default file:// URIs are supported.")
	fs.StringVar(&target_uri, "target-uri", "", "A valid GoCloud blob URI to specify where your final picturebook PDF file should be written to. By default file:// URIs are supported.")

	// Deprecated flags

	fs.Var(&preprocess_uris, "pre-process", "DEPRECATED: Please use -process {PROCESS_NAME}:// flag instead.")
	fs.Var(&include, "include", "A valid regular expression to use for testing whether a file should be included in your picturebook. DEPRECATED: Please use -filter regexp://include/?pattern={REGULAR_EXPRESSION} flag instead.")
	fs.Var(&exclude, "exclude", "A valid regular expression to use for testing whether a file should be excluded from your picturebook. DEPRECATED: Please use -filter regexp://exclude/?pattern={REGULAR_EXPRESSION} flag instead.")

	fs.StringVar(&target, "target", "", "Valid targets are: cooperhewitt; flickr; orthis. If defined this flag will set the -filter and -caption flags accordingly. DEPRECATED: Please use specific -filter and -caption flags as needed.")

	return fs, nil
}

func NewApplication(ctx context.Context, fs *flag.FlagSet) (application.Application, error) {

	app := &CommandLineApplication{
		flagset: fs,
	}

	return app, nil
}

func (app *CommandLineApplication) Run(ctx context.Context) error {

	flagset.Parse(app.flagset)

	// get flags here...

	if debug {

		log.Println("WARNING The -debug flag is deprecated. Please use the -verbose flag instead.")
		verbose = debug
	}

	if target != "" {

		log.Println("WARNING The -target flag is deprecated. Please use specific -filter and -caption flags as needed.")

		str_filter := fmt.Sprintf("%s://", target)
		str_caption := fmt.Sprintf("%s://", target)

		err := filter_uris.Set(str_filter)

		if err != nil {
			msg := fmt.Sprintf("Failed to assign filter '%s', %v", str_filter, err)
			return errors.New(msg)
		}

		if caption_uri != "" {
			msg := fmt.Sprintf("Can not assign -caption using -target since -caption is already defined.")
			return errors.New(msg)
		}

		caption_uri = str_caption
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

	source_bucket, err := blob.OpenBucket(ctx, source_uri)

	if err != nil {
		return err
	}

	target_bucket, err := blob.OpenBucket(ctx, target_uri)

	if err != nil {
		return err
	}

	opts, err := picturebook.NewPictureBookDefaultOptions(ctx)

	if err != nil {
		msg := fmt.Sprintf("Failed to create default picturebook options, %v", err)
		return errors.New(msg)
	}

	opts.Orientation = orientation
	opts.Size = size
	opts.Width = width
	opts.Height = height
	opts.DPI = dpi
	opts.Border = border
	opts.FillPage = fill_page
	opts.Verbose = verbose
	opts.OCRAFont = ocra_font

	processed := make([]string, 0)

	defer func() {

		for _, p := range processed {

			go func(p string) {

				_, err := os.Stat(p)

				// FIX ME...

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

	if caption_uri != "" {

		if !uri_re.MatchString(caption_uri) {
			caption_uri = fmt.Sprintf("%s://", caption_uri)
		}

		c, err := caption.NewCaption(ctx, caption_uri)

		if err != nil {
			log.Fatal(err)
		}

		opts.Caption = c
	}

	if sort_uri != "" {

		s, err := sort.NewSorter(ctx, sort_uri)

		if err != nil {
			return err
		}

		opts.Sort = s
	}

	sources := app.flagset.Args()

	if len(sources) == 0 {

		base := filepath.Base(source_uri)
		root := filepath.Dir(source_uri)

		sb, err := blob.OpenBucket(ctx, root)

		if err != nil {
			return err
		}

		source_bucket = sb
		sources = []string{base}
	}

	opts.Source = source_bucket
	opts.Target = target_bucket

	pb, err := picturebook.NewPictureBook(ctx, opts)

	if err != nil {
		msg := fmt.Sprintf("Failed to create new picturebook, %v", err)
		return errors.New(msg)
	}

	err = pb.AddPictures(ctx, sources)

	if err != nil {
		return err
	}

	err = pb.Save(ctx, filename)

	if err != nil {
		return err
	}

	return nil
}
