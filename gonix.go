// Package gonix provides bindings to the nix APIs. It can work with any store
// supported by nix, including working with `/nix/store` directly or over the
// nix-daemon.
package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
import "C"

func init() {
	C.nix_libutil_init(nil)
	C.nix_libexpr_init(nil)
	C.nix_libstore_init(nil)
}
