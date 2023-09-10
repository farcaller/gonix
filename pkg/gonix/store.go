package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <stdio.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
// void nixStoreBuildCallback_cgo(void * userdata, const char * outname, const char * out);
import "C"

import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

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

//export nixStoreBuildCallback
func nixStoreBuildCallback(userdata unsafe.Pointer, outname, out *C.char) {
	h := cgo.Handle(userdata)
	bo := h.Value().(BuildOutputs)
	bo[C.GoString(outname)] = C.GoString(out)
}

type StorePath struct {
	store      *Store
	cstorepath *C.StorePath
}

func (s *Store) ParsePath(path string) (*StorePath, error) {
	csp := C.nix_store_parse_path(s.ctx.ccontext, s.cstore, C.CString(path))
	err := s.Context().LastError()
	if err != nil {
		return nil, err
	}
	return &StorePath{s, csp}, nil
}

func (sp *StorePath) Valid() bool {
	return bool(C.nix_store_is_valid_path(sp.store.ctx.ccontext, sp.store.cstore, sp.cstorepath))
}

type BuildOutputs map[string]string

func (sp *StorePath) Build() (BuildOutputs, error) {
	bo := BuildOutputs{}
	cerr := C.nix_store_build(
		sp.store.ctx.ccontext,
		sp.store.cstore,
		sp.cstorepath,
		unsafe.Pointer(cgo.NewHandle(bo)),
		(*[0]byte)(C.nixStoreBuildCallback_cgo),
	)
	err := nixError(cerr, sp.store.Context())
	if err != nil {
		return nil, err
	}
	return bo, nil
}
