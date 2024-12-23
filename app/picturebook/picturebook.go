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
	"github.com/aaronland/go-picturebook/progress"
	"github.com/aaronland/go-picturebook/sort"
	"github.com/aaronland/go-picturebook/text"
)

// Regular expression for validating filter and caption URIs.
var uri_re *regexp.Regexp

func init() {
	uri_re = regexp.MustCompile(`(?:[a-z0-9_]+):\/\/.*`)
}

// Run will run the `picturebook` application configured using the default flagset and options.
func Run(ctx context.Context) error {

	fs, err := DefaultFlagSet(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create default flag set, %w", err)
	}

	return RunWithFlagSet(ctx, fs)
}

// Run will run the `picturebook` application configured using 'fs'.
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

// Run will run the `picturebook` application configured using 'app_opts'.
func RunWithOptions(ctx context.Context, app_opts *RunOptions) error {

	if app_opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	// START OF unfortunate bit of hoop-jumping to (re) register gocloud stuff
	// because of the way Go imports are ordered.

	err := bucket.RegisterGoCloudBuckets(ctx)

	if err != nil {
		return fmt.Errorf("Failed to register gocloud buckets, %w", err)
	}

	// END OF unfortunate bit of hoop-jumping to (re) register gocloud stuff

	source_uri := app_opts.SourceBucketURI
	target_uri := app_opts.TargetBucketURI
	tmpfile_uri := app_opts.TempBucketURI

	source_uri, err = ensureScheme(source_uri)

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

	target_uri, err = ensureScheme(target_uri)

	if err != nil {
		return fmt.Errorf("Failed to ensure scheme for target URI %s, %w", target_uri, err)
	}

	target_uri, err = ensureSkipMetadata(target_uri)

	if err != nil {
		return fmt.Errorf("Failed to ensure ?metadata=skip for target URI %s, %w", target_uri, err)
	}

	tmpfile_uri, err = ensureScheme(tmpfile_uri)

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

	pb_opts, err := pb.NewPictureBookDefaultOptions(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create default picturebook options, %w", err)
	}

	pb_opts.Orientation = app_opts.Orientation
	pb_opts.Size = app_opts.Size
	pb_opts.Width = app_opts.Width
	pb_opts.Height = app_opts.Height
	pb_opts.Units = app_opts.Units
	pb_opts.DPI = app_opts.DPI
	pb_opts.Border = app_opts.Border
	pb_opts.Bleed = app_opts.Bleed
	pb_opts.MarginTop = app_opts.MarginTop
	pb_opts.MarginBottom = app_opts.MarginBottom
	pb_opts.MarginLeft = app_opts.MarginLeft
	pb_opts.MarginRight = app_opts.MarginRight
	pb_opts.FillPage = app_opts.FillPage
	pb_opts.Verbose = app_opts.Verbose
	pb_opts.OCRAFont = app_opts.OCRAFont
	pb_opts.EvenOnly = app_opts.EvenOnly
	pb_opts.OddOnly = app_opts.OddOnly
	pb_opts.MaxPages = app_opts.MaxPages

	processed := make([]string, 0)

	defer func() {

		for _, p := range processed {

			go func(p string) {

				_, err := os.Stat(p)

				if os.IsNotExist(err) {
					return
				}

				slog.Debug("Remove temporary file", "path", p)
				os.Remove(p)
			}(p)
		}
	}()

	if len(app_opts.FilterURIs) > 0 {

		filters := make([]filter.Filter, len(app_opts.FilterURIs))

		for idx, filter_uri := range app_opts.FilterURIs {

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

		pb_opts.Filter = multi
	}

	if len(app_opts.ProcessURIs) > 0 {

		processes := make([]process.Process, len(app_opts.ProcessURIs))
		rotatetofill_processes := make([]process.Process, 0)

		for idx, process_uri := range app_opts.ProcessURIs {

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

		pb_opts.PreProcess = multi

		if len(rotatetofill_processes) > 0 {

			rotatetofill_multi, err := process.NewMultiProcess(ctx, rotatetofill_processes...)

			if err != nil {
				return fmt.Errorf("Failed to create multi process for rotate to fill post processing, %w", err)
			}

			pb_opts.RotateToFillPostProcess = rotatetofill_multi
		}
	}

	if len(app_opts.CaptionURIs) > 0 {

		captions := make([]caption.Caption, len(app_opts.CaptionURIs))

		for idx, c_uri := range app_opts.CaptionURIs {

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

		pb_opts.Caption = c
	}

	if app_opts.TextURI != "" {

		if !uri_re.MatchString(app_opts.TextURI) {
			app_opts.TextURI = fmt.Sprintf("%s://", app_opts.TextURI)
		}

		t, err := text.NewText(ctx, app_opts.TextURI)

		if err != nil {
			return fmt.Errorf("Failed to create new text, %w", err)
		}

		pb_opts.Text = t
	}

	if app_opts.SortURI != "" {

		s, err := sort.NewSorter(ctx, app_opts.SortURI)

		if err != nil {
			return fmt.Errorf("Failed to create new sorter, %w", err)
		}

		pb_opts.Sort = s
	}

	if len(app_opts.Sources) == 0 {

		base := filepath.Base(source_uri)
		root := filepath.Dir(source_uri)

		sb, err := bucket.NewBucket(ctx, root)

		if err != nil {
			return fmt.Errorf("Failed to open bucket for %s, %w", root, err)
		}

		source_bucket = sb
		app_opts.Sources = []string{base}
	}

	monitor, err := progress.NewMonitor(ctx, app_opts.ProgressMonitorURI)

	if err != nil {
		return fmt.Errorf("Failed to create new progress monitor, %w", err)
	}

	pb_opts.Source = source_bucket
	pb_opts.Target = target_bucket
	pb_opts.Temporary = tmpfile_bucket
	pb_opts.Monitor = monitor

	pb, err := pb.NewPictureBook(ctx, pb_opts)

	if err != nil {
		return fmt.Errorf("Failed to create new picturebook, %v", err)
	}

	err = pb.AddPictures(ctx, app_opts.Sources)

	if err != nil {
		return fmt.Errorf("Failed to add pictures to picturebook, %w", err)
	}

	err = pb.Save(ctx, app_opts.Filename)

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
