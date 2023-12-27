package text

import (
	"testing"
	
	"github.com/go-pdf/fpdf"
)

func TestPrepareText(t *testing.T) {

	tests := map[string]int{
		"This page left intentionally blank.\nWoo woo\nFoobar\nIt’s good to see Mastodon and Bluesky showing a lot of life, but I will say that I was secretly hoping Twitter would die without a replacement and we’d start sending personal e-mails again. I miss those.":                                                                                                                                                                                                                                                                         5,
		`<path fill="#e0f0ff" stroke="#000000" d="M254,328L258,324L260,324L264,324L267,324L269,324L271,325L272,327L272,330L273,333L273,337L273,341L274,343L274,345L275,346L275,347L275,348L273,349L269,350L266,351L262,352L260,352L258,353L260,353L263,353L266,353L269,350L271,348L272,345L273,342L273,339L274,337L274,335L275,333L275,330Z" stroke-opacity="0" fill-opacity="0.7" transform="matrix(1,0,0,1,0,0)" style="-webkit-tap-highlight-color: rgba(0, 0, 0, 0); stroke-opacity: 0; fill-opacity: 0.7; -webkit-user-select: text;"></path>`: 5,
	}

	max_w := 972.0

	var pdf *fpdf.Fpdf

	sz := fpdf.SizeType{
		Wd: 8.5,
		Ht: 11.0,
	}

	init := fpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "in",
		SizeStr:        "",
		Size:           sz,
		FontDirStr:     "",
	}

	pdf = fpdf.NewCustom(&init)

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
