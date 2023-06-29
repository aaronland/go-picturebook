package text

import (
	"fmt"
	"log"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func PrepareText(pdf *gofpdf.Fpdf, max_w float64, txt string) []string {

	return prepareTextWithSeparator(pdf, max_w, txt, "\n")
}

func prepareTextWithSeparator(pdf *gofpdf.Fpdf, max_w float64, txt string, sep string) []string {

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

			prepped_ln := prepareTextWithSeparator(pdf, max_w, ln, " ")

			if len(prepped_ln) == 1 {
				prepped_ln = prepareTextWithLength(pdf, max_w, txt)
			}

			prepped = append(prepped, prepped_ln[:]...)
			continue
		}

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
				phrase_prepped := PrepareText(pdf, max_w, new_phrase)

				prepped = append(prepped, phrase_prepped[:]...)
				break
			}

			last_phrase = phrase
		}
	}

	return prepped

}

func prepareTextWithLength(pdf *gofpdf.Fpdf, max_w float64, txt string) []string {

	prepped := make([]string, 0)
	// Please write me...
	return prepped
}
