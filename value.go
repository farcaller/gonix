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

func NewValue(state *State) (*Value, error) {
	cvalue := C.nix_alloc_value(state.context().ccontext, state.cstate)
	if cvalue == nil {
		return nil, NewContext().lastError()
	}
	runtime.SetFinalizer(&cvalue, func(v *unsafe.Pointer) {
		C.nix_gc_decref(state.context().ccontext, *v)
	})
	return &Value{state, cvalue}, nil
}

func (v *Value) context() *Context {
	return v.state.store.context()
}

func (v *Value) Force() error {
	cerr := C.nix_value_force(v.context().ccontext, v.state.cstate, v.cvalue)
	return nixError(cerr, v.context())
}

func (v *Value) GetType() int {
	t := C.nix_get_type(v.context().ccontext, v.cvalue)

	return (int)(t)
}

func (v *Value) GetString() (*string, error) {
	cctx := v.context().ccontext

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

func (v *Value) GetInt() (int64, error) {
	cctx := v.context().ccontext

	typ := v.GetType()
	if typ != C.NIX_TYPE_INT {
		return 0, fmt.Errorf("expected an int, got %v", typ)
	}

	i := C.nix_get_int(cctx, v.cvalue)
	if err := v.context().lastError(); err != nil {
		return 0, err
	}
	return int64(i), nil
}

func (v *Value) SetBool(b bool) error {
	cerr := C.nix_set_bool(v.context().ccontext, v.cvalue, C.bool(b))
	return nixError(cerr, v.context())
}

func (v *Value) SetString(s string) error {
	cerr := C.nix_set_string(v.context().ccontext, v.cvalue, C.CString(s))
	return nixError(cerr, v.context())
}

func (v *Value) SetPathString(ps string) error {
	cerr := C.nix_set_path_string(v.context().ccontext, v.cvalue, C.CString(ps))
	return nixError(cerr, v.context())
}

func (v *Value) SetFloat(f float64) error {
	cerr := C.nix_set_float(v.context().ccontext, v.cvalue, C.double(f))
	return nixError(cerr, v.context())
}

func (v *Value) SetInt(i int64) error {
	cerr := C.nix_set_int(v.context().ccontext, v.cvalue, C.longlong(i))
	return nixError(cerr, v.context())
}

func (v *Value) SetNull() error {
	cerr := C.nix_set_null(v.context().ccontext, v.cvalue)
	return nixError(cerr, v.context())
}

func (v *Value) SetExternalValue(ev externalValueHandle) error {
	cerr := C.nix_set_external(v.context().ccontext, v.cvalue, ev.cev)
	return nixError(cerr, v.context())
}
