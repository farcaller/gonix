package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
import "C"

func ExprEvalFromString(state *State, expr, path string) (*Value, error) {
	retval := NewValue(state)
	cexpr := C.CString(expr)
	cpath := C.CString(path)
	cerr := C.nix_expr_eval_from_string(state.Context().ccontext, state.cstate, cexpr, cpath, retval.cvalue)
	err := nixError(cerr)
	if err != nil {
		return nil, err
	}
	return retval, nil
}
