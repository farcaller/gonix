package gonix_test

import (
	"fmt"

	"github.com/farcaller/gonix"
)

func ExampleRegisterGlobalPrimOp() {
	ctx := gonix.NewContext()
	store, _ := gonix.NewStore(ctx, "dummy", nil)

	gonix.RegisterGlobalPrimOp(
		ctx,
		"summer",
		[]string{"arg0", "arg1", "arg2"},
		"docs",
		func(ctx *gonix.Context, state *gonix.State, args ...*gonix.Value) *gonix.Value {
			var sum int64 = 0
			for _, val := range args {
				i, _ := val.GetInt()
				sum += i
			}
			r, _ := state.NewInt(sum)
			return r
		})

	state := store.NewState(nil)

	res, _ := state.EvalExpr(`builtins.summer 1 2 3`, ".")

	fmt.Println(res)
	// Output: 6
}

func ExampleState_NewPrimOp() {
	ctx := gonix.NewContext()
	store, _ := gonix.NewStore(ctx, "dummy", nil)
	state := store.NewState(nil)

	op, _ := gonix.NewPrimOp(
		ctx,
		"helloworlder",
		[]string{"target"},
		"docs",
		func(ctx *gonix.Context, state *gonix.State, args ...*gonix.Value) *gonix.Value {
			v, _ := state.NewString(fmt.Sprintf("hello, %s", args[0]))
			return v
		})
	vop, _ := state.NewPrimOp(op)

	world, _ := state.NewString("world")

	res, err := state.Call(vop, world)
	if err != nil {
		panic(fmt.Errorf("failed to call: %v", err))
	}

	fmt.Println(res)
	// Output: hello, world
}
