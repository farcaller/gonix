//go:build darwin

package gonix

// #cgo pkg-config: nix-expr-c
// #include <nix_api_util.h>
// #include <nix_api_value.h>
import "C"

func (v *Value) SetInt(i int64) error {
	cerr := C.nix_init_int(v.context().ccontext, v.cvalue, C.longlong(i))
	return nixError(cerr, v.context())
}
