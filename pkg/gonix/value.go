package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
import "C"
import (
	"runtime"
	"unsafe"
)

type Value struct {
	state  *State
	cvalue unsafe.Pointer
}

func NewValue(state *State) *Value {
	cvalue := C.nix_alloc_value(state.Context().ccontext, state.cstate)
	runtime.SetFinalizer(&cvalue, func(v *unsafe.Pointer) {
		C.nix_gc_decref(state.Context().ccontext, *v)
	})
	return &Value{state, cvalue}
}

func (v *Value) Force() error {
	cerr := C.nix_value_force(v.state.Context().ccontext, v.state.cstate, v.cvalue)
	return nixError(cerr)
}

func (v *Value) String() string {
	cctx := v.state.Context().ccontext

	s := C.nix_get_string(cctx, v.cvalue)
	if s == nil {
		return ""
	}
	return C.GoString(s)
}
