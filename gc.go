package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_expr.h>
import "C"

// GcNow triggers the nix garbage collector manually.
func GcNow() {
	C.nix_gc_now()
}
