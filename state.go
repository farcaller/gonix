package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_expr.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// State is the execution state in a given [Store].
type State struct {
	_store *Store
	ctx    *Context
	cstate *C.EvalState
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
	if cstate == nil {
		return nil
	}
	state := &State{s, s.context(), cstate}
	runtime.SetFinalizer(state, finalizeState)
	return state
}

func finalizeState(state *State) {
	C.nix_state_free(state.cstate)
}

func (s *State) context() *Context {
	return s.ctx
}

// EvalExpr evaluates the expression and returns the result.
func (s *State) EvalExpr(expr, path string) (*Value, error) {
	ret, err := newValue(s)
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
	ret, err := newValue(s)
	if err != nil {
		return nil, err
	}

	var carg unsafe.Pointer
	if argument != nil {
		carg = argument.cvalue
	}

	cerr := C.nix_value_call(s.context().ccontext, s.cstate, fun.cvalue, carg, ret.cvalue)
	err = nixError(cerr, s.context())
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// NewInt returns a Value containing an integer.
func (s *State) NewInt(i int64) (*Value, error) {
	v, err := newValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetInt(i)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NewInt returns a Value containing a float.
func (s *State) NewFloat(f float64) (*Value, error) {
	v, err := newValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetFloat(f)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NewBool returns a Value containing a bool.
func (s *State) NewBool(b bool) (*Value, error) {
	v, err := newValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetBool(b)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NewString returns a Value containing a string.
func (s *State) NewString(st string) (*Value, error) {
	v, err := newValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetString(st)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NewPath returns a Value containing a path.
func (s *State) NewPath(p string) (*Value, error) {
	v, err := newValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetPath(p)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NewNull returns a Value containing a null.
func (s *State) NewNull() (*Value, error) {
	v, err := newValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetNull()
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NewList returns a Value containing a list of passed items.
func (s *State) NewList(items []*Value) (*Value, error) {
	v, err := newValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetList(items)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NewAttrs returns a Value containing an attrset.
func (s *State) NewAttrs(attrs map[string]*Value) (*Value, error) {
	v, err := newValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetAttrs(attrs)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// NewPrimOp returns a Value containing a [PrimOp].
func (s *State) NewPrimOp(op *PrimOp) (*Value, error) {
	v, err := newValue(s)
	if err != nil {
		return nil, err
	}
	err = v.SetPrimOp(op)
	if err != nil {
		return nil, err
	}
	return v, nil
}
