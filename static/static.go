// package static provides an `embed.FS` containing JavaScript and WebAssembly binaries used by the go-sfomuseum-airfield-wasm tools and methods.
package static

import (
	"embed"
)

//go:embed wasm/*
var FS embed.FS
