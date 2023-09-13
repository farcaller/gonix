package gonix_test

import (
	"fmt"

	"github.com/farcaller/gonix"
)

func ExampleState_Call() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	f, err := state.EvalExpr("value: value * 2", ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval a function: %v", err))
	}
	v, err := state.NewInt(21)
	if err != nil {
		panic(fmt.Errorf("failed to create an int: %v", err))
	}
	res, err := state.Call(f, v)
	if err != nil {
		panic(fmt.Errorf("failed to call a function: %v", err))
	}
	resInt, err := res.GetInt()
	if err != nil {
		panic(fmt.Errorf("failed to convert result to int: %v", err))
	}
	fmt.Println(resInt)
	// Output: 42
}

func ExampleState_EvalExpr() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	val, err := state.EvalExpr("101 - 59", ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}
	i, err := val.GetInt()
	if err != nil {
		panic(fmt.Errorf("failed to convert result to int: %v", err))
	}
	fmt.Println(i)
	// Output: 42
}
