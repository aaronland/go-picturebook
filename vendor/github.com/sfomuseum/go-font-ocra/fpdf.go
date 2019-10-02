package ocra

import (
	"github.com/sfomuseum/go-font-ocra/fonts"
)

type FPDFFont struct {
	Family string
	Style string
	JSON []byte
	Z []byte
}

func LoadFPDFFont() (*FPDFFont, error) {

	body_json, err := fonts.Asset("OCRA.json")

	if err != nil {
		return nil, err
	}
	
	body_z, err := fonts.Asset("OCRA.z")	

	if err != nil {
		return nil, err		
	}

	f := &FPDFFont{
		Family: "OCRA",
		Style: "",
		JSON: body_json,
		Z: body_z,
	}

	return f, nil
}
