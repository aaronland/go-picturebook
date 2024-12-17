// package picturebook provides a command-line application for creating picturebooks.
package picturebook

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	pb "github.com/aaronland/go-picturebook"
	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-picturebook/caption"
	"github.com/aaronland/go-picturebook/filter"
	"github.com/aaronland/go-picturebook/process"
	"github.com/aaronland/go-picturebook/sort"
	"github.com/aaronland/go-picturebook/text"
	"github.com/sfomuseum/go-flags/flagset"
)

// Regular expression for validating filter and caption URIs.
var uri_re *regexp.Regexp

func init() {
	uri_re = regexp.MustCompile(`(?:[a-z0-9_]+):\/\/.*`)
}

func Run(ctx context.Context) error {

	fs, err := DefaultFlagSet(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create default flag set, %w", err)
	}

	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	logger := slog.Default()

	// START OF unfortunate bit of hoop-jumping to (re) register gocloud stuff
	// because of the way Go imports are ordered.

	err := bucket.RegisterGoCloudBuckets(ctx)

	if err != nil {
		return fmt.Errorf("Failed to register gocloud buckets, %w", err)
	}

	// END OF unfortunate bit of hoop-jumping to (re) register gocloud stuff

	if tmpfile_uri == "" {

		tmpfile_uri = fmt.Sprintf("file://%s", os.TempDir())

		if verbose {
			logger.Debug("Using operating system temporary directory for processing files", "uri", tmpfile_uri)
		}
	}

	if margin != 0.0 {
		margin_top = margin
		margin_bottom = margin
		margin_left = margin
		margin_right = margin
	}

	source_uri, err := ensureScheme(source_uri)

	if err != nil {
		return fmt.Errorf("Failed to ensure scheme for source URI %s, %w", source_uri, err)
	}

	if target_uri == "" {

		cwd, err := os.Getwd()

		if err != nil {
			return fmt.Errorf("Failed to determine current working directory, %w", err)
		}

		target_uri = cwd
	}

	target_uri, err := ensureScheme(target_uri)

	if err != nil {
		return fmt.Errorf("Failed to ensure scheme for target URI %s, %w", target_uri, err)
	}

	target_uri, err = ensureSkipMetadata(target_uri)

	if err != nil {
		return fmt.Errorf("Failed to ensure ?metadata=skip for target URI %s, %w", target_uri, err)
	}

	tmpfile_uri, err := ensureScheme(tmpfile_uri)

	if err != nil {
		return fmt.Errorf("Failed to ensure scheme for tmpfile URI %s, %w", tmpfile_uri, err)
	}

	tmpfile_uri, err = ensureSkipMetadata(tmpfile_uri)

	if err != nil {
		return fmt.Errorf("Failed to ensure ?metadata=skip for tmpfile URI %s, %w", tmpfile_uri, err)
	}

	source_bucket, err := bucket.NewBucket(ctx, source_uri)

	if err != nil {
		return fmt.Errorf("Failed to open source bucket, %w", err)
	}

	target_bucket, err := bucket.NewBucket(ctx, target_uri)

	if err != nil {
		return fmt.Errorf("Failed to open target bucket, %w", err)
	}

	tmpfile_bucket, err := bucket.NewBucket(ctx, tmpfile_uri)

	if err != nil {
		return fmt.Errorf("Failed to open tmpfile bucket, %w", err)
	}

	opts, err := pb.NewPictureBookDefaultOptions(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create default picturebook options, %w", err)
	}

	opts.Orientation = orientation
	opts.Size = size
	opts.Width = width
	opts.Height = height
	opts.Units = units
	opts.DPI = dpi
	opts.Border = border
	opts.Bleed = bleed
	opts.MarginTop = margin_top
	opts.MarginBottom = margin_bottom
	opts.MarginLeft = margin_left
	opts.MarginRight = margin_right
	opts.FillPage = fill_page
	opts.Verbose = verbose
	opts.OCRAFont = ocra_font
	opts.EvenOnly = even_only
	opts.OddOnly = odd_only
	opts.MaxPages = max_pages
	opts.Logger = logger

	processed := make([]string, 0)

	defer func() {

		for _, p := range processed {

			go func(p string) {

				_, err := os.Stat(p)

				if os.IsNotExist(err) {
					return
				}

				logger.Debug("Remove temporary file", "path", p)
				os.Remove(p)
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
				return fmt.Errorf("Failed to create filter '%s', %w", filter_uri, err)
			}

			filters[idx] = f
		}

		multi, err := filter.NewMultiFilter(ctx, filters...)

		if err != nil {
			return fmt.Errorf("Failed to create multi filter, %w", err)
		}

		opts.Filter = multi
	}

	if len(process_uris) > 0 {

		processes := make([]process.Process, len(process_uris))
		rotatetofill_processes := make([]process.Process, 0)

		for idx, process_uri := range process_uris {

			pr, err := process.NewProcess(ctx, process_uri)

			if err != nil {
				return fmt.Errorf("Failed to create process '%s', %w", process_uri, err)
			}

			processes[idx] = pr

			if strings.HasPrefix(process_uri, "colorspace://") || strings.HasPrefix(process_uri, "colourspace://") {
				rotatetofill_processes = append(rotatetofill_processes, pr)
			}
		}

		multi, err := process.NewMultiProcess(ctx, processes...)

		if err != nil {
			return fmt.Errorf("Failed to create multi process, %w", err)
		}

		opts.PreProcess = multi

		if len(rotatetofill_processes) > 0 {

			rotatetofill_multi, err := process.NewMultiProcess(ctx, rotatetofill_processes...)

			if err != nil {
				return fmt.Errorf("Failed to create multi process for rotate to fill post processing, %w", err)
			}

			opts.RotateToFillPostProcess = rotatetofill_multi
		}
	}

	if len(caption_uris) > 0 {

		captions := make([]caption.Caption, len(caption_uris))

		for idx, c_uri := range caption_uris {

			if !uri_re.MatchString(c_uri) {
				c_uri = fmt.Sprintf("%s://", c_uri)
			}

			c, err := caption.NewCaption(ctx, c_uri)

			if err != nil {
				return fmt.Errorf("Failed to create new caption for '%s', %w", c_uri, err)
			}

			captions[idx] = c
		}

		c_opts := &caption.MultiCaptionOptions{
			Captions:   captions,
			Combined:   false,
			AllowEmpty: true,
		}

		c, err := caption.NewMultiCaptionWithOptions(ctx, c_opts)

		if err != nil {
			return fmt.Errorf("Failed to create multi caption, %w", err)
		}

		opts.Caption = c
	}

	if text_uri != "" {

		if !uri_re.MatchString(text_uri) {
			text_uri = fmt.Sprintf("%s://", text_uri)
		}

		t, err := text.NewText(ctx, text_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new text, %w", err)
		}

		opts.Text = t
	}

	if sort_uri != "" {

		s, err := sort.NewSorter(ctx, sort_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new sorter, %w", err)
		}

		opts.Sort = s
	}

	sources := fs.Args()

	if len(sources) == 0 {

		base := filepath.Base(source_uri)
		root := filepath.Dir(source_uri)

		sb, err := bucket.NewBucket(ctx, root)

		if err != nil {
			return fmt.Errorf("Failed to open bucket for %s, %w", root, err)
		}

		source_bucket = sb
		sources = []string{base}
	}

	opts.Source = source_bucket
	opts.Target = target_bucket
	opts.Temporary = tmpfile_bucket

	pb, err := pb.NewPictureBook(ctx, opts)

	if err != nil {
		return fmt.Errorf("Failed to create new picturebook, %v", err)
	}

	err = pb.AddPictures(ctx, sources)

	if err != nil {
		return fmt.Errorf("Failed to add pictures to picturebook, %w", err)
	}

	err = pb.Save(ctx, filename)

	if err != nil {
		return fmt.Errorf("Failed to save picturebook, %w", err)
	}

	return nil
}

// ensureScheme ensures that 'uri' has a valid URI scheme. If the scheme is empty then a default of "file" is applied to 'uri'.
func ensureScheme(uri string) (string, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return "", fmt.Errorf("Failed to parse URI '%s', %w", uri, err)
	}

	if u.Scheme == "" {
		u.Scheme = "file"
	}

	return u.String(), nil
}

// ensureScheme ensures that 'uri' has a '?metadata=skip' query parameter, adding one if necessary.
func ensureSkipMetadata(uri string) (string, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return "", fmt.Errorf("Failed to parse URI '%s', %w", uri, err)
	}

	q := u.Query()

	m := q.Get("metadata")

	if m == "skip" {
		return uri, nil
	}

	q.Del("metadata")
	q.Set("metadata", "skip")

	u.RawQuery = q.Encode()
	return u.String(), nil
}
