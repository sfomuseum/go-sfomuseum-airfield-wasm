package http

import (
	"net/http"

	aa_static "github.com/aaronland/go-http-static"
	"github.com/sfomuseum/go-sfomuseum-airfield-wasm/static"
)

// WASMOptions provides a list of JavaScript and CSS link to include with HTML output.
type WASMOptions struct {
	Prefix string
}

// Return a *WASMOptions struct with default paths and URIs.
func DefaultWASMOptions() *WASMOptions {
	opts := &WASMOptions{}
	return opts
}

// Append all the files in the net/http FS instance containing the embedded WASM assets to an *http.ServeMux instance.
func AppendAssetHandlers(mux *http.ServeMux, opts *WASMOptions) error {

	return aa_static.AppendStaticAssetHandlersWithPrefix(mux, static.FS, opts.Prefix)
}
