package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
import "C"
import (
	"fmt"
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
	return nixError(cerr, v.state.Context())
}

func (v *Value) GetType() int {
	cctx := v.state.Context().ccontext

	t := C.nix_get_type(cctx, v.cvalue)

	return (int)(t)
}

func (v *Value) GetString() (*string, error) {
	cctx := v.state.Context().ccontext

	typ := v.GetType()
	if typ != C.NIX_TYPE_STRING {
		return nil, fmt.Errorf("expected a string, got %v", typ)
	}

	s := C.nix_get_string(cctx, v.cvalue)
	if s == nil {
		return nil, nil
	}
	str := C.GoString(s)
	return &str, nil
}
