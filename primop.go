package gonix

// #cgo pkg-config: nix-expr-c nix-main-c
// #include <stdlib.h>
// #include <stdio.h>
// #include <nix_api_util.h>
// #include <nix_api_value.h>
// #include <nix_api_expr.h>
// #include <nix_api_main.h>
/*
void nixPrimOp_cgo(void * user_data, nix_c_context * context, EvalState * state, nix_value ** args, nix_value * ret) {
	void nixPrimOp(void * user_data, nix_c_context * context, EvalState * state, nix_value ** args, nix_value * ret);
	nixPrimOp(user_data, context, state, args, ret);
}
void finalizePrimOp_cgo(void * obj, void * cd) {
	void finalizePrimOp(void * obj, void * cd);
	finalizePrimOp(obj, cd);
}
*/
import "C"

import (
	"errors"
	"runtime"
	"runtime/cgo"
	"unsafe"
)

// PrimOp is a wrapper around a nix primop.
type PrimOp struct {
	cprimop *C.PrimOp
}

type PrimOpFunc func(ctx *Context, state *State, args []*Value, ret *Value) error

type primOpHandle struct {
	numArgs int
	fun     PrimOpFunc
}

func NewPrimOp(ctx *Context, name string, args []string, doc string, fun PrimOpFunc) (*PrimOp, error) {
	if len(args) == 0 {
		return nil, errors.New("a function must have at least one argument")
	}

	argNamesPtrs := make([]*C.char, len(args)+1)
	for i, str := range args {
		cStr := C.CString(str)
		defer C.free(unsafe.Pointer(cStr))
		argNamesPtrs[i] = cStr
	}
	argNames := (**C.char)(unsafe.Pointer(&argNamesPtrs[0]))

	handle := cgo.NewHandle(primOpHandle{len(args), fun})
	handlePtr := &handle

	cprimop := C.nix_alloc_primop(ctx.ccontext, (*[0]byte)(C.nixPrimOp_cgo), C.int(len(args)), C.CString(name), argNames, C.CString(doc), unsafe.Pointer(handlePtr))
	if cprimop == nil {
		return nil, NewContext().lastError()
	}
	C.nix_gc_register_finalizer(unsafe.Pointer(cprimop), unsafe.Pointer(handlePtr), (*[0]byte)(C.finalizePrimOp_cgo))

	po := &PrimOp{cprimop}
	runtime.SetFinalizer(po, func(v *PrimOp) {
		C.nix_gc_decref(nil, unsafe.Pointer(v.cprimop))
	})

	return po, nil
}

var globalPrimops []*PrimOp = nil

// RegisterGlobalPrimOp registers the primop in the `builtins` attrset. Only applies
// to the [State]s created after the call.
func RegisterGlobalPrimOp(ctx *Context, name string, args []string, doc string, fun PrimOpFunc) error {
	po, err := NewPrimOp(ctx, name, args, doc, fun)
	if err != nil {
		return err
	}
	cerr := C.nix_register_primop(ctx.ccontext, po.cprimop)
	if err := nixError(cerr, ctx); err != nil {
		return err
	}
	globalPrimops = append(globalPrimops, po)
	return nil
}
