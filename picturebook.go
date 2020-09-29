package picturebook

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/aaronland/go-image-rotate"
	"github.com/aaronland/go-image-tools/util"
	"github.com/aaronland/go-picturebook/caption"
	"github.com/aaronland/go-picturebook/filter"
	"github.com/aaronland/go-picturebook/picture"
	"github.com/aaronland/go-picturebook/process"
	"github.com/aaronland/go-picturebook/sort"
	"github.com/aaronland/go-picturebook/tempfile"
	"github.com/jung-kurt/gofpdf"
	"github.com/rainycape/unidecode"
	"github.com/sfomuseum/go-font-ocra"
	"gocloud.dev/blob"
	"io"
	"log"
	"sync"
)

type PictureBookOptions struct {
	Orientation string
	Size        string
	Width       float64
	Height      float64
	DPI         float64
	Border      float64
	Filter      filter.Filter
	PreProcess  process.Process
	Caption     caption.Caption
	Sort        sort.Sorter
	FillPage    bool
	Verbose     bool
	OCRAFont    bool
	Source      *blob.Bucket
	Target      *blob.Bucket
}

type PictureBookBorder struct {
	Top    float64
	Bottom float64
	Left   float64
	Right  float64
}

type PictureBookCanvas struct {
	Width  float64
	Height float64
}

type PictureBookText struct {
	Font   string
	Style  string
	Size   float64
	Margin float64
	Colour []int
}

type PictureBook struct {
	PDF      *gofpdf.Fpdf
	Mutex    *sync.Mutex
	Border   PictureBookBorder
	Canvas   PictureBookCanvas
	Text     PictureBookText
	Options  *PictureBookOptions
	pages    int
	tmpfiles []string
}

func NewPictureBookDefaultOptions(ctx context.Context) (*PictureBookOptions, error) {

	opts := &PictureBookOptions{
		Orientation: "P",
		Size:        "letter",
		Width:       0.0,
		Height:      0.0,
		DPI:         150.0,
		Border:      0.01,
		Verbose:     false,
	}

	return opts, nil
}

func NewPictureBook(ctx context.Context, opts *PictureBookOptions) (*PictureBook, error) {

	var pdf *gofpdf.Fpdf

	if opts.Size == "custom" {

		sz := gofpdf.SizeType{
			Wd: opts.Width,
			Ht: opts.Height,
		}

		init := gofpdf.InitType{
			OrientationStr: opts.Orientation,
			UnitStr:        "in",
			SizeStr:        "",
			Size:           sz,
			FontDirStr:     "",
		}

		pdf = gofpdf.NewCustom(&init)

	} else {

		pdf = gofpdf.New(opts.Orientation, "in", opts.Size, "")
	}

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
			return nil, err
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

	border_top := 1.0 * opts.DPI
	border_bottom := border_top * 1.5
	border_left := border_top * 1.0
	border_right := border_top * 1.0

	canvas_w := page_w - (border_left + border_right)
	canvas_h := page_h - (border_top + border_bottom)

	pdf.SetAutoPageBreak(false, border_bottom)

	b := PictureBookBorder{
		Top:    border_top,
		Bottom: border_bottom,
		Left:   border_left,
		Right:  border_right,
	}

	c := PictureBookCanvas{
		Width:  canvas_w,
		Height: canvas_h,
	}

	tmpfiles := make([]string, 0)
	mu := new(sync.Mutex)

	pb := PictureBook{
		PDF:      pdf,
		Mutex:    mu,
		Border:   b,
		Canvas:   c,
		Text:     t,
		Options:  opts,
		pages:    0,
		tmpfiles: tmpfiles,
	}

	return &pb, nil
}

func (pb *PictureBook) AddPictures(ctx context.Context, paths []string) error {

	pictures, err := pb.GatherPictures(ctx, paths)

	if err != nil {
		return err
	}

	if pb.Options.Verbose {
		log.Printf("Count pictures gathered: %d\n", len(pictures))
	}

	if pb.Options.Sort != nil {

		sorted, err := pb.Options.Sort.Sort(ctx, pb.Options.Source, pictures)

		if err != nil {
			return err
		}

		pictures = sorted
	}

	for _, pic := range pictures {

		pb.Mutex.Lock()
		pb.pages += 1
		pagenum := pb.pages
		pb.Mutex.Unlock()

		err = pb.AddPicture(ctx, pagenum, pic.Path, pic.Caption)

		if err != nil && pb.Options.Verbose {
			log.Printf("Failed to add %s, %v", pic.Path, err)
		}
	}

	return nil
}

func (pb *PictureBook) GatherPictures(ctx context.Context, paths []string) ([]*picture.PictureBookPicture, error) {

	pictures := make([]*picture.PictureBookPicture, 0)

	var list func(context.Context, *blob.Bucket, string) error

	file := func(ctx context.Context, b *blob.Bucket, path string) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		abs_path := path

		if pb.Options.Filter != nil {

			ok, err := pb.Options.Filter.Continue(ctx, pb.Options.Source, abs_path)

			if err != nil {
				log.Printf("Failed to filter %s, %v\n", abs_path, err)
				return nil
			}

			if !ok {
				return nil
			}

			if pb.Options.Verbose {
				log.Printf("Include %s\n", abs_path)
			}
		}

		caption := ""

		if pb.Options.Caption != nil {

			txt, err := pb.Options.Caption.Text(ctx, pb.Options.Source, abs_path)

			if err != nil {
				log.Printf("Failed to generate caption text for %s, %v\n", abs_path, err)
				return nil
			}

			caption = txt
		}

		final_path := abs_path

		if pb.Options.PreProcess != nil {

			if pb.Options.Verbose {
				log.Printf("Processing %s\n", abs_path)
			}

			processed_path, err := pb.Options.PreProcess.Transform(ctx, pb.Options.Source, abs_path)

			if err != nil {
				log.Printf("Failed to process %s, %v\n", abs_path, err)
				return nil
			}

			pb.tmpfiles = append(pb.tmpfiles, processed_path)
			final_path = processed_path
		}

		pb.Mutex.Lock()
		defer pb.Mutex.Unlock()

		pic := &picture.PictureBookPicture{
			Source:  abs_path,
			Path:    final_path,
			Caption: caption,
		}

		pictures = append(pictures, pic)
		return nil
	}

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
				return err
			}

			path := obj.Key

			if obj.IsDir {

				err := list(ctx, bucket, path)

				if err != nil {
					return err
				}

				continue
			}

			err = file(ctx, bucket, path)

			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, path := range paths {

		err := list(ctx, pb.Options.Source, path)

		if err != nil {
			return nil, err
		}
	}

	return pictures, nil
}

func (pb *PictureBook) AddPicture(ctx context.Context, pagenum int, abs_path string, caption string) error {

	pb.Mutex.Lock()
	defer pb.Mutex.Unlock()

	im_r, err := pb.Options.Source.NewReader(ctx, abs_path, nil)

	if err != nil {
		return err
	}

	defer im_r.Close()

	im, format, err := util.DecodeImageFromReader(im_r)

	if err != nil {
		return err
	}

	// trap gofpdf "16-bit depth not supported in PNG file" errors

	if format == "png" {

		buf := new(bytes.Buffer)

		err = util.EncodeImage(im, format, buf)

		if err != nil {
			return err
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

			tmpfile_path, tmpfile_format, err := tempfile.TempFileWithImage(ctx, pb.Options.Source, im)

			if err != nil {
				return err
			}

			if pb.Options.Verbose {
				log.Printf("%s converted to a JPG (%s)\n", abs_path, tmpfile_path)
			}

			pb.tmpfiles = append(pb.tmpfiles, tmpfile_path)

			abs_path = tmpfile_path
			format = tmpfile_format
		}
	}

	dims := im.Bounds()

	w := float64(dims.Max.X)
	h := float64(dims.Max.Y)

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

			if pb.Options.Verbose {
				log.Printf("Rotate %s\b", abs_path)
			}

			new_im, err := rotate.RotateImageWithDegrees(ctx, im, 90.0)

			if err != nil {
				return err
			}

			im = new_im
			dims = im.Bounds()

			w = float64(dims.Max.X)
			h = float64(dims.Max.Y)

			// now save to disk...

			tmpfile_path, tmpfile_format, err := tempfile.TempFileWithImage(ctx, pb.Options.Source, im)

			if err != nil {
				return err
			}

			pb.tmpfiles = append(pb.tmpfiles, tmpfile_path)

			if pb.Options.Verbose {
				log.Printf("%s converted to a JPG (%s)\n", abs_path, tmpfile_path)
			}

			abs_path = tmpfile_path
			format = tmpfile_format
		}
	}

	info := pb.PDF.GetImageInfo(abs_path)

	if info == nil {

		opts := gofpdf.ImageOptions{
			ReadDpi:   false,
			ImageType: format,
		}

		r, err := pb.Options.Source.NewReader(ctx, abs_path, nil)

		if err != nil {
			return err
		}

		defer r.Close()

		info = pb.PDF.RegisterImageOptionsReader(abs_path, opts, r)

	}

	if info == nil {
		return errors.New("unable to determine info")
	}

	info.SetDpi(pb.Options.DPI)

	if pb.Options.Verbose {
		log.Printf("[%d] %s %02.f x %02.f\n", pagenum, abs_path, w, h)
	}

	if w == 0.0 || h == 0.0 {
		msg := fmt.Sprintf("[%d] %s has zero-sized dimension", pagenum, abs_path)
		return errors.New(msg)
	}

	x := pb.Border.Left
	y := pb.Border.Top

	_, line_h := pb.PDF.GetFontSize()

	max_w := pb.Canvas.Width
	max_h := pb.Canvas.Height - (pb.Text.Margin + line_h)

	if pb.Options.Verbose {
		log.Printf("[%d] canvas: %0.2f (%0.2f) x %0.2f (%0.2f) image: %0.2f x %0.2f\n", pagenum, max_w, pb.Canvas.Width, max_h, pb.Canvas.Height, w, h)
	}

	for {

		if w > max_w || h > max_h {

			// log.Printf("[%d] WTF 1 %0.2f x %0.2f (%0.2f x %0.2f) \n", pagenum, w, h, max_w, max_h)

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

		if w <= max_w && h <= max_h {
			break
		}
	}

	if w < max_w {

		padding := max_w - w
		x = x + (padding / 2.0)
	}

	// if max_h > max_w && h < (max_h - pb.Border.Top) {

	if h < (max_h - pb.Border.Top) {

		y = y + pb.Border.Top
	}

	if pb.Options.Verbose {
		log.Printf("[%d] final %0.2f x %0.2f (%0.2f x %0.2f)\n", pagenum, w, h, x, y)
	}

	pb.PDF.AddPage()

	// https://godoc.org/github.com/jung-kurt/gofpdf#ImageOptions

	opts := gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: format,
	}

	x = x / pb.Options.DPI
	y = y / pb.Options.DPI
	w = w / pb.Options.DPI
	h = h / pb.Options.DPI

	r_border := pb.Options.Border

	if r_border > 0.0 {
		pb.PDF.SetFillColor(0, 0, 0)
		pb.PDF.Rect((x - r_border), (y - r_border), (w + (r_border * 2)), (h + (r_border * 2)), "FD")
	}

	pb.PDF.ImageOptions(abs_path, x, y, w, h, false, opts, 0, "")

	if caption != "" {

		cur_x, cur_y := pb.PDF.GetXY()

		txt := caption

		txt_w := pb.PDF.GetStringWidth(txt)
		txt_h := line_h

		txt_w = txt_w + +(pb.Text.Margin * 2)
		txt_h = txt_h + +(pb.Text.Margin * 2)

		cur_x = (x - r_border)
		cur_y = (y - r_border) + (h + (r_border * 2))

		// please do this in the constructor...
		// (20171128/thisisaaronland)

		font_sz, _ := pb.PDF.GetFontSize()
		pb.PDF.SetFontSize(font_sz + 2)

		_, line_h := pb.PDF.GetFontSize()

		pb.PDF.SetFontSize(font_sz)

		txt_x := cur_x
		txt_y := cur_y + line_h

		if pb.Options.Verbose {
			log.Printf("[%d] text at %0.2f x %0.2f (%0.2f x %0.2f)\n", pagenum, txt_x, txt_y, txt_w, txt_h)
		}

		// pb.PDF.SetFillColor(255, 255, 255)
		// pb.PDF.Rect(txt_x, txt_y, txt_w, txt_h, "FD")

		pb.PDF.SetXY(txt_x, txt_y)
		// pb.PDF.Cell(txt_w, txt_h, txt)

		pb.PDF.SetLeftMargin(x)
		pb.PDF.SetRightMargin(pb.Border.Right / pb.Options.DPI)

		// please account for lack of utf-8 support (20171128/thisisaaronland)
		// https://github.com/jung-kurt/gofpdf/blob/cc7f4a2880e224dc55d15289863817df6d9f6893/fpdf_test.go#L1440-L1478
		// tr := pb.PDF.UnicodeTranslatorFromDescriptor("utf8")
		// txt = tr(txt)

		txt = unidecode.Unidecode(txt)

		html := pb.PDF.HTMLBasicNew()
		html.Write(line_h, txt)
	}

	return nil
}

func (pb *PictureBook) Save(ctx context.Context, path string) error {

	if pb.Options.Target == nil {
		return errors.New("Missing or invalid target bucket")
	}

	// move this out of here...

	defer func() {

		for _, path := range pb.tmpfiles {

			if pb.Options.Verbose {
				log.Printf("Remove tmp file '%s'\n", path)
			}

			err := pb.Options.Source.Delete(ctx, path)

			if err != nil {
				log.Printf("Failed to delete %s, %v\n", path, err)
			}
		}
	}()

	if pb.Options.Verbose {
		log.Printf("Save %s\n", path)
	}

	wr, err := pb.Options.Target.NewWriter(ctx, path, nil)

	if err != nil {
		return err
	}

	err = pb.PDF.Output(wr)

	if err != nil {
		return err
	}

	err = wr.Close()

	if err != nil {
		return err
	}

	return nil
}
