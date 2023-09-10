package main

import (
	"fmt"

	"github.com/farcaller/gonix/pkg/gonix"
)

func mustEvalString(state *gonix.State, expr string) string {
	val, err := gonix.ExprEvalFromString(state, expr, ".")
	if err != nil {
		panic(fmt.Errorf("faled to eval: %v", err))
	}
	return val.String()
}

func main() {
	ctx := gonix.NewContext()
	store := gonix.NewStore(ctx, "", nil)
	state := gonix.NewState(store, nil)

	fmt.Printf("running nix v%s on %s\n",
		mustEvalString(state, "builtins.nixVersion"),
		mustEvalString(state, "builtins.currentSystem"))
}
