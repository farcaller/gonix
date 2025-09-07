package gonix

// #cgo pkg-config: nix-expr-c nix-main-c
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
// #include <nix_api_main.h>
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
func nixPrimOp(funh unsafe.Pointer, cctx *C.nix_c_context, cstate *C.EvalState, cargs unsafe.Pointer, cret *C.nix_value) {
	h := (*cgo.Handle)(funh)
	poh := h.Value().(primOpHandle)

	ctx := &Context{cctx}
	state := &State{nil, ctx, cstate}

	doError := func(err error) {
		C.nix_set_err_msg(cctx, C.NIX_ERR_UNKNOWN, C.CString(err.Error()))
	}

	args := make([]*Value, 0, poh.numArgs)
	for idx := 0; idx < poh.numArgs; idx++ {
		cargPtr := (**C.nix_value)(unsafe.Pointer(uintptr(cargs) + uintptr(uintptr(idx)*unsafe.Sizeof(cret))))

		val, err := wrapValue(state, *cargPtr)
		if err != nil {
			doError(fmt.Errorf("failed to wrap cvalue during a primop call: %v", err))
		}
		err = val.Force()
		if err != nil {
			doError(fmt.Errorf("failed to force cvalue during a primop call: %v", err))
		}
		args = append(args, val)
	}

	retValue, err := wrapValue(state, cret)
	if err != nil {
		doError(fmt.Errorf("failed to wrap cvalue during a primop call: %v", err))
		return
	}

	err = poh.fun(ctx, state, args, retValue)
	if err != nil {
		doError(fmt.Errorf("error during Go callback: %v", err))
	}
}
