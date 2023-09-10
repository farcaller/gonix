package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
import "C"
import "runtime"

type Store struct {
	ctx    *Context
	cstore *C.Store
}

func NewStore(ctx *Context, uri string, params map[string]string) *Store {
	cstore := C.nix_store_open(ctx.ccontext, nil, nil)
	runtime.SetFinalizer(cstore, finalizeStore)
	return &Store{ctx, cstore}
}

func finalizeStore(cstore *C.Store) {
	C.nix_store_unref(cstore)
}

func (s *Store) Context() *Context {
	return s.ctx
}
