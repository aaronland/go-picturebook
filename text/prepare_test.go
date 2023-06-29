package text

import (
	"fmt"
	"testing"

	"github.com/jung-kurt/gofpdf"
)

func TestPrepareText(t *testing.T) {

	txt := "This page left intentionally blank.\nWoo woo\nFoobar\nIt’s good to see Mastodon and Bluesky showing a lot of life, but I will say that I was secretly hoping Twitter would die without a replacement and we’d start sending personal e-mails again. I miss those."

	max_w := 972.0

	var pdf *gofpdf.Fpdf

	sz := gofpdf.SizeType{
		Wd: 8.5,
		Ht: 11.0,
	}

	init := gofpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "in",
		SizeStr:        "",
		Size:           sz,
		FontDirStr:     "",
	}

	pdf = gofpdf.NewCustom(&init)

	pdf.SetFont("Helvetica", "", 8.0)

	prepped := PrepareText(pdf, txt, max_w)

	for i, ln := range prepped {
		fmt.Printf("%d '%s'\n", i, ln)
	}
}
