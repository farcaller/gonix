package gonix

// #cgo pkg-config: nix-expr-c nix-main-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
// #include <nix_api_main.h>
import "C"
import (
	"errors"
	"fmt"
)

func nixError(err C.nix_err, ctx *Context) error {
	switch err {
	case C.NIX_OK:
		return nil
	case C.NIX_ERR_UNKNOWN:
		return errors.New("unknown nix error")
	case C.NIX_ERR_OVERFLOW:
		return errors.New("nix overflow")
	case C.NIX_ERR_KEY:
		return errors.New("nix key error")
	case C.NIX_ERR_NIX_ERROR:
		return fmt.Errorf("nix error: %s", C.GoString(C.nix_err_msg(nil, ctx.ccontext, nil)))
	default:
		return fmt.Errorf("unknown nix error code %d", err)
	}
}
