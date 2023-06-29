package text

import (
	"fmt"
	"log"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func PrepareText(pdf *gofpdf.Fpdf, txt string, max_w float64) []string {

	prepped := make([]string, 0)

	for _, ln := range strings.Split(txt, "\n") {

		txt_w := pdf.GetStringWidth(ln)

		if txt_w <= max_w {
			prepped = append(prepped, ln)
			continue
		}

		words := strings.Split(ln, " ")
		count := len(words)

		if count == 1 {

			last_phrase := ""
			phrase := ""

			for i := 0; i < count; i++ {

				if i == 0 {
					phrase = words[i]
				} else {
					phrase = fmt.Sprintf("%s %s", phrase, words[i])
				}

				phrase_w := pdf.GetStringWidth(phrase)

				log.Printf("'%s' %f (%f)\n", phrase, phrase_w, max_w)

				if phrase_w > max_w {
					prepped = append(prepped, last_phrase)

					new_phrase := strings.Join(words[i:], " ")
					prepped = append(prepped, PrepareText(pdf, new_phrase, max_w)[:]...)
					break
				}

				last_phrase = phrase
			}
		}
	}

	return prepped

}
