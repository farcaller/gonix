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

type State struct {
	store  *Store
	cstate *C.State
}

// NewState creates a new execution state.
func (s *Store) NewState(searchPath []string) *State {
	searchPathPtrs := make([]*C.char, len(searchPath))
	for i, str := range searchPath {
		cStr := C.CString(str)
		defer C.free(unsafe.Pointer(cStr))
		searchPathPtrs[i] = cStr
	}
	var cSearchPath **C.char
	if len(searchPath) > 0 {
		cSearchPath = (**C.char)(unsafe.Pointer(&searchPathPtrs[0]))
	}

	cstate := C.nix_state_create(s.context().ccontext, cSearchPath, s.cstore)
	runtime.SetFinalizer(cstate, finalizeState)
	return &State{s, cstate}
}

func finalizeState(cstate *C.State) {
	C.nix_state_free(cstate)
}

func (s *State) context() *Context {
	return s.store.context()
}

// EvalExpr evaluates the expression and returns the result.
func (s *State) EvalExpr(expr, path string) (*Value, error) {
	ret, err := NewValue(s)
	if err != nil {
		return nil, err
	}
	cexpr := C.CString(expr)
	cpath := C.CString(path)
	cerr := C.nix_expr_eval_from_string(s.context().ccontext, s.cstate, cexpr, cpath, ret.cvalue)
	err = nixError(cerr, s.context())
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Call calls a given function with a given argument and returns tne result.
func (s *State) Call(fun, argument *Value) (*Value, error) {
	ret, err := NewValue(s)
	if err != nil {
		return nil, err
	}
	cerr := C.nix_value_call(s.context().ccontext, s.cstate, fun.cvalue, argument.cvalue, ret.cvalue)
	err = nixError(cerr, s.context())
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// NewInt returns a Value containing an integer.
func (s *State) NewInt(i int64) (*Value, error) {
	v, err := NewValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetInt(i)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NewBool returns a Value containing a bool.
func (s *State) NewBool(b bool) (*Value, error) {
	v, err := NewValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetBool(b)
	if err != nil {
		return nil, err
	}
	return v, nil
}
