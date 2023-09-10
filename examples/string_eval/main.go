package main

import (
	"fmt"

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
	state := gonix.NewState(store, nil)

	fmt.Printf("running nix v%s on %s\n",
		mustEvalString(state, "builtins.nixVersion"),
		mustEvalString(state, "builtins.currentSystem"))
}
