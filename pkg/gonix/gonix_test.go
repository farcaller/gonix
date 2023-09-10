package gonix_test

import (
	"fmt"

	"github.com/farcaller/gonix/pkg/gonix"
)

func Example() {
	ctx := gonix.NewContext()
	store := gonix.NewStore(ctx, "", nil)
	state := store.NewState(nil)

	val, err := state.EvalExpr("builtins.toJSON { answer = 42; }", ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}

	strVal, err := val.GetString()
	if err != nil {
		panic(fmt.Errorf("failed to convert the value to string: %v", err))
	}

	fmt.Println(*strVal)
	// Output: {"answer":42}
}
