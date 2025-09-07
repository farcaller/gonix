package gonix_test

import (
	"fmt"

	"github.com/farcaller/gonix"
)

type simpleExternal struct{}

func (v *simpleExternal) Print() string {
	return "printed"
}

func (v *simpleExternal) ShowType() string {
	return "simpleexternal"
}
func (v *simpleExternal) TypeOf() string {
	return "simpleexternal"
}

func (v *simpleExternal) CoerceToString(addContext func(string) error, copyMore, copyToStore bool) (string, error) {
	return fmt.Sprintf("coerced with copyMore=%v", copyMore), nil
}

func (v *simpleExternal) Equal(other gonix.ExternalValueProvider) bool {
	return false
}

func (v *simpleExternal) PrintValueAsJSON(strict, copyToStore bool) string {
	return `{ "cvalue" : "simpleexternal" }`
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

	strInterpolation, err := state.EvalExpr(`v: "${v}"`, ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}
	res2, err := state.Call(strInterpolation, externalVal)
	if err != nil {
		panic(fmt.Errorf("failed to call a function: %v", err))
	}

	fmt.Println(res)
	fmt.Println(res2)
	// Output:
	// coerced with copyMore=true
	// coerced with copyMore=false
}

func ExampleExternalValueProvider_PrintValueAsJSON() {
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
	toJSON, err := state.EvalExpr("builtins.toJSON", ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}
	res, err := state.Call(toJSON, externalVal)
	if err != nil {
		panic(fmt.Errorf("failed to call a function: %v", err))
	}

	fmt.Println(res)
	// Output: {"cvalue":"simpleexternal"}
}
