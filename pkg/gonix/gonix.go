package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
import "C"

func init() {
	C.nix_libexpr_init(nil)
	C.nix_libstore_init(nil)
	C.nix_libutil_init(nil)
}
