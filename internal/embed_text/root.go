package embed_text

import _ "embed"

//go:embed root/root.txt
var Root string

//go:embed root/version.txt
var Version string
