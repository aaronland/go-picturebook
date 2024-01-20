package fonts

import (
	"embed"
)

//go:embed *.json *.otf *.ttf *.z
var FS embed.FS
