package gonix_test

import (
	"fmt"

	"github.com/farcaller/gonix"
)

type simpleExternal struct{}

func (v *simpleExternal) Print() string {
	return "printed"
}

func (v *simpleExternal) CoerceToString(addContext func(string) error, copyMore, copyToStore bool) (string, error) {
	return fmt.Sprintf("coerced with copyMore=%v", copyMore), nil
}

func ExampleState_NewExternalValue() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	externalVal, err := state.NewExternalValue(&simpleExternal{})
	if err != nil {
		panic(fmt.Errorf("failed to create an external value: %v", err))
	}
	toStr, err := state.EvalExpr("builtins.toString", ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}
	res, err := state.Call(toStr, externalVal)
	if err != nil {
		panic(fmt.Errorf("failed to call a function: %v", err))
	}
	resStr, err := res.GetString()
	if err != nil {
		panic(fmt.Errorf("failed to convert result to int: %v", err))
	}
	fmt.Println(*resStr)
	// Output: coerced with copyMore=true
}
