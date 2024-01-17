package picturebook

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image/png"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aaronland/go-image/decode"
	"github.com/aaronland/go-image/rotate"
	"github.com/aaronland/go-mimetypes"
	"github.com/aaronland/go-picturebook/caption"
	"github.com/aaronland/go-picturebook/filter"
	"github.com/aaronland/go-picturebook/picture"
	"github.com/aaronland/go-picturebook/process"
	"github.com/aaronland/go-picturebook/sort"
	"github.com/aaronland/go-picturebook/tempfile"
	"github.com/aaronland/go-picturebook/text"
	"github.com/go-pdf/fpdf"
	"github.com/sfomuseum/go-font-ocra"
	"gocloud.dev/blob"
)

// MM2INCH defines the number if millimeters in an inch.
const MM2INCH float64 = 25.4

// PictureBookOptions defines a struct containing configuration information for a given picturebook instance.
type PictureBookOptions struct {
	// The orientation of the final picturebook. Valid options are "P" and "L" for portrait and landscape respectively.
	Orientation string
	// A string label corresponding to known size. Valid options are "a1", "a2", "a3", "a4", "a5", "a6", "a7", "letter", "legal" and "tabloid".
	Size string
	// The width of the final picturebook.
	Width float64
	// The height of the final picturebook.
	Height float64
	// The unit of measurement to use for the `Width` and `Height` options.
	Units string
	// The number dots per inch to use when calculating the size of the final picturebook. Valid options are "inches", "centimeters", "millimeters".
	DPI float64
	// The size of any border to apply to each image in the final picturebook.
	Border float64
	// The size of any additional bleed to apply to the final picturebook.
	Bleed float64
	// The size of any margin to add to the top of each page.
	MarginTop float64
	// The size of any margin to add to the bottom of each page.
	MarginBottom float64
	// The size of any margin to add to the left-hand side of each page.
	MarginLeft float64
	// The size of any margin to add to the right-hand side of each page.
	MarginRight float64
	// An optional `filter.Filter` instance used to determine whether or not an image should be included in the final picturebook.
	Filter filter.Filter
	// Zero or more optional `process.Process` instance used to transform images being included in the final picturebook.
	PreProcess process.Process
	// Zero or more optional `process.Process` instance used to transform images after having been rotated to fill the page and before being included in the final picturebook.
	RotateToFillPostProcess process.Process
	// An optional `caption.Caption` instance used to derive a caption string for each image added to the final picturebook.
	Caption caption.Caption
	// An optional `text.Text` instance used to derive a text string for each image added to the final picturebook.
	Text text.Text
	// An optional `sort.Sorter` instance used to sort images before they are added to the final picturebook.
	Sort sort.Sorter
	// A boolean value signaling that an image should be rotated if necessary to fill the maximum amount of any given page.
	FillPage bool
	// A boolean value to enable verbose logging during the creation of a picturebook.
	Verbose bool
	// A boolean value to enable to use of an OCRA font for writing captions.
	OCRAFont bool
	// A gocloud.dev/blob `Bucket` instance where source images are stored.
	Source *blob.Bucket
	// A gocloud.dev/blob `Bucket` instance where the final picturebook is written to.
	Target *blob.Bucket
	// A gocloud.dev/blob `Bucket` instance where are temporary files necessary in the creation of the picturebook are written to.
	Temporary *blob.Bucket
	// A boolean value signaling that images should only be added on even-numbered pages.
	EvenOnly bool
	// A boolean value signaling that images should only be added on odd-numbered pages.
	OddOnly bool
	// An optional value to indicate that a picturebook should not exceed this number of pages
	MaxPages int
	Logger   *slog.Logger
}

// type PictureBookMargins defines a struct for storing margins to be applied to a picturebook
type PictureBookMargins struct {
	// The size of any margin to add to the top of each page.
	Top float64
	// The size of any margin to add to the bottom of each page.
	Bottom float64
	// The size of any margin to add to the left-hand side of each page.
	Left float64
	// The size of any margin to add to the right-hand side of each page.
	Right float64
}

// type PictureBookBorders defines a struct for storing borders to be applied to a images in a picturebook.
type PictureBookBorders struct {
	// The size of any border to add to the top of each image.
	Top float64
	// The size of any border to add to the bottom of each image.
	Bottom float64
	// The size of any border to add to the left-hand side of each image.
	Left float64
	// The size of any border to add to the right-hand side of each image.
	Right float64
}

// type PictureBookCanvas defines a struct for storing canvas information for a picturebook.
type PictureBookCanvas struct {
	// The width of the internal picturebook canvas.
	Width float64
	// The height of the internal picturebook canvas.
	Height float64
}

// type PictureBookText defines a struct for storing information for how text should be displayed in a picturebook.
type PictureBookText struct {
	// The name of the font to use for text strings.
	Font string
	// The style of the font to use for text strings.
	Style string
	// The size of the font to use for text strings.
	Size float64
	// The margin to apply to text strings.
	Margin float64
	// The colour of the font to use for text strings.
	Colour []int
}

// type PictureBook provides a struct for creating a PDF file from a folder of images (a picturebook).
type PictureBook struct {
	// A `fpdf.Fpdf` instance used to produce the picturebook PDF file.
	PDF *fpdf.Fpdf
	// A `sync.Mutex` instance used to add images in an orderly fashion.
	Mutex *sync.Mutex
	// The `PictureBookBorders` definition to use for this picturebook
	Borders *PictureBookBorders
	// The `PictureBookMargins` definition to use for this picturebook
	Margins *PictureBookMargins
	// The `PictureBookCanvas` definition to use for this picturebook
	Canvas PictureBookCanvas
	// The `PictureBookText` definition to use for this picturebook
	Text PictureBookText
	// The `PictureBookOptions` used to create this picturebook
	Options *PictureBookOptions
	// The `GatherPicturesProcessFunc` function used to determine whether an image is included in a picturebook
	ProcessFunc GatherPicturesProcessFunc
	// The number of pages in this picturebook
	pages int
	// A list of temporary files used in the creation of a picturebook and to be removed when the picturebook is saved
	tmpfiles []string
}

// type GatherPicturesProcessFunc defines a method for processing the path to an image file in to a `picture.PictureBookPicture` instance.
type GatherPicturesProcessFunc func(context.Context, string) (*picture.PictureBookPicture, error)

//	DefaultGatherPicturesProcessFunc returns a default GatherPicturesProcessFunc used to derive a `picture.PictureBookPicture` instance
//
// from the path to an image file. It applies any filters and transformation processes and derives caption data per settings defined in 'pb_opts'.
func DefaultGatherPicturesProcessFunc(pb_opts *PictureBookOptions) (GatherPicturesProcessFunc, error) {

	fn := func(ctx context.Context, path string) (*picture.PictureBookPicture, error) {

		select {
		case <-ctx.Done():
			return nil, nil
		default:
			// pass
		}

		abs_path := path
		is_image := false

		logger := pb_opts.Logger.With("path", abs_path)

		ext := filepath.Ext(abs_path)
		ext = strings.ToLower(ext)

		for _, t := range mimetypes.TypesByExtension(ext) {
			if strings.HasPrefix(t, "image/") {
				is_image = true
				break
			}
		}

		if !is_image {
			logger.Debug("File does not appear to be an image, skipping", "extension", ext)
			return nil, nil
		}

		if pb_opts.Filter != nil {

			ok, err := pb_opts.Filter.Continue(ctx, pb_opts.Source, abs_path)

			if err != nil {
				logger.Error("Failed to filter image", "error", err)
				return nil, nil
			}

			if !ok {
				return nil, nil
			}

			logger.Debug("Include image")
		}

		caption := ""
		text_body := ""

		if pb_opts.Caption != nil {

			txt, err := pb_opts.Caption.Text(ctx, pb_opts.Source, abs_path)

			if err != nil {
				logger.Error("Failed to derive caption text", "error", err)
				return nil, nil
			}

			caption = txt
		}

		if pb_opts.Text != nil {

			txt, err := pb_opts.Text.Body(ctx, pb_opts.Source, abs_path)

			if err != nil {
				logger.Error("Failed to derive text body", "error", err)
				return nil, nil
			}

			text_body = txt
		}

		var final_bucket *blob.Bucket
		final_path := abs_path

		var tmpfile_path string

		if pb_opts.PreProcess != nil {

			logger.Debug("Pre-processing image")

			processed_path, err := pb_opts.PreProcess.Transform(ctx, pb_opts.Source, pb_opts.Temporary, abs_path)

			if err != nil {
				logger.Error("Failed to process image", "error", err)
				return nil, nil
			}

			logger.Debug("After pre-processing path becomes", "processed_path", processed_path)

			if processed_path != "" && processed_path != abs_path {
				final_path = processed_path
				final_bucket = pb_opts.Temporary
				tmpfile_path = processed_path
			}
		}

		logger.Debug("Append path to list for processing", "final_path", final_path)

		pic := &picture.PictureBookPicture{
			Source:   abs_path,
			Bucket:   final_bucket,
			Path:     final_path,
			Caption:  caption,
			Text:     text_body,
			TempFile: tmpfile_path,
		}

		return pic, nil
	}

	return fn, nil
}

// NewPictureBookDefaultOptions returns a `PictureBookOptions` with default settings.
func NewPictureBookDefaultOptions(ctx context.Context) (*PictureBookOptions, error) {

	opts := &PictureBookOptions{
		Orientation:  "P",
		Size:         "letter",
		Width:        0.0,
		Height:       0.0,
		Units:        "inches",
		DPI:          150.0,
		Border:       0.01,
		Bleed:        0.0,
		MarginTop:    1.0,
		MarginBottom: 1.0,
		MarginLeft:   1.0,
		MarginRight:  1.0,
		Verbose:      false,
		Logger:       slog.Default(),
	}

	return opts, nil
}

// NewPictureBook returns a new `PictureBook` instances configured according to the settings in 'opts'.
func NewPictureBook(ctx context.Context, opts *PictureBookOptions) (*PictureBook, error) {

	var pdf *fpdf.Fpdf

	// opts_w := opts.Width
	// opts_h := opts.Height
	// opts_b := opts.Bleed

	// Start by convert everything to inches - not because it's better but
	// just because it's expedient right now (20210218/straup)

	if opts.Width == 0.0 && opts.Height == 0.0 {

		switch strings.ToLower(opts.Size) {
		case "a1":
			opts.Width = 584.0 / MM2INCH
			opts.Height = 841.0 / MM2INCH
		case "a2":
			opts.Width = 420 / MM2INCH
			opts.Height = 594 / MM2INCH
		case "a3":
			opts.Width = 297 / MM2INCH
			opts.Height = 420 / MM2INCH
		case "a4":
			opts.Width = 210.0 / MM2INCH
			opts.Height = 297.0 / MM2INCH
		case "a5":
			opts.Width = 148 / MM2INCH
			opts.Height = 210 / MM2INCH
		case "a6":
			opts.Width = 105 / MM2INCH
			opts.Height = 148 / MM2INCH
		case "a7":
			opts.Width = 74 / MM2INCH
			opts.Height = 105 / MM2INCH
		case "letter":
			opts.Width = 8.5
			opts.Height = 11.0
		case "legal":
			opts.Width = 11.0
			opts.Height = 17.0
		case "tabloid":
			opts.Width = 11.0
			opts.Height = 17.0
		default:
			return nil, fmt.Errorf("Unrecognized page size '%s'", opts.Size)
		}
	} else {

		switch opts.Units {
		case "inches":
			// pass
		case "millimeters":
			opts.Width = opts.Width / MM2INCH
			opts.Height = opts.Height / MM2INCH
		case "centimeters":
			opts.Width = (opts.Width * 10.0) / MM2INCH
			opts.Height = (opts.Height * 10.0) / MM2INCH
		default:
			return nil, fmt.Errorf("Invalid or unsupported unit '%s'", opts.Units)
		}
	}

	// log.Printf("%0.2f x %0.2f (%s)\n", opts.Width, opts.Height, opts.Size)

	sz := fpdf.SizeType{
		Wd: opts.Width + (opts.Bleed * 2.0),
		Ht: opts.Height + (opts.Bleed * 2.0),
	}

	init := fpdf.InitType{
		OrientationStr: opts.Orientation,
		UnitStr:        "in",
		SizeStr:        "",
		Size:           sz,
		FontDirStr:     "",
	}

	pdf = fpdf.NewCustom(&init)

	t := PictureBookText{
		Font:   "Helvetica",
		Style:  "",
		Size:   8.0,
		Margin: 0.1,
		Colour: []int{128, 128, 128},
	}

	if opts.OCRAFont {

		font, err := ocra.LoadFPDFFont()

		if err != nil {
			return nil, fmt.Errorf("Failed to load OCRA font, %w", err)
		}

		pdf.AddFontFromBytes(font.Family, font.Style, font.JSON, font.Z)
		pdf.SetFont(font.Family, "", 8.0)

		pdf.SetTextColor(t.Colour[0], t.Colour[1], t.Colour[2])

	} else {

		pdf.SetFont(t.Font, t.Style, t.Size)
	}

	w, h, _ := pdf.PageSize(1)

	page_w := w * opts.DPI
	page_h := h * opts.DPI

	// https://github.com/aaronland/go-picturebook/issues/22

	// margin around each page (inclusive of page bleed)

	margin_top := (opts.MarginTop + (opts.Bleed * 2.0)) * opts.DPI
	margin_bottom := (opts.MarginBottom + (opts.Bleed * 2.0)) * opts.DPI
	margin_left := (opts.MarginLeft + (opts.Bleed * 2.0)) * opts.DPI
	margin_right := (opts.MarginRight + (opts.Bleed * 2.0)) * opts.DPI

	margins := &PictureBookMargins{
		Top:    margin_top,
		Bottom: margin_bottom,
		Left:   margin_left,
		Right:  margin_right,
	}

	// border around each image

	border_top := opts.Border * opts.DPI
	border_bottom := opts.Border * opts.DPI
	border_left := opts.Border * opts.DPI
	border_right := opts.Border * opts.DPI

	borders := &PictureBookBorders{
		Top:    border_top,
		Bottom: border_bottom,
		Left:   border_left,
		Right:  border_right,
	}

	// Remember: margins have been calculated inclusive of page bleeds

	canvas_w := page_w - (margin_left + margin_right + border_left + border_right)
	canvas_h := page_h - (margin_top + margin_bottom + border_top + border_bottom)

	pdf.SetAutoPageBreak(false, border_bottom)

	canvas := PictureBookCanvas{
		Width:  canvas_w,
		Height: canvas_h,
	}

	tmpfiles := make([]string, 0)
	mu := new(sync.Mutex)

	process_func, err := DefaultGatherPicturesProcessFunc(opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to return DefaultGatherPicturesProcessFunc, %w", err)
	}

	pb := PictureBook{
		PDF:         pdf,
		Mutex:       mu,
		Borders:     borders,
		Margins:     margins,
		Canvas:      canvas,
		Text:        t,
		Options:     opts,
		ProcessFunc: process_func,
		pages:       0,
		tmpfiles:    tmpfiles,
	}

	return &pb, nil
}

// AddPictures adds images founds in one or more folders defined 'paths' to the picturebook instance.
func (pb *PictureBook) AddPictures(ctx context.Context, paths []string) error {

	pictures, err := pb.GatherPictures(ctx, paths)

	if err != nil {
		return fmt.Errorf("Failed to gather pictures, %w", err)
	}

	pb.Options.Logger.Debug("Pictures gathered", "count", len(pictures))

	if pb.Options.Sort != nil {

		sorted, err := pb.Options.Sort.Sort(ctx, pb.Options.Source, pictures)

		if err != nil {
			return fmt.Errorf("Failed to sort pictures, %w", err)
		}

		pictures = sorted
	}

	for _, pic := range pictures {

		pb.Mutex.Lock()
		pb.pages += 1
		pagenum := pb.pages
		pb.Mutex.Unlock()

		var err error

		if pb.Options.EvenOnly {

			if pagenum%2 != 0 {
				pb.AddBlankPage(ctx, pagenum)
				pb.pages += 1
				pagenum = pb.pages
			}

			if pic.Text != "" {
				pb.AddText(ctx, pagenum, pic)
				pb.pages += 1
				pagenum = pb.pages

				pb.AddBlankPage(ctx, pagenum)
				pb.pages += 1
				pagenum = pb.pages
			}

			err = pb.AddPicture(ctx, pagenum, pic)

		} else if pb.Options.OddOnly {

			if pagenum == 1 {
				pb.AddBlankPage(ctx, pagenum)
				pb.pages += 1
				pagenum = pb.pages
			}

			if pagenum%2 == 0 {
				err = pb.AddBlankPage(ctx, pagenum)
				pb.pages += 1
				pagenum = pb.pages
			}

			if pic.Text != "" {
				pb.AddText(ctx, pagenum, pic)
				pb.pages += 1
				pagenum = pb.pages

				pb.AddBlankPage(ctx, pagenum)
				pb.pages += 1
				pagenum = pb.pages
			}

			err = pb.AddPicture(ctx, pagenum, pic)

		} else {

			if pic.Text != "" {
				pb.AddText(ctx, pagenum, pic)
				pb.pages += 1
				pagenum = pb.pages
			}

			err = pb.AddPicture(ctx, pagenum, pic)
		}

		if err != nil {
			pb.Options.Logger.Debug("Failed to add picture", "path", pic.Path, "error", err)
		}
	}

	return nil
}

// GatherPictures collects all the images in one or more folders defined by 'paths' and returns a list of `picture.PictureBookPicture` instances.
func (pb *PictureBook) GatherPictures(ctx context.Context, paths []string) ([]*picture.PictureBookPicture, error) {

	pictures := make([]*picture.PictureBookPicture, 0)

	var list func(context.Context, *blob.Bucket, string) error

	list = func(ctx context.Context, bucket *blob.Bucket, prefix string) error {

		iter := bucket.List(&blob.ListOptions{
			Delimiter: "/",
			Prefix:    prefix,
		})

		for {
			obj, err := iter.Next(ctx)

			if err == io.EOF {
				break
			}

			if err != nil {
				return fmt.Errorf("Failed to iterate next in bucket for %s, %w", prefix, err)
			}

			path := obj.Key

			if obj.IsDir {

				err := list(ctx, bucket, path)

				if err != nil {
					return fmt.Errorf("Failed to list bucket for %s, %w", path, err)
				}

				continue
			}

			pic, err := pb.ProcessFunc(ctx, path)

			if err != nil {
				return fmt.Errorf("Failed to apply ProcessFunc for %s, %w", path, err)
			}

			if pic == nil {
				continue
			}

			if pic.TempFile != "" {
				pb.tmpfiles = append(pb.tmpfiles, pic.TempFile)
			}

			pb.Mutex.Lock()

			pictures = append(pictures, pic)
			count_pictures := len(pictures)

			pb.Mutex.Unlock()

			if pb.Options.MaxPages > 0 && count_pictures >= pb.Options.MaxPages {
				break
			}
		}

		return nil
	}

	for _, path := range paths {

		err := list(ctx, pb.Options.Source, path)

		if err != nil {
			return nil, fmt.Errorf("Failed to list bucket for %s, %w", path, err)
		}
	}

	return pictures, nil
}

// AddBlankPage add a blank page the final PDF document at page 'pagenum'.
func (pb *PictureBook) AddBlankPage(ctx context.Context, pagenum int) error {
	pb.PDF.AddPage()
	return nil
}

// AddText add the value of `pic.Text` on the adjacent page to `pic`.
func (pb *PictureBook) AddText(ctx context.Context, pagenum int, pic *picture.PictureBookPicture) error {

	pb.Mutex.Lock()
	defer pb.Mutex.Unlock()

	pb.PDF.AddPage()

	_, line_h := pb.PDF.GetFontSize()

	max_w := pb.Canvas.Width
	// max_h := pb.Canvas.Height - (pb.Text.Margin + line_h)

	/*
		w := max_w
		h := max_h
	*/

	margins := pb.Margins

	current_x := margins.Left
	current_y := margins.Top

	// START OF reconcile me with code for rendering captions...

	prepped := text.PrepareText(pb.PDF, pb.Options.DPI, max_w, pic.Text)

	for _, txt := range prepped {

		txt = strings.TrimSpace(txt)

		// txt_w := pb.PDF.GetStringWidth(txt)
		txt_h := line_h

		/*
			txt_w = txt_w + pb.Text.Margin
		*/
		txt_h = txt_h + pb.Text.Margin

		// log.Printf("DEBUG %d max: %f03 w: %f03 %s\n", len(txt), max_w, txt_w*pb.Options.DPI, txt)
		// please do this in the constructor...
		// (20171128/thisisaaronland)

		font_sz, _ := pb.PDF.GetFontSize()
		pb.PDF.SetFontSize(font_sz + 2)

		_, line_h := pb.PDF.GetFontSize()

		pb.PDF.SetFontSize(font_sz)

		txt_x := current_x / pb.Options.DPI
		txt_y := (current_y / pb.Options.DPI)

		// pb.Options.Logger.Debug("[%d][%s] text at %0.2f x %0.2f (- x %0.2f)\n", pagenum, pic.Path, txt_x, txt_y, txt_h)

		pb.PDF.SetXY(txt_x, txt_y)

		html := pb.PDF.HTMLBasicNew()
		html.Write(line_h, txt)

		current_y += ((txt_h * pb.Options.DPI) * .65)
	}

	// END OF reconcile me with code for rendering captions...

	return nil
}

func (pb *PictureBook) AddPicture(ctx context.Context, pagenum int, pic *picture.PictureBookPicture) error {

	pb.Mutex.Lock()
	defer pb.Mutex.Unlock()

	abs_path := pic.Path
	caption := pic.Caption

	is_tempfile := false

	picture_bucket := pb.Options.Source

	if pic.Bucket != nil {
		picture_bucket = pic.Bucket
	}

	logger := pb.Options.Logger.With("path", abs_path, "page_number", pagenum)

	im_r, err := picture_bucket.NewReader(ctx, abs_path, nil)

	if err != nil {
		return fmt.Errorf("Failed to create new bucket for %s, %w", abs_path, err)
	}

	defer im_r.Close()

	dec, err := decode.NewDecoder(ctx, abs_path)

	if err != nil {
		return fmt.Errorf("Failed to create new decoder for %s, %w", abs_path, err)
	}

	im, format, err := dec.Decode(ctx, im_r)

	if err != nil {
		return fmt.Errorf("Failed to decode image for %s, %w", abs_path, err)
	}

	// START OF put me somewhere in aaronland/go-image ... maybe?
	// trap fpdf "16-bit depth not supported in PNG file" errors

	if format == "png" {

		buf := new(bytes.Buffer)

		err = png.Encode(buf, im)

		if err != nil {
			return fmt.Errorf("Failed to encode PNG image for %s, %w", abs_path, err)
		}

		// this bit is cribbed from https://github.com/jung-kurt/gofpdf/blob/7d57599b9d9c5fb48ea733596cbb812d7f84a8d6/png.go
		// (20181231/thisisaaronland)

		_ = buf.Next(12)

		var bpc int32
		err := binary.Read(buf, binary.BigEndian, &bpc)

		if err != nil {
			return err
		}

		if bpc > 8 {

			tmpfile_path, tmpfile_format, err := tempfile.TempFileWithImage(ctx, pb.Options.Temporary, im)

			if err != nil {
				return fmt.Errorf("Failed to generate tempfile for %s, %w", abs_path, err)
			}

			logger.Debug("PNG converted to a JPG", "tmpfile_path", tmpfile_path)

			pb.tmpfiles = append(pb.tmpfiles, tmpfile_path)

			abs_path = tmpfile_path
			format = tmpfile_format

			is_tempfile = true
		}
	}

	// END OF put me somewhere in aaronland/go-image ... maybe?

	dims := im.Bounds()

	w := float64(dims.Max.X)
	h := float64(dims.Max.Y)

	logger.Debug("Dimensions", slog.Float64("width", w), slog.Float64("height", h))

	if pb.Options.FillPage {

		image_orientation := "U" // unknown

		if dims.Max.Y > dims.Max.X {
			image_orientation = "P"
		} else if dims.Max.X > dims.Max.Y {
			image_orientation = "L"
		} else {
			// pass
		}

		_, line_h := pb.PDF.GetFontSize()

		max_w := pb.Canvas.Width
		max_h := pb.Canvas.Height - (pb.Text.Margin + line_h)

		rotate_to_fill := false

		if pb.Options.Orientation == "P" && image_orientation == "L" && w > max_w {
			rotate_to_fill = true
		}

		if pb.Options.Orientation == "L" && image_orientation == "P" && h > max_h {
			rotate_to_fill = true
		}

		if rotate_to_fill {

			pb.Options.Logger.Debug("Rotate image to fill path", "path", abs_path)

			new_im, err := rotate.RotateImageWithDegrees(ctx, im, 90.0)

			if err != nil {
				return err
			}

			im = new_im
			dims = im.Bounds()

			w = float64(dims.Max.X)
			h = float64(dims.Max.Y)

			// now save to disk...

			tmpfile_path, tmpfile_format, err := tempfile.TempFileWithImage(ctx, pb.Options.Temporary, im)

			if err != nil {
				return fmt.Errorf("Failed to create temporary file (rotate to fill) for %s, %w", abs_path, err)
			}

			if pb.Options.RotateToFillPostProcess != nil {

				tmpfile_path, err = pb.Options.RotateToFillPostProcess.Transform(ctx, pb.Options.Temporary, pb.Options.Temporary, tmpfile_path)

				if err != nil {
					return fmt.Errorf("Failed to apply colour space transformations to temporary file (rotate to fill), %w", err)
				}
			}

			pb.tmpfiles = append(pb.tmpfiles, tmpfile_path)

			logger.Debug("Append rotated image", "tmpfile_path", tmpfile_path)

			abs_path = tmpfile_path
			format = tmpfile_format

			is_tempfile = true
		}
	}

	opts := fpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: format,
	}

	var r io.ReadCloser

	if is_tempfile {
		r, err = pb.Options.Temporary.NewReader(ctx, abs_path, nil)
	} else {
		r, err = picture_bucket.NewReader(ctx, abs_path, nil)
	}

	if err != nil {
		return fmt.Errorf("Failed to create new reader (info) for %s, %v", abs_path, err)
	}

	defer r.Close()

	info := pb.PDF.RegisterImageOptionsReader(abs_path, opts, r)

	if info == nil {
		return fmt.Errorf("unable to determine info for %s", abs_path)
	}

	info.SetDpi(pb.Options.DPI)

	logger.Debug("Dimensions", slog.Float64("width", w), slog.Float64("height", h))

	if w == 0.0 || h == 0.0 {
		return fmt.Errorf("[%d] %s has zero-sized dimension", pagenum, abs_path)
	}

	// Remember: margins have been calculated inclusive of page bleeds

	margins := pb.Margins

	x := margins.Left
	y := margins.Top

	_, line_h := pb.PDF.GetFontSize()

	logger.Debug("margins", slog.Float64("left_and_right", (margins.Left+margins.Right)))
	logger.Debug("margins", slog.Float64("top_and_bottom", (margins.Top+margins.Bottom)))
	logger.Debug("margins", slog.Float64("caption", (pb.Text.Margin+line_h)))

	max_w := pb.Canvas.Width
	max_h := pb.Canvas.Height

	logger.Debug("max dimensions", slog.Float64("max_width", max_w), slog.Float64("width", w), slog.Float64("max_height", max_h), slog.Float64("height", h))

	for {

		if w >= max_w || h >= max_h {

			if w > max_w {

				ratio := max_w / w
				w = max_w
				h = h * ratio

			}

			if h > max_h {

				ratio := max_h / h
				w = w * ratio
				h = max_h

			}

		}

		// TO DO: ENSURE ! h < max_h && ! w < max_w

		if w <= max_w && h <= max_h {
			break
		}
	}

	if w < max_w {
		padding := max_w - w
		x = x + (padding / 2.0)
	}

	if h < max_h {
		padding := max_h - h
		y = y + (padding / 2.0)
	}

	logger.Debug("final dimensions", slog.Float64("width", w), slog.Float64("height", h), slog.Float64("x", x), slog.Float64("y", y))

	pb.PDF.AddPage()

	// logger.Debug("final dimensions %0.2f x %0.2f (%0.2f x %0.2f)", w, h, x, y)

	// draw margins

	mx := x / pb.Options.DPI
	my := y / pb.Options.DPI
	mw := w / pb.Options.DPI
	mh := h / pb.Options.DPI

	logger.Debug("margin", slog.Float64("x", mx), slog.Float64("y", my), slog.Float64("width", mw), slog.Float64("height", mh))

	pb.PDF.SetFillColor(0, 0, 0)
	pb.PDF.Rect(mx, my, mw, mh, "FD")

	// draw borders

	borders := pb.Borders
	r_border := borders.Right

	if r_border > 0.0 {

		bx := (x - borders.Left) / pb.Options.DPI
		by := (y - borders.Top) / pb.Options.DPI
		bw := (w + borders.Left + borders.Right) / pb.Options.DPI
		bh := (h + borders.Top + borders.Bottom) / pb.Options.DPI

		logger.Debug("border", slog.Float64("x", bx), slog.Float64("y", by), slog.Float64("width", bw), slog.Float64("height", bh))

		pb.PDF.SetFillColor(0, 0, 0)
		pb.PDF.Rect(bx, by, bw, bh, "FD")
	}

	// draw the image

	// https://godoc.org/github.com/jung-kurt/fpdf#ImageOptions

	image_opts := fpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: format,
	}

	image_x := x / pb.Options.DPI
	image_y := y / pb.Options.DPI
	image_w := w / pb.Options.DPI
	image_h := h / pb.Options.DPI

	logger.Debug("image", slog.Float64("x", image_x), slog.Float64("y", image_y), slog.Float64("width", image_w), slog.Float64("height", image_h))

	pb.PDF.ImageOptions(abs_path, image_x, image_y, image_w, image_h, false, image_opts, 0, "")

	if caption != "" {

		current_x := x
		current_y := y

		for _, txt := range strings.Split(caption, "\n") {

			txt = strings.TrimSpace(txt)

			txt_w := pb.PDF.GetStringWidth(txt)
			txt_h := line_h

			txt_w = txt_w + pb.Text.Margin
			txt_h = txt_h + pb.Text.Margin

			// please do this in the constructor...
			// (20171128/thisisaaronland)

			font_sz, _ := pb.PDF.GetFontSize()
			pb.PDF.SetFontSize(font_sz + 2)

			_, line_h := pb.PDF.GetFontSize()

			logger.Debug("line height %0.2f", line_h)

			pb.PDF.SetFontSize(font_sz)

			txt_x := ((current_x + w) / pb.Options.DPI) - txt_w
			txt_y := ((current_y + h) / pb.Options.DPI) + line_h

			logger.Debug("Text", slog.Float64("x", txt_x), slog.Float64("y", txt_y), slog.Float64("width", txt_w), slog.Float64("height", txt_h))

			// pb.PDF.SetFillColor(255, 255, 255)
			// pb.PDF.Rect(txt_x, txt_y, txt_w, txt_h, "FD")

			pb.PDF.SetXY(txt_x, txt_y)

			logger.Debug("caption", "text", txt)

			html := pb.PDF.HTMLBasicNew()
			html.Write(line_h, txt)

			current_y += ((txt_h * pb.Options.DPI) * .65)
		}
	}

	return nil
}

// Save will write the picturebook to 'path' in the `Target` bucket specified in the `PictureBookOptions`
// used to create the picturebook option.
func (pb *PictureBook) Save(ctx context.Context, path string) error {

	if pb.Options.Target == nil {
		return fmt.Errorf("Missing or invalid target bucket")
	}

	// move this out of here...

	defer func() {

		for _, path := range pb.tmpfiles {

			fname := filepath.Base(path)

			// This shouldn't be necessary and points to a larger problem
			// but this bandaid-fix will have to do for now...
			// (20210103/straup)

			if !strings.HasPrefix(fname, "picturebook-") {
				continue
			}

			pb.Options.Logger.Debug("Remove tmp file", "path", path)

			err := pb.Options.Temporary.Delete(ctx, path)

			if err != nil {
				pb.Options.Logger.Error("Failed to delete tmp file", "path", path, "error", err)
			}
		}
	}()

	pb.Options.Logger.Debug("Save picturebook", "path", path)

	wr, err := pb.Options.Target.NewWriter(ctx, path, nil)

	if err != nil {
		return fmt.Errorf("Failed to create a new writer for %s, %w", path, err)
	}

	err = pb.PDF.Output(wr)

	if err != nil {
		return fmt.Errorf("Failed to output PDF file for %s, %w", path, err)
	}

	err = wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close writer for %s, %w", path, err)
	}

	return nil
}
