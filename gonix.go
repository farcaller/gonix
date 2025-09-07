// Package gonix provides bindings to the nix APIs. It can work with any store
// supported by nix, including working with `/nix/store` directly or over the
// nix-daemon.
package gonix

// #cgo pkg-config: nix-expr-c nix-util-c nix-store-c nix-main-c
// #include <stdlib.h>
// #include <stdbool.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
// #include <nix_api_main.h>
/*
typedef const char cchar_t;
void nixGetCallbackString_cgo(cchar_t * start, unsigned int n, char ** user_data);
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

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
//
// Warning: since gonix started binding against Nix 2.30+, this always returns NIX_ERR_KEY
func GetSetting(ctx *Context, name string) (string, error) {
	var str *string = new(string)
	strh := cgo.NewHandle(str)
	defer strh.Delete()

	nameCStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameCStr))

	cerr := C.nix_setting_get(ctx.ccontext, nameCStr, (*[0]byte)(C.nixGetCallbackString_cgo), nil)
	if cerr != C.NIX_OK {
		return "", nixError(cerr, ctx)
	}
	return *str, nil
}

// SetSetting sets the setting to the passed cvalue. this value affects all the
// calls done to nix API within the lifetime of the executable (regardless of
// the context passed).
//
// Warning: since gonix started binding against Nix 2.30+, this always returns NIX_ERR_KEY
func SetSetting(ctx *Context, name, value string) error {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	valueStr := C.CString(value)
	defer C.free(unsafe.Pointer(valueStr))

	cerr := C.nix_setting_set(ctx.ccontext, nameStr, valueStr)
	return nixError(cerr, ctx)
}

const maxBufferSize = 1024 * 10

// type  nixStringCallback func(start *char[0] , n int, data *void)
