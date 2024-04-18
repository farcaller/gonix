// Package gonix provides bindings to the nix APIs. It can work with any store
// supported by nix, including working with `/nix/store` directly or over the
// nix-daemon.
package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
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
func GetSetting(ctx *Context, name string) (string, error) {
	var str *string = new(string)
	strh := cgo.NewHandle(str)
	defer strh.Delete()
	cerr := C.nix_setting_get(ctx.ccontext, C.CString(name), (*[0]byte)(C.nixGetCallbackString_cgo), unsafe.Pointer(strh))
	if cerr != C.NIX_OK {
		return "", nixError(cerr, ctx)
	}
	return *str, nil
}

// SetSetting sets the setting to the passed value. This value affects all the
// calls done to nix API within the lifetime of the executable (irregardless of
// the context passed).
func SetSetting(ctx *Context, name, value string) error {
	cerr := C.nix_setting_set(ctx.ccontext, C.CString(name), C.CString(value))
	return nixError(cerr, ctx)
}

const maxBufferSize = 1024 * 10

// type  nixStringCallback func(start *char[0] , n int, data *void)

// InitPlugins loads the plugins specified in Nix's plugin-files setting.
//
// This function should be called once (if needed), after all the required settings were set.
func InitPlugins(ctx *Context) error {
	cerr := C.nix_init_plugins(ctx.ccontext)
	return nixError(cerr, ctx)
}
