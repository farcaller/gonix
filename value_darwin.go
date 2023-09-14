//go:build darwin

package gonix

// #cgo pkg-config: nix-expr-c
// #include <nix_api_value.h>
import "C"

func (v *Value) SetInt(i int64) error {
	v.ev = nil
	cerr := C.nix_set_int(v.context().ccontext, v.cvalue, C.longlong(i))
	return nixError(cerr, v.context())
}
