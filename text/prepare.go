package text

import (
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
)

func PrepareText(pdf *fpdf.Fpdf, dpi float64, max_w float64, txt string) []string {

	return prepareTextWithSeparator(pdf, dpi, max_w, txt, "\n")
}

func prepareTextWithSeparator(pdf *fpdf.Fpdf, dpi float64, max_w float64, txt string, sep string) []string {

	prepped := make([]string, 0)

	lines := strings.Split(txt, "\n")

	for _, ln := range lines {

		ln_w := pdf.GetStringWidth(ln) * dpi

		if ln_w <= max_w {
			prepped = append(prepped, ln)
			continue
		}

		words := strings.Split(ln, " ")
		count := len(words)

		if count == 1 {

			prepped_ln := prepareTextWithSeparator(pdf, dpi, max_w, ln, " ")

			if len(prepped_ln) == 1 {
				prepped_ln = prepareTextWithLength(pdf, dpi, max_w, txt)
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

			phrase_w := pdf.GetStringWidth(phrase) * dpi

			if phrase_w > max_w {

				prepped = append(prepped, last_phrase)
				last_phrase = ""

				word_w := pdf.GetStringWidth(words[i]) * dpi

				if word_w > max_w {
					word_prepped := prepareTextWithLength(pdf, dpi, max_w, words[i])
					prepped = append(prepped, word_prepped[:]...)
				} else {

					new_phrase := strings.Join(words[i:], " ")
					phrase_prepped := prepareTextWithSeparator(pdf, dpi, max_w, new_phrase, " ")

					prepped = append(prepped, phrase_prepped[:]...)
					break
				}

			} else {
				last_phrase = phrase
			}
		}

		if last_phrase != "" {
			prepped = append(prepped, last_phrase)
		}
	}

	return prepped

}

func prepareTextWithLength(pdf *fpdf.Fpdf, dpi float64, max_w float64, txt string) []string {

	prepped := make([]string, 0)

	runes := []rune(txt)
	buf := make([]rune, 0)

	for _, r := range runes {

		buf = append(buf, r)
		ln := string(buf)

		ln_w := pdf.GetStringWidth(ln) * dpi

		if ln_w == max_w {
			prepped = append(prepped, ln)
			buf = make([]rune, 0)
		}
	}

	if len(buf) > 0 {
		ln := string(buf)
		prepped = append(prepped, ln)
	}

	return prepped
}
