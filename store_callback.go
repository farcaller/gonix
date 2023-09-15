package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <nix_api_store.h>
/*
typedef const char cchar_t;
void nixStoreBuildCallback_cgo(void * userdata, cchar_t * outname, cchar_t * out) {
	void nixStoreBuildCallback(void *, cchar_t *, cchar_t *);
	nixStoreBuildCallback(userdata, outname, out);
}
*/
import "C"
