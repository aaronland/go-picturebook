package ocra

import (
	"fmt"

	"github.com/sfomuseum/go-font-ocra/fonts"
)

type FPDFFont struct {
	Family string
	Style  string
	JSON   []byte
	Z      []byte
}

func LoadFPDFFont() (*FPDFFont, error) {

	body_json, err := fonts.FS.ReadFile("OCRA.json")

	if err != nil {
		return nil, fmt.Errorf("Failed to read OCRA.json, %w", err)
	}

	body_z, err := fonts.FS.ReadFile("OCRA.z")

	if err != nil {
		return nil, fmt.Errorf("Failed to read OCRA.z, %w", err)
	}

	f := &FPDFFont{
		Family: "OCRA",
		Style:  "",
		JSON:   body_json,
		Z:      body_z,
	}

	return f, nil
}
