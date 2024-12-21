package picturebook

import (
	"context"
	"flag"
	"fmt"
	gosort "sort"
	"strings"

	"github.com/aaronland/go-picturebook/caption"
	"github.com/aaronland/go-picturebook/filter"
	"github.com/aaronland/go-picturebook/process"
	"github.com/aaronland/go-picturebook/sort"
	"github.com/aaronland/go-picturebook/text"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"gocloud.dev/blob"
)

// String label defining the orientation of picturebook PDF files. Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively.
var orientation string

// A common paper size to use for the size of your picturebook. Valid sizes are: "a3", "a4", "a5", "letter", "legal", or "tabloid".
var size string

// A width height to use as the size for a picturebook PDF file.
var width float64

// A custom height to use as the size for a picturebook PDF file.
var height float64

// The unit of measurement to apply to the height and width of a picturebook PDF file.
var units string

// The "dots per inch" (DPI) resolution for a picturebook PDF file.
var dpi float64

// The size of the border to apply to each image in a picturebook PDF file.
var border float64

// The size of the margin to be applied to all sides of a picturebook.
var margin float64

// The size of the top margin for a picturebook.
var margin_top float64

// The size of the bottom margin for a picturebook.
var margin_bottom float64

// The size of the left margin for a picturebook.
var margin_left float64

// The size of the right margin for a picturebook.
var margin_right float64

// The size of an exterior "bleed" margin for a picturebook.
var bleed float64

// A valid aaronland/go-picturebook/bucket.Bucket URI for where source input images are read from.
var source_uri string

// A valid aaronland/go-picturebook/bucket.Bucket URI for where the final picturebook file will be written to.
var target_uri string

// A valid aaronland/go-picturebook/bucket.Bucket URI for where temporary picturebook-related images will be written to and read from.
var tmpfile_uri string

// A boolean flag indicating that, when necessary, an image should be rotated 90 degrees to use the most available page space.
var fill_page bool

// The base filename of the finished picturebook document.
var filename string

// Boolean flag to indicate that images should only be included on even-numbered pages.
var even_only bool

// Boolean flag to indicate that images should only be included on odd-numbered pages.
var odd_only bool

// Boolean flag to signal verbose logging during the creation of a picturebook.
var verbose bool

// Boolean flag to signal that all the steps to create a picturebook should be taken but without creating a final picturebook document.
var debug bool

// Zero or more valid `caption.Caption` URIs.
var caption_uris multi.MultiString

// A valid `text.Text` URI.
var text_uri string

// A valid `sort.Sorter` URI.
var sort_uri string

// One or more valid `filter.Filter` URIs.
var filter_uris multi.MultiString

// One or more valid `process.Process` URIs.
var process_uris multi.MultiString

// A boolean flag indicating that the OCR-69 font should be used for text.
var ocra_font bool

// The maximum number of pages a picturebook can have.
var max_pages int

// A registered `aaronland/go-picturebook/progress.Monitor` URI used to signal picturebook creation progress.
var progress_monitor_uri string

// formatSchemes takes a list of (scheme) strings and ensure that they are lower-cased and have a trailing "://" string.
func formatSchemes(schemes []string) []string {

	for idx, scheme := range schemes {
		scheme = strings.ToLower(scheme)
		scheme = fmt.Sprintf("%s://", scheme)
		schemes[idx] = scheme
	}

	gosort.Strings(schemes)

	return schemes
}

// formatSchemesAsString takes a list of strings and returns a single comma-separated string.
func formatSchemesAsString(schemes []string) string {
	schemes = formatSchemes(schemes)
	return strings.Join(schemes, ", ")
}

// DefaultFlagSet returns a `flag.FlagSet` with required flags and default values for a `CommandLineApplication` instance.
func DefaultFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("picturebook")

	available_buckets := blob.DefaultURLMux().BucketSchemes()
	available_buckets_str := formatSchemesAsString(available_buckets)

	available_filters := filter.AvailableFilters()
	available_filters_str := formatSchemesAsString(available_filters)

	available_captions := caption.AvailableCaptions()
	available_captions_str := formatSchemesAsString(available_captions)

	available_texts := text.AvailableTexts()
	available_texts_str := formatSchemesAsString(available_texts)

	available_processes := process.AvailableProcesses()
	available_processes_str := formatSchemesAsString(available_processes)

	available_sorters := sort.AvailableSorters()
	available_sorters_str := formatSchemesAsString(available_sorters)

	desc_filters := fmt.Sprintf("A valid filter.Filter URI. Valid schemes are: %s.", available_filters_str)
	desc_captions := fmt.Sprintf("Zero or more valid caption.Caption URIs. Valid schemes are: %s.", available_captions_str)
	desc_texts := fmt.Sprintf("A valid text.Text URI. Valid schemes are: %s.", available_texts_str)
	desc_processes := fmt.Sprintf("A valid process.Process URI. Valid schemes are: %s.", available_processes_str)
	desc_sorters := fmt.Sprintf("A valid sort.Sorter URI. Valid schemes are: %s.", available_sorters_str)

	desc_buckets := fmt.Sprintf("A valid GoCloud blob URI to specify where files should be read from. Available schemes are: %s. If no URI scheme is included then the file:// scheme is assumed.", available_buckets_str)

	fs.StringVar(&orientation, "orientation", "P", "The orientation of your picturebook. Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively.")
	fs.StringVar(&size, "size", "letter", `A common paper size to use for the size of your picturebook. Valid sizes are: "a3", "a4", "a5", "letter", "legal", or "tabloid".`)
	fs.Float64Var(&width, "width", 0.0, "A custom height to use as the size of your picturebook. Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -height flag.")
	fs.Float64Var(&height, "height", 0.0, "A custom width to use as the size of your picturebook. Units are defined in inches by default. This flag overrides the -size flag when used in combination with the -width flag.")
	fs.StringVar(&units, "units", "inches", "The unit of measurement to apply to the -height and -width flags. Valid options are inches, millimeters, centimeters")
	fs.Float64Var(&dpi, "dpi", 150, "The DPI (dots per inch) resolution for your picturebook.")
	fs.Float64Var(&border, "border", 0.01, "The size of the border around images.")

	fs.Float64Var(&margin_top, "margin-top", 1.0, "The margin around the top of each page.")
	fs.Float64Var(&margin_bottom, "margin-bottom", 1.0, "The margin around the bottom of each page.")
	fs.Float64Var(&margin_left, "margin-left", 1.0, "The margin around the left-hand side of each page.")
	fs.Float64Var(&margin_right, "margin-right", 1.0, "The margin around the right-hand side of each page.")
	fs.Float64Var(&margin, "margin", 0.0, "The margin around all sides of a page. If non-zero this value will be used to populate all the other -margin-(N) flags.")

	fs.Float64Var(&bleed, "bleed", 0.0, "An additional bleed area to add (on all four sides) to the size of your picturebook.")

	fs.BoolVar(&fill_page, "fill-page", false, "If necessary rotate image 90 degrees to use the most available page space. Note that any '-process' flags involving colour space manipulation will automatically be applied to images after they have been rotated.")

	fs.StringVar(&filename, "filename", "picturebook.pdf", "The filename (path) for your picturebook.")

	fs.BoolVar(&verbose, "verbose", false, "Display verbose output as the picturebook is created.")

	fs.BoolVar(&even_only, "even-only", false, "Only include images on even-numbered pages.")
	fs.BoolVar(&odd_only, "odd-only", false, "Only include images on odd-numbered pages.")

	fs.Var(&caption_uris, "caption", desc_captions)

	fs.StringVar(&text_uri, "text", "", desc_texts)

	fs.StringVar(&sort_uri, "sort", "", desc_sorters)

	fs.BoolVar(&ocra_font, "ocra-font", false, "Use an OCR-compatible font for captions.")

	fs.Var(&filter_uris, "filter", desc_filters)
	fs.Var(&process_uris, "process", desc_processes)

	fs.StringVar(&source_uri, "source-uri", "", desc_buckets)

	desc_buckets_target := fmt.Sprintf("%s If empty then the code will try to use the operating system's 'current working directory' where applicable.", desc_buckets)
	fs.StringVar(&target_uri, "target-uri", "", desc_buckets_target)

	desc_buckets_tmp := fmt.Sprintf("%s If empty the operating system's temporary directory will be used.", desc_buckets)
	fs.StringVar(&tmpfile_uri, "tmpfile-uri", "", desc_buckets_tmp)

	fs.IntVar(&max_pages, "max-pages", 0, "An optional value to indicate that a picturebook should not exceed this number of pages")

	fs.StringVar(&progress_monitor_uri, "progress-monitor-uri", "progressbar://", "A registered aaronland/go-picturebook/progress.Monitor URI")
	return fs, nil
}
