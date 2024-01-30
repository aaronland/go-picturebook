# go-font-ocra

Go package exporting the OCR-A font.

## Examples

### go-fpdf

```
import (
	"github.com/go-pdf/fpdf"
	"github.com/sfomuseum/go-font-ocra"
)

func main() {

	font, _ := ocra.LoadFPDFFont()

	pdf := fpdf.New("P", "mm", "A4", "")
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
* https://pkg.go.dev/github.com/go-pdf/fpdf#Fpdf.AddFontFromBytes