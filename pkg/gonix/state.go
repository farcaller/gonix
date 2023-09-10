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

func NewState(store *Store, searchPath []string) *State {
	searchPathPtrs := make([]*C.char, len(searchPath))
	for i, str := range searchPath {
		cStr := C.CString(str)
		defer C.free(unsafe.Pointer(cStr)) // Free the C string when done
		searchPathPtrs[i] = cStr
	}
	var csearchPath **C.char
	if len(searchPath) > 0 {
		csearchPath = (**C.char)(unsafe.Pointer(&searchPathPtrs[0]))
	}

	cstate := C.nix_state_create(store.Context().ccontext, csearchPath, store.cstore)
	runtime.SetFinalizer(cstate, finalizeState)
	return &State{store, cstate}
}

func finalizeState(cstate *C.State) {
	C.nix_state_free(cstate)
}

func (s *State) Context() *Context {
	return s.store.Context()
}

func (s *State) EvalExpr(expr, path string) (*Value, error) {
	retval := NewValue(s)
	cexpr := C.CString(expr)
	cpath := C.CString(path)
	cerr := C.nix_expr_eval_from_string(s.Context().ccontext, s.cstate, cexpr, cpath, retval.cvalue)
	err := nixError(cerr, s.Context())
	if err != nil {
		return nil, err
	}
	return retval, nil
}
