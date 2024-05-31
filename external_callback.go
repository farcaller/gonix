package gonix

// #cgo pkg-config: nix-expr-c
// #include <nix_api_external.h>
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

//export externalValuePrint
func externalValuePrint(self unsafe.Pointer, pr *C.nix_printer) {
	h := cgo.Handle(self)
	v := h.Value().(ExternalValueProvider)

	C.nix_external_print(nil, pr, C.CString(v.Print()))
}

//export externalValueShowType
func externalValueShowType(self unsafe.Pointer, res *C.nix_string_return) {
	h := cgo.Handle(self)
	v := h.Value().(ExternalValueProvider)

	C.nix_set_string_return(nil, C.CString(v.ShowType()))
}

//export externalValueTypeOf
func externalValueTypeOf(self unsafe.Pointer, res *C.nix_string_return) {
	h := cgo.Handle(self)
	v := h.Value().(ExternalValueProvider)

	C.nix_set_string_return(nil, C.CString(v.TypeOf()))
}

//export externalValueCoerceToString
func externalValueCoerceToString(self unsafe.Pointer, c *C.nix_string_context, coerceMore, copyToStore C.int, res *C.nix_string_return) {
	h := cgo.Handle(self)
	v := h.Value().(ExternalValueProvider)

	addCtx := func(s string) error {
		cerr := C.nix_external_add_string_context(nil, c, C.CString(s))
		return nixError(cerr, nil)
	}
	ret, err := v.CoerceToString(addCtx, coerceMore != 0, copyToStore != 0)
	if err != nil {
		return
	}
	C.nix_set_string_return(res, C.CString(ret))
}

//export externalValueEqual
func externalValueEqual(self unsafe.Pointer, other unsafe.Pointer) C.int {
	h := cgo.Handle(self)
	v := h.Value().(ExternalValueProvider)

	oh := cgo.Handle(other)
	ov := oh.Value().(ExternalValueProvider)

	if v.Equal(ov) {
		return 1
	}
	return 0
}

//export externalValuePrintValueAsJSON
func externalValuePrintValueAsJSON(self unsafe.Pointer, state *C.EvalState, strict int, c *C.nix_string_context, copyToStore bool, res *C.nix_string_return) {
	h := cgo.Handle(self)
	v := h.Value().(ExternalValueProvider)

	C.nix_set_string_return(res, C.CString(v.PrintValueAsJSON(strict != 0, copyToStore)))
}

//export finalizeExternalValue
func finalizeExternalValue(obj, cd unsafe.Pointer) {
	h := cgo.Handle(cd)
	h.Delete()
}
