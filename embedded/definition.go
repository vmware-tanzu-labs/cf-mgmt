package embedded

import (
	"embed"
	_ "embed"
)

//go:embed files
var Files embed.FS
