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
	ev := h.Value().(externalValue)
	ev.print(pr)
	// TODO: ev.print can fail tho?
}

//export externalValueCoerceToString
func externalValueCoerceToString(self unsafe.Pointer, c *C.nix_string_context, coerceMore, copyToStore C.int, res *C.nix_string_return) {
	h := cgo.Handle(self)
	ev := h.Value().(externalValue)
	ev.coerceToString(c, coerceMore, copyToStore, res)
	// TODO: ev.coerceToString can fail tho?
}
