package gonix

//
// #cgo pkg-config: nix-expr-c
// #include <nix_api_util.h>
// #include <nix_api_expr.h>
// #include <nix_api_value.h>
// #include <string.h>
/*
typedef const char cchar_t;
void nixGetCallbackString_cgo(cchar_t * start, unsigned int n, char * user_data) {
   void nixGetCallbackString (const char * start, unsigned int n, char * user_data);
  nixGetCallbackString(start, n, user_data);
}
*/
import "C"
