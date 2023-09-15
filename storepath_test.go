package gonix_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/farcaller/gonix"
)

func ExampleStorePath_Build() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(strings.Split(os.Getenv("NIX_PATH"), ":"))

	val, err := state.EvalExpr("(import <nixpkgs> {}).hello.drvPath", ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}

	sp, err := store.ParsePath(val.String())
	if err != nil {
		panic(fmt.Errorf("failed to parse path: %v", err))
	}
	if !sp.Valid() {
		panic(fmt.Errorf("store path isn't valid"))
	}

	outs, err := sp.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build: %v", err))
	}

	for k, _ := range outs {
		fmt.Printf("%s\n", k)
	}
	// Output:
	// out
}
