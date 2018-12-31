package flags

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

type RGBAColor []color.Color

func (c *RGBAColor) String() string {
	return fmt.Sprintf("%v", *c)
}

func (c *RGBAColor) Set(value string) error {

	rgb := strings.Split(value, ",")

	if len(rgb) < 3 {
		return errors.New("Invalid R,G,B (,A) count")
	}

	r, err := strconv.Atoi(rgb[0])

	if err != nil {
		return err
	}

	g, err := strconv.Atoi(rgb[1])

	if err != nil {
		return err
	}

	b, err := strconv.Atoi(rgb[2])

	if err != nil {
		return err
	}

	var a int

	if len(rgb) == 4 {

		a, err = strconv.Atoi(rgb[3])

		if err != nil {
			return err
		}

	} else {
		a = 1
	}

	clr := color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}

	*c = append(*c, clr)
	return nil
}
