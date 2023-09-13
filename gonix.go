// Package gonix provides bindings to the nix APIs. It can work with any store
// supported by nix, including working with `/nix/store` directly or over the
// nix-daemon.
package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
import "C"
import "unsafe"

func init() {
	C.nix_libutil_init(nil)
	C.nix_libexpr_init(nil)
	C.nix_libstore_init(nil)
}

// Version returns the API version.
//
// Example:
//
//	v := Version()
//
// Output:
//
//	2.18.0pre20230828_af7d89a
func Version() string {
	return C.GoString(C.nix_version_get())
}

// GetSetting returns the value of a setting.
func GetSetting(ctx *Context, name string) (string, error) {
	return reallocatingBufferRead(ctx, func(ctx *Context, buf *C.char, bl C.int) C.int {
		return C.nix_setting_get(ctx.ccontext, C.CString(name), buf, bl)
	})
}

// SetSetting sets the setting to the passed value. This value affects all the
// calls done to nix API within the lifetime of the executable (irregardless of
// the context passed).
func SetSetting(ctx *Context, name, value string) error {
	cerr := C.nix_setting_set(ctx.ccontext, C.CString(name), C.CString(value))
	return nixError(cerr, ctx)
}

const maxBufferSize = 1024 * 10

func reallocatingBufferRead(ctx *Context, call func(*Context, *C.char, C.int) C.int) (string, error) {
	sz := 1
	for {
		currentSize := 1024 * sz
		if currentSize > maxBufferSize {
			return "", nixError(C.NIX_ERR_OVERFLOW, nil)
		}
		buf := make([]byte, currentSize)
		cerr := call(ctx, (*C.char)(unsafe.Pointer(&buf[0])), C.int(len(buf)))
		if cerr == C.NIX_OK {
			return C.GoString((*C.char)(unsafe.Pointer(&buf[0]))), nil
		}
		if cerr == C.NIX_ERR_OVERFLOW {
			sz += 1
			continue
		}
		return "", nixError(cerr, ctx)
	}
}
