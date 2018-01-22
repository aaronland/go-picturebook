/* Adapted from https://github.com/bradfitz/camlistore */
package exifutil

import (
	"image"
	"image/draw"
	"log"
)

const (
	TopLeftSide     = 1
	TopRightSide    = 2
	BottomRightSide = 3
	BottomLeftSide  = 4
	LeftSideTop     = 5
	RightSideTop    = 6
	RightSideBottom = 7
	LeftSideBottom  = 8
)

// The FlipDirection type is used by the Flip option in DecodeOpts
// to indicate in which direction to flip an image.
type FlipDirection int

// FlipVertical and FlipHorizontal are two possible FlipDirections
// values to indicate in which direction an image will be flipped.
const (
	FlipVertical FlipDirection = 1 << iota
	FlipHorizontal
)

func Rotate(im image.Image, angle int) image.Image {
	var rotated *image.NRGBA
	// trigonometric (i.e counter clock-wise)
	switch angle {
	case 90:
		newH, newW := im.Bounds().Dx(), im.Bounds().Dy()
		rotated = image.NewNRGBA(image.Rect(0, 0, newW, newH))
		for y := 0; y < newH; y++ {
			for x := 0; x < newW; x++ {
				rotated.Set(x, y, im.At(newH-1-y, x))
			}
		}
	case -90:
		newH, newW := im.Bounds().Dx(), im.Bounds().Dy()
		rotated = image.NewNRGBA(image.Rect(0, 0, newW, newH))
		for y := 0; y < newH; y++ {
			for x := 0; x < newW; x++ {
				rotated.Set(x, y, im.At(y, newW-1-x))
			}
		}
	case 180, -180:
		newW, newH := im.Bounds().Dx(), im.Bounds().Dy()
		rotated = image.NewNRGBA(image.Rect(0, 0, newW, newH))
		for y := 0; y < newH; y++ {
			for x := 0; x < newW; x++ {
				rotated.Set(x, y, im.At(newW-1-x, newH-1-y))
			}
		}
	default:
		return im
	}
	return rotated
}

// flip returns a flipped version of the image im, according to
// the direction(s) in dir.
// It may flip the input im in place and return it, or it may allocate a
// new NRGBA (if im is an *image.YCbCr).
func Flip(im image.Image, dir FlipDirection) image.Image {
	if dir == 0 {
		return im
	}
	ycbcr := false
	var nrgba image.Image
	dx, dy := im.Bounds().Dx(), im.Bounds().Dy()
	di, ok := im.(draw.Image)
	if !ok {
		if _, ok := im.(*image.YCbCr); !ok {
			log.Printf("failed to flip image: input does not satisfy draw.Image")
			return im
		}
		// because YCbCr does not implement Set, we replace it with a new NRGBA
		ycbcr = true
		nrgba = image.NewNRGBA(image.Rect(0, 0, dx, dy))
		di, ok = nrgba.(draw.Image)
		if !ok {
			log.Print("failed to flip image: could not cast an NRGBA to a draw.Image")
			return im
		}
	}
	if dir&FlipHorizontal != 0 {
		for y := 0; y < dy; y++ {
			for x := 0; x < dx/2; x++ {
				old := im.At(x, y)
				di.Set(x, y, im.At(dx-1-x, y))
				di.Set(dx-1-x, y, old)
			}
		}
	}
	if dir&FlipVertical != 0 {
		for y := 0; y < dy/2; y++ {
			for x := 0; x < dx; x++ {
				old := im.At(x, y)
				di.Set(x, y, im.At(x, dy-1-y))
				di.Set(x, dy-1-y, old)
			}
		}
	}
	if ycbcr {
		return nrgba
	}
	return im
}

func ProcessOrientation(orientation int64) (angle int, flipMode FlipDirection, switchDimensions bool) {
	switch orientation {
	case TopLeftSide:
		// do nothing
	case TopRightSide:
		flipMode = 2
	case BottomRightSide:
		angle = 180
	case BottomLeftSide:
		angle = 180
		flipMode = 2
	case LeftSideTop:
		angle = -90
		flipMode = 2
		switchDimensions = true
	case RightSideTop:
		angle = -90
		switchDimensions = true
	case RightSideBottom:
		angle = 90
		flipMode = 2
		switchDimensions = true
	case LeftSideBottom:
		angle = 90
		switchDimensions = true
	}
	return
}
