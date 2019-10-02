# go-font-ocra

Work in progress

## Examples

### go-fpdf

```
import (
	"github.com/jung-kurt/gofpdf"
	"github.com/sfomuseum/go-font-ocra"
)

func main() {

	font, _ := ocra.LoadFPDFFont()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddFontFromBytes(font.Family, font.Style, font.JSON, font.Z)

	pdf.AddPage()
	pdf.SetFont(font.Family, "", 16)

	pdf.Cell(0, 10, "This page left intentionally blank.")
	pdf.OutputFileAndClose("test.pdf")
}
```

_Error handling removed for the sake of brevity._

## See also

* https://github.com/opensourcedesign/fonts/tree/master/OCR
* https://sourceforge.net/projects/ocr-a-font/files/OCR-A/1.0/
* https://godoc.org/github.com/jung-kurt/gofpdf#example-Fpdf-AddFontFromBytes
* https://godoc.org/github.com/jung-kurt/gofpdf#hdr-Nonstandard_Fonts