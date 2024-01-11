package contour

import (
	"context"
	"image"

	"github.com/fogleman/colormap"
	"github.com/fogleman/contourmap"
	"github.com/fogleman/gg"
)

// ContourImageSVG generate contour data derived from 'im' and draws it to a new `image.Image` instances.
// The scale and number of contours is adjusted relative to the 'scale' and 'n' values respectively.
func ContourImage(ctx context.Context, im image.Image, n int, scale float64) (image.Image, error) {

	m := contourmap.FromImage(im).Closed()
	z0 := m.Min
	z1 := m.Max

	w := int(float64(m.W) * scale)
	h := int(float64(m.H) * scale)

	dc := gg.NewContext(w, h)
	dc.SetRGB(1, 1, 1)
	dc.SetColor(colormap.ParseColor("FFFFFF"))
	dc.Clear()
	dc.Scale(scale, scale)

	for i := 0; i < n; i++ {

		t := float64(i) / (float64(n) - 1)
		z := z0 + (z1-z0)*t
		contours := m.Contours(z + 1e-9)

		for _, c := range contours {

			dc.NewSubPath()

			for _, p := range c {
				dc.LineTo(p.X, p.Y)
			}
		}

		dc.SetRGB(0, 0, 0)
		dc.SetLineWidth(z)
		dc.Stroke()
	}

	return dc.Image(), nil
}
