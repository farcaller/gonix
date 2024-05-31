package gonix

// #cgo pkg-config: nix-expr-c
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
import "C"

import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

//export finalizePrimOp
func finalizePrimOp(obj, cd unsafe.Pointer) {
	h := cgo.Handle(cd)
	h.Delete()
}

//export nixPrimOp
func nixPrimOp(funh unsafe.Pointer, cctx *C.nix_c_context, cstate *C.EvalState, cargs unsafe.Pointer, cret unsafe.Pointer) {
	h := cgo.Handle(funh)
	poh := h.Value().(primOpHandle)

	ctx := &Context{cctx}
	state := &State{nil, ctx, cstate}

	args := make([]*Value, poh.numArgs)
	for idx := 0; idx < poh.numArgs; idx++ {
		cargPtr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(cargs) + uintptr(uintptr(idx)*unsafe.Sizeof(cargs))))
		carg := *cargPtr

		val, err := wrapValue(state, unsafe.Pointer(carg))
		if err != nil {
			panic(fmt.Errorf("failed to wrap value during a primop call: %v", err))
		}
		err = val.Force()
		if err != nil {
			panic(fmt.Errorf("failed to force value during a primop call: %v", err))
		}
		args[idx] = val
	}

	ret := poh.fun(ctx, state, args...)
	if ret != nil {
		C.nix_copy_value(cctx, unsafe.Pointer(cret), ret.cvalue)
	}
}
