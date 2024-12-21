package picturebook

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

// RunOptions is a struct containing details about a picturebook to create.
type RunOptions struct {
	// A valid aaronland/go-picturebook/bucket.Bucket URI for where source input images are read from.
	SourceBucketURI string
	// A valid aaronland/go-picturebook/bucket.Bucket URI for where the final picturebook file will be written to.
	TargetBucketURI string
	// A valid aaronland/go-picturebook/bucket.Bucket URI for where temporary picturebook-related images will be written to and read from.
	TempBucketURI string
	// String label defining the orientation of picturebook PDF files. Valid orientations are: 'P' and 'L' for portrait and landscape mode respectively.
	Orientation string
	// A common paper size to use for the size of your picturebook. Valid sizes are: "a3", "a4", "a5", "letter", "legal", or "tabloid".
	Size string
	// A width height to use as the size for a picturebook PDF file.
	Width float64
	// A custom height to use as the size for a picturebook PDF file.
	Height float64
	// The unit of measurement to apply to the height and width of a picturebook PDF file.
	Units string
	// The "dots per inch" (DPI) resolution for a picturebook PDF file.
	DPI float64
	// A boolean flag indicating that the OCR-69 font should be used for text.
	OCRAFont bool
	// The size of the border to apply to each image in a picturebook PDF file.
	Border float64
	// The size of an exterior "bleed" margin for a picturebook.
	Bleed float64
	// A boolean flag indicating that, when necessary, an image should be rotated 90 degrees to use the most available page space.
	FillPage bool
	// Boolean flag to indicate that images should only be included on even-numbered pages.
	EvenOnly bool
	// Boolean flag to indicate that images should only be included on odd-numbered pages.
	OddOnly bool
	// The maximum number of pages a picturebook can have.
	MaxPages int
	// The size of the top margin for a picturebook.
	MarginTop float64
	// The size of the bottom margin for a picturebook.
	MarginBottom float64
	// The size of the left margin for a picturebook.
	MarginLeft float64
	// The size of the right margin for a picturebook.
	MarginRight float64
	// Zero or more valid `filter.Filter` URIs.
	FilterURIs []string
	// Zero or more valid `process.Process` URIs.
	ProcessURIs []string
	// Zero or more valid `caption.Caption` URIs.
	CaptionURIs []string
	// A valid `text.Text` URI.
	TextURI string
	// A valid `sort.Sorter` URI.
	SortURI string
	// One or more paths to crawl for images to add to a picturebook.
	Sources []string
	// The base filename of the finished picturebook document.
	Filename string
	// A registered `aaronland/go-picturebook/progress.Monitor` URI used to signal picturebook creation progress.
	ProgressMonitorURI string
	// Boolean flag to signal verbose logging during the creation of a picturebook.
	Verbose bool
}

// Derive a new `RunOptions` instances from 'fs'.
func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	if tmpfile_uri == "" {

		tmpfile_uri = fmt.Sprintf("file://%s", os.TempDir())
	}

	if margin != 0.0 {
		margin_top = margin
		margin_bottom = margin
		margin_left = margin
		margin_right = margin
	}

	opts := &RunOptions{
		SourceBucketURI: source_uri,
		TargetBucketURI: target_uri,
		TempBucketURI:   tmpfile_uri,

		Orientation: orientation,
		Size:        size,
		Width:       width,
		Height:      height,
		Units:       units,
		DPI:         dpi,

		MarginTop:    margin_top,
		MarginBottom: margin_bottom,
		MarginRight:  margin_right,
		MarginLeft:   margin_left,

		Border:   border,
		Bleed:    bleed,
		FillPage: fill_page,

		EvenOnly:    even_only,
		OddOnly:     odd_only,
		OCRAFont:    ocra_font,
		FilterURIs:  filter_uris,
		ProcessURIs: process_uris,
		CaptionURIs: caption_uris,
		TextURI:     text_uri,
		SortURI:     sort_uri,

		Sources:            fs.Args(),
		Filename:           filename,
		ProgressMonitorURI: progress_monitor_uri,
		Verbose:            verbose,
	}

	return opts, nil
}
