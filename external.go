package gonix

// #cgo pkg-config: nix-expr-c
// #include <nix_api_external.h>
/* void externalValuePrint_cgo(void * self, nix_printer * printer) {
	void externalValuePrint(void * self, nix_printer * printer);
	externalValuePrint(self, printer);
}
void externalValueCoerceToString_cgo(void * self, nix_string_context * c, int coerceMore, int copyToStore, nix_string_return * res) {
	void externalValueCoerceToString(void * self, nix_string_context * c, int coerceMore, int copyToStore, nix_string_return * res);
	externalValueCoerceToString(self, c, coerceMore, copyToStore, res);
}
*/
import "C"
import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

type ExternalValue interface {
	Print() string
	CoerceToString(addContext func(string) error, copyMore, copyToStore bool) (string, error)
}

type externalValue struct {
	ctx *Context
	v   ExternalValue
}

type externalValueHandle struct {
	goev externalValue
	cev  *C.ExternalValue
}

var desc C.NixCExternalValueDesc

func init() {
	desc = C.NixCExternalValueDesc{
		print:          (*[0]byte)(C.externalValuePrint_cgo),
		coerceToString: (*[0]byte)(C.externalValueCoerceToString_cgo),
	}
}

func (state *State) NewExternalValue(val ExternalValue) (*Value, error) {
	ev := externalValue{state.context(), val}
	cev := C.nix_create_external_value(state.context().ccontext, &desc, unsafe.Pointer(cgo.NewHandle(ev)))
	if cev == nil {
		return nil, NewContext().lastError()
	}
	// It's a mess of ownership so just to be on the safe side:
	// h refs both the val (that we sent into the C api via a cgo.Handle) and the
	// C.ExternalValue. We pass the handle into the Value, which (supposedly)
	// refcounts the C.ExternalValue on its own, but it will also ref the handle.
	// Now if everything works well, the Value retains a reference to the go's
	// val. So everything works. I think.
	//
	// What a mess.
	h := externalValueHandle{ev, cev}
	runtime.SetFinalizer(&h, func(v *externalValueHandle) {
		C.nix_gc_decref(state.context().ccontext, unsafe.Pointer(v.cev))
	})

	v, err := NewValue(state)
	if err != nil {
		return nil, err
	}
	err = v.SetExternalValue(h)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (ev *externalValue) print(pr *C.nix_printer) error {
	cerr := C.nix_external_print(ev.ctx.ccontext, pr, C.CString(ev.v.Print()))
	return nixError(cerr, ev.ctx)
}

func (ev *externalValue) coerceToString(c *C.nix_string_context, coerceMore, copyToStore C.int, res *C.nix_string_return) error {
	addCtx := func(s string) error {
		cerr := C.nix_external_add_string_context(ev.ctx.ccontext, c, C.CString(s))
		return nixError(cerr, ev.ctx)
	}
	ret, err := ev.v.CoerceToString(addCtx, coerceMore != 0, copyToStore != 0)
	if err != nil {
		return err
	}
	C.nix_set_string_return(res, C.CString(ret))
	return nil
}
