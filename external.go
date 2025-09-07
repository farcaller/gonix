package gonix

// #cgo pkg-config: nix-expr-c
// #include <nix_api_external.h>
/*
void externalValuePrint_cgo(void * self, nix_printer * printer) {
	void externalValuePrint(void * self, nix_printer * printer);
	externalValuePrint(self, printer);
}
void externalValueShowType_cgo(void * self, nix_string_return * res) {
	void externalValueShowType(void * self, nix_string_return * res);
	externalValueShowType(self, res);
}
void externalValueTypeOf_cgo(void * self, nix_string_return * res) {
	void externalValueTypeOf(void * self, nix_string_return * res);
	externalValueTypeOf(self, res);
}
void externalValueCoerceToString_cgo(void * self, nix_string_context * c, int coerceMore, int copyToStore, nix_string_return * res) {
	void externalValueCoerceToString(void * self, nix_string_context * c, int coerceMore, int copyToStore, nix_string_return * res);
	externalValueCoerceToString(self, c, coerceMore, copyToStore, res);
}
int externalValueEqual_cgo(void * self, void * other) {
	int externalValueEqual(void * self, void * other);
	return externalValueEqual(self, other);
}
void externalValuePrintValueAsJSON_cgo(void * self, EvalState * state, int strict, nix_string_context * c, bool copyToStore, nix_string_return * res) {
	void externalValuePrintValueAsJSON(void * self, EvalState * state, int strict, nix_string_context * c, bool copyToStore, nix_string_return * res);
	externalValuePrintValueAsJSON(self, state, strict, c, copyToStore, res);
}
void finalizeExternalValue_cgo(void * obj, void * cd) {
	void finalizeExternalValue(void * obj, void * cd);
	finalizeExternalValue(obj, cd);
}
*/
import "C"

import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

// ExternalValueProvider is the interface that external values must implement to
// be passed into the nix [Value].
type ExternalValueProvider interface {
	// Print is called when printing the external value.
	Print() string

	// ShowType is called when `:t` is invoked on the value.
	ShowType() string

	// TypeOf is called when `builtins.typeOf` is invoked on the value.
	TypeOf() string

	// CoerceToString is called when "${str}" or `builtins.toString` is invoked on the value.
	//
	// if coerceMore is true, CoerceToString was invoked trough a `builtins.toString` call
	// (which also converts nulls, integers, booleans and lists to a string).
	//
	// if copyToStore is true, referenced paths are copied to the Nix store as a side effect.
	CoerceToString(addContext func(string) error, copyMore, copyToStore bool) (string, error)

	// Equal is called for a comparison of two external values. If `other` is not an external
	// cvalue, Equal isn't called and instead `false` is always returned.
	Equal(other ExternalValueProvider) bool

	// PrintValueAsJSON is called when the value is converted to JSON. The result must
	// be valid JSON.
	PrintValueAsJSON(strict, copyToStore bool) string

	// TODO: this api seems somewhat convoluted
	// PrintValueAsXML(strict, location, copyToStore bool) string
}

// ExternalValue is a wrapper around a go cvalue passed into nix.
type ExternalValue struct {
	cev *C.ExternalValue
}

var desc C.NixCExternalValueDesc

func init() {
	desc = C.NixCExternalValueDesc{
		print:            (*[0]byte)(C.externalValuePrint_cgo),
		showType:         (*[0]byte)(C.externalValueShowType_cgo),
		typeOf:           (*[0]byte)(C.externalValueTypeOf_cgo),
		coerceToString:   (*[0]byte)(C.externalValueCoerceToString_cgo),
		equal:            (*[0]byte)(C.externalValueEqual_cgo),
		printValueAsJSON: (*[0]byte)(C.externalValuePrintValueAsJSON_cgo),
	}
}

// NewExternalValue returns a Value containing an external value.
func (state *State) NewExternalValue(val ExternalValueProvider) (*Value, error) {
	h := cgo.NewHandle(val)
	cev := C.nix_create_external_value(state.context().ccontext, &desc, unsafe.Pointer(h))
	if cev == nil {
		return nil, NewContext().lastError()
	}
	C.nix_gc_register_finalizer(unsafe.Pointer(cev), unsafe.Pointer(h), (*[0]byte)(C.finalizeExternalValue_cgo))

	ev := &ExternalValue{cev}
	runtime.SetFinalizer(ev, func(v *ExternalValue) {
		C.nix_gc_decref(nil, unsafe.Pointer(v.cev))
	})

	v, err := newValue(state)
	if err != nil {
		return nil, err
	}
	err = v.SetExternalValue(ev)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Content returns the [ExternalValueProvider] of this value.
func (ev *ExternalValue) Content() ExternalValueProvider {
	// TODO: context?
	hp := C.nix_get_external_value_content(nil, ev.cev)
	if hp == nil {
		return nil
	}
	h := cgo.Handle(hp)
	return h.Value().(ExternalValueProvider)
}

func (ev ExternalValue) String() string {
	return ev.Content().Print()
}
