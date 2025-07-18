package contour

import (
	"context"
	"fmt"
	"image"
	"io"

	"github.com/fogleman/contourmap"
)

// ContourImageSVG generate contour data derived from 'im' and writes to 'wr' as SVG data. The scale and number of contours
// is adjusted relative to the 'scale' and 'n' values respectively.
func ContourImageSVG(ctx context.Context, wr io.Writer, im image.Image, n int, scale float64) error {

	m := contourmap.FromImage(im).Closed()
	z0 := m.Min
	z1 := m.Max

	w := int(float64(m.W) * scale)
	h := int(float64(m.H) * scale)

	fmt.Fprintf(wr, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, w, h, w, h)

	for i := 0; i < n; i++ {

		t := float64(i) / (float64(n) - 1)
		z := z0 + (z1-z0)*t
		contours := m.Contours(z + 1e-9)

		for _, c := range contours {

			fmt.Fprintf(wr, `<path stroke="%s" stroke-width="%02f" stroke-opacity="1" fill-opacity="0" d="M`, "#000000", z)

			for i, p := range c {

				if i > 0 {
					fmt.Fprintf(wr, `L`)
				}

				fmt.Fprintf(wr, `%d,%d`, int(p.X), int(p.Y))
			}

			fmt.Fprintf(wr, `Z"></path>`)
		}
	}

	fmt.Fprintf(wr, `</svg>`)
	return nil
}
