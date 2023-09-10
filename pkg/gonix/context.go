package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
import "C"
import "runtime"

type Context struct {
	ccontext *C.nix_c_context
}

func NewContext() *Context {
	cctx := C.nix_c_context_create()
	runtime.SetFinalizer(cctx, finalizeContext)
	return &Context{cctx}
}

func finalizeContext(cctx *C.nix_c_context) {
	C.nix_c_context_free(cctx)
}
