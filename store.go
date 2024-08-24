package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_store.h>
/*
typedef const char cchar_t;
void nixStoreBuildCallback_cgo(void * userdata, cchar_t * outname, cchar_t * out);
*/
import "C"

import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

// Store is the interface to the nix store.
type Store struct {
	ctx    *Context
	cstore *C.Store
}

// NewStore creates a new store. See [nix help-stores] for details on the store
// options.
//
// [nix help-stores]: https://nixos.org/manual/nix/stable/command-ref/new-cli/nix3-help-stores.html
func NewStore(ctx *Context, uri string, params map[string]string) (*Store, error) {
	cstore := C.nix_store_open(ctx.ccontext, nil, nil)
	if cstore == nil {
		return nil, ctx.lastError()
	}
	s := &Store{ctx, cstore}
	runtime.SetFinalizer(s, finalizeStore)
	return s, nil
}

func finalizeStore(store *Store) {
	C.nix_store_free(store.cstore)
}

func (s *Store) context() *Context {
	return s.ctx
}

//export nixStoreBuildCallback
func nixStoreBuildCallback(userdata unsafe.Pointer, outname, out *C.cchar_t) {
	h := cgo.Handle(userdata)
	bo := h.Value().(BuildOutputs)
	bo[C.GoString(outname)] = C.GoString(out)
}

// StorePath defines a path in the nix store.
type StorePath struct {
	store      *Store
	cstorepath *C.StorePath
}

// ParsePath parses a path in the nix store.
func (s *Store) ParsePath(path string) (*StorePath, error) {
	csp := C.nix_store_parse_path(s.ctx.ccontext, s.cstore, C.CString(path))
	err := s.context().lastError()
	if err != nil {
		return nil, err
	}
	return &StorePath{s, csp}, nil
}

// Valid returns true if the store path is valid for the owning store.
func (sp *StorePath) Valid() bool {
	return bool(C.nix_store_is_valid_path(sp.store.ctx.ccontext, sp.store.cstore, sp.cstorepath))
}

// BuildOutputs is a map for the derivation outputs to output paths.
type BuildOutputs map[string]string

// Build builds the store path, returning the realised outputs.
func (sp *StorePath) Build() (BuildOutputs, error) {
	bo := BuildOutputs{}
	boh := cgo.NewHandle(bo)
	defer boh.Delete()

	cerr := C.nix_store_realise(
		sp.store.ctx.ccontext,
		sp.store.cstore,
		sp.cstorepath,
		unsafe.Pointer(boh),
		(*[0]byte)(C.nixStoreBuildCallback_cgo),
	)
	err := nixError(cerr, sp.store.context())
	if err != nil {
		return nil, err
	}
	return bo, nil
}

// URI returns the URI of the store.
//
// Example:
//
//	uri, err := store.URI()
//
// Output:
//
//	"daemon", nil
func (s *Store) URI() (string, error) {
	var buf *[0]byte
	cerr := C.nix_store_get_uri(s.context().ccontext, s.cstore, buf, unsafe.Pointer(nil))
	if cerr != C.NIX_OK {
		return "", s.context().lastError()
	}
	return C.GoString((*C.char)(unsafe.Pointer(buf))), nil
}

// Version returns the version of the store.
//
// Example:
//
//	ver, err := store.Version()
//
// Output:
//
//	"2.13.5", nil
func (s *Store) Version() (string, error) {
	var buf *[0]byte
	cerr := C.nix_store_get_version(s.context().ccontext, s.cstore, buf, unsafe.Pointer(nil))
	if cerr != C.NIX_OK {
		return "", s.context().lastError()
	}
	return C.GoString((*C.char)(unsafe.Pointer(buf))), nil
}
