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
	"path/filepath"
	"strings"
	"sync"
)

const MM2INCH float64 = 25.4

type PictureBookOptions struct {
	Orientation  string
	Size         string
	Width        float64
	Height       float64
	DPI          float64
	Border       float64
	Bleed        float64
	MarginTop    float64
	MarginBottom float64
	MarginLeft   float64
	MarginRight  float64
	Filter       filter.Filter
	PreProcess   process.Process
	Caption      caption.Caption
	Sort         sort.Sorter
	FillPage     bool
	Verbose      bool
	OCRAFont     bool
	Source       *blob.Bucket
	Target       *blob.Bucket
	EvenOnly     bool
	OddOnly      bool
}

type PictureBookMargins struct {
	Top    float64
	Bottom float64
	Left   float64
	Right  float64
}

type PictureBookBorders struct {
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
	Borders  *PictureBookBorders
	Margins  *PictureBookMargins
	Canvas   PictureBookCanvas
	Text     PictureBookText
	Options  *PictureBookOptions
	pages    int
	tmpfiles []string
}

func NewPictureBookDefaultOptions(ctx context.Context) (*PictureBookOptions, error) {

	opts := &PictureBookOptions{
		Orientation:  "P",
		Size:         "letter",
		Width:        0.0,
		Height:       0.0,
		DPI:          150.0,
		Border:       0.01,
		Bleed:        0.0,
		MarginTop:    1.0,
		MarginBottom: 1.0,
		MarginLeft:   1.0,
		MarginRight:  1.0,
		Verbose:      false,
	}

	return opts, nil
}

func NewPictureBook(ctx context.Context, opts *PictureBookOptions) (*PictureBook, error) {

	var pdf *gofpdf.Fpdf

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
	}

	// log.Printf("%0.2f x %0.2f (%s)\n", opts.Width, opts.Height, opts.Size)

	sz := gofpdf.SizeType{
		Wd: opts.Width + (opts.Bleed * 2.0),
		Ht: opts.Height + (opts.Bleed * 2.0),
	}

	init := gofpdf.InitType{
		OrientationStr: opts.Orientation,
		UnitStr:        "in",
		SizeStr:        "",
		Size:           sz,
		FontDirStr:     "",
	}

	pdf = gofpdf.NewCustom(&init)

	/*
		} else {

			// TO DO: ACCOUNT FOR BLEED
			// func (f *Fpdf) GetPageSizeStr(sizeStr string) (size SizeType) {

			pdf = gofpdf.New(opts.Orientation, "in", opts.Size, "")
		}
	*/

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

	pb := PictureBook{
		PDF:      pdf,
		Mutex:    mu,
		Borders:  borders,
		Margins:  margins,
		Canvas:   canvas,
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

		var err error

		if pb.Options.EvenOnly {

			if pagenum%2 != 0 {
				pb.AddBlankPage(ctx, pagenum)
				pb.pages += 1
				pagenum = pb.pages
			}

			err = pb.AddPicture(ctx, pagenum, pic.Path, pic.Caption)

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

			err = pb.AddPicture(ctx, pagenum, pic.Path, pic.Caption)

		} else {
			err = pb.AddPicture(ctx, pagenum, pic.Path, pic.Caption)
		}

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

			if processed_path != "" {
				pb.tmpfiles = append(pb.tmpfiles, processed_path)
				final_path = processed_path
			}
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

func (pb *PictureBook) AddBlankPage(ctx context.Context, pagenum int) error {
	pb.PDF.AddPage()
	return nil
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

	if pb.Options.Verbose {
		log.Printf("[%d][%s] dimensions %0.2f x %0.2f\n", pagenum, abs_path, w, h)
	}

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
		log.Printf("[%d][%s] dimensions %02.f x %02.f\n", pagenum, abs_path, w, h)
	}

	if w == 0.0 || h == 0.0 {
		msg := fmt.Sprintf("[%d] %s has zero-sized dimension", pagenum, abs_path)
		return errors.New(msg)
	}

	// Remember: margins have been calculated inclusive of page bleeds

	margins := pb.Margins

	x := margins.Left
	y := margins.Top

	_, line_h := pb.PDF.GetFontSize()

	if pb.Options.Verbose {
		log.Printf("[%d][%s] margins, left and right %0.2f\n", pagenum, abs_path, (margins.Left + margins.Right))
		log.Printf("[%d][%s] margins, top and bottom %0.2f\n", pagenum, abs_path, (margins.Top + margins.Bottom))
		log.Printf("[%d][%s] margins, caption %0.2f\n", pagenum, abs_path, (pb.Text.Margin + line_h))
	}

	max_w := pb.Canvas.Width
	max_h := pb.Canvas.Height

	if pb.Options.Verbose {
		log.Printf("[%d][%s] max dimensions %0.2f (%0.2f) x %0.2f (%0.2f)\n", pagenum, abs_path, max_w, w, max_h, h)
	}

	for {

		if w >= max_w || h >= max_h {

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

		// TO DO: ENSURE ! h < max_h && ! w < max_w

		if w <= max_w && h <= max_h {
			break

			if h < max_h {
				h = max_h
			}

		}
	}

	// log.Printf("[%d][%s] max dimensions (1) %0.2f (%0.2f) H  %0.2f (%0.2f)\n", pagenum, abs_path, max_w, w, max_h, h)

	if w < max_w {
		padding := max_w - w
		x = x + (padding / 2.0)
	}

	if h < max_h {
		padding := max_h - h
		y = y + (padding / 2.0)
	}

	// log.Printf("[%d][%s] max dimensions (2) %0.2f (%0.2f) H  %0.2f (%0.2f)\n", pagenum, abs_path, max_w, w, max_h, h)

	if pb.Options.Verbose {
		log.Printf("[%d][%s] final %0.2f x %0.2f (%0.2f x %0.2f)\n", pagenum, abs_path, w, h, x, y)
	}

	pb.PDF.AddPage()

	if pb.Options.Verbose {
		log.Printf("[%d][%s] final dimensions %0.2f x %0.2f (%0.2f x %0.2f)\n", pagenum, abs_path, w, h, x, y)
	}

	// draw margins

	mx := x / pb.Options.DPI
	my := y / pb.Options.DPI
	mw := w / pb.Options.DPI
	mh := h / pb.Options.DPI

	if pb.Options.Verbose {
		log.Printf("[%d][%s] margin  %0.2f x %0.2f @ %0.2f x %0.2f\n", pagenum, abs_path, mx, my, mw, mh)
	}

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

		if pb.Options.Verbose {
			log.Printf("[%d][%s] border  %0.2f x %0.2f @ %0.2f x %0.2f\n", pagenum, abs_path, bx, by, bw, bh)
		}

		pb.PDF.SetFillColor(0, 0, 0)
		pb.PDF.Rect(bx, by, bw, bh, "FD")
	}

	// draw the image

	// https://godoc.org/github.com/jung-kurt/gofpdf#ImageOptions

	image_opts := gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: format,
	}

	image_x := x / pb.Options.DPI
	image_y := y / pb.Options.DPI
	image_w := w / pb.Options.DPI
	image_h := h / pb.Options.DPI

	if pb.Options.Verbose {
		// log.Printf("[%d][%s] image  %0.2f x %0.2f @ %0.2f x %0.2f\n", pagenum, abs_path, x, y, w, h)
		log.Printf("[%d][%s] image  %0.2f x %0.2f @ %0.2f x %0.2f\n", pagenum, abs_path, image_x, image_y, image_w, image_h)
	}

	pb.PDF.ImageOptions(abs_path, image_x, image_y, image_w, image_h, false, image_opts, 0, "")

	if caption != "" {

		txt := caption

		txt_w := pb.PDF.GetStringWidth(txt)
		txt_h := line_h

		txt_w = txt_w + pb.Text.Margin
		txt_h = txt_h + pb.Text.Margin

		// please do this in the constructor...
		// (20171128/thisisaaronland)

		font_sz, _ := pb.PDF.GetFontSize()
		pb.PDF.SetFontSize(font_sz + 2)

		_, line_h := pb.PDF.GetFontSize()

		if pb.Options.Verbose {
			log.Printf("[%d][%s] line height %0.2f\n", pagenum, abs_path, line_h)
		}

		pb.PDF.SetFontSize(font_sz)

		txt_x := ((x + w) / pb.Options.DPI) - txt_w
		txt_y := ((y + h) / pb.Options.DPI) + line_h

		if pb.Options.Verbose {
			log.Printf("[%d][%s] text at %0.2f x %0.2f (%0.2f x %0.2f)\n", pagenum, abs_path, txt_x, txt_y, txt_w, txt_h)
		}

		// pb.PDF.SetFillColor(255, 255, 255)
		// pb.PDF.Rect(txt_x, txt_y, txt_w, txt_h, "FD")

		pb.PDF.SetXY(txt_x, txt_y)

		// please account for lack of utf-8 support (20171128/thisisaaronland)
		// https://github.com/jung-kurt/gofpdf/blob/cc7f4a2880e224dc55d15289863817df6d9f6893/fpdf_test.go#L1440-L1478
		// tr := pb.PDF.UnicodeTranslatorFromDescriptor("utf8")
		// txt = tr(txt)

		txt = unidecode.Unidecode(txt)

		if pb.Options.Verbose {
			log.Printf("[%d][%s] caption '%s'\n", pagenum, abs_path, txt)
		}

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

			fname := filepath.Base(path)

			// This shouldn't be necessary and points to a larger problem
			// but this bandaid-fix will have to do for now...
			// (20210103/straup)

			if !strings.HasPrefix(fname, "picturebook-") {
				continue
			}

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
