package gonix

// #cgo pkg-config: nix-expr-c
// #include <stdlib.h>
// #include <stdio.h>
/* void nixStoreBuildCallback_cgo(void * userdata, const char * outname, const char * out)
{
	void nixStoreBuildCallback(void *, const char *, const char *);
	nixStoreBuildCallback(userdata, outname, out);
}
*/
import "C"
