package gonix_test

import (
	"fmt"

	"github.com/farcaller/gonix"
)

func Example() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	val, err := state.EvalExpr("builtins.toJSON { answer = 42; }", ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}

	strVal, err := val.GetString()
	if err != nil {
		panic(fmt.Errorf("failed to convert the value to string: %v", err))
	}

	fmt.Println(strVal)
	// Output: {"answer":42}
}
func ExampleGetSetting() {
	ctx := gonix.NewContext()
	val, err := gonix.GetSetting(ctx, "trace-verbose")
	if err != nil {
		panic(fmt.Errorf("failed to read the setting: %v", err))
	}

	fmt.Println(val)
	// Output: false
}
