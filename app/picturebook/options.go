package picturebook

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	SourceBucketURI string
	TargetBucketURI string
	TempBucketURI   string

	Orientation string
	Size        string
	Width       float64
	Height      float64
	Units       string
	DPI         float64
	OCRAFont    bool
	Border      float64
	Bleed       float64
	FillPage    bool

	EvenOnly bool
	OddOnly  bool

	MaxPages int

	MarginTop    float64
	MarginBottom float64
	MarginLeft   float64
	MarginRight  float64

	FilterURIs  []string
	ProcessURIs []string
	CaptionURIs []string
	TextURI     string
	SortURI     string

	Sources []string

	Filename string
	Verbose  bool
}

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

		Sources:  fs.Args(),
		Filename: filename,

		Verbose: verbose,
	}

	return opts, nil
}
