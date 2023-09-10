package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/farcaller/gonix/pkg/gonix"
)

func mustEvalString(state *gonix.State, expr string) string {
	val, err := state.EvalExpr(expr, ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}
	s, err := val.GetString()
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}
	return *s
}

func main() {
	ctx := gonix.NewContext()
	store := gonix.NewStore(ctx, "", nil)
	state := store.NewState(strings.Split(os.Getenv("NIX_PATH"), ":"))

	path := mustEvalString(state, "(import <nixpkgs> {}).hello.drvPath")
	sp, err := store.ParsePath(path)
	if err != nil {
		panic(fmt.Errorf("failed to parse path: %v", err))
	}
	fmt.Printf("drv %s vaild? %v\n", path, sp.Valid())
	outs, err := sp.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build: %v", err))
	}
	for k, v := range outs {
		fmt.Printf("%s -> %s\n", k, v)
	}
}
