package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
import "C"
import "runtime"

// Context contains the current execution context and the last raised error.
type Context struct {
	ccontext *C.nix_c_context
}

// NewContext creates a new context.
func NewContext() *Context {
	cctx := C.nix_c_context_create()
	runtime.SetFinalizer(cctx, finalizeContext)
	return &Context{cctx}
}

func finalizeContext(cctx *C.nix_c_context) {
	C.nix_c_context_free(cctx)
}

func (c *Context) lastError() error {
	cerr := C.nix_err_code(c.ccontext)
	return nixError(cerr, c)
}
