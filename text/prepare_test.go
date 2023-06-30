package text

import (
	"testing"

	"github.com/jung-kurt/gofpdf"
)

func TestPrepareText(t *testing.T) {

	tests := map[string]int{
		"This page left intentionally blank.\nWoo woo\nFoobar\nIt’s good to see Mastodon and Bluesky showing a lot of life, but I will say that I was secretly hoping Twitter would die without a replacement and we’d start sending personal e-mails again. I miss those.": 5,
	}

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

	dpi := 150.0

	for txt, expected_count := range tests {

		prepped := PrepareText(pdf, dpi, max_w, txt)
		count := len(prepped)

		if count != expected_count {
			t.Fatalf("Unexpected count %d (expected %d) for '%s'", count, expected_count, txt)
		}
	}

}
