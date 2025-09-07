package gonix_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/farcaller/gonix"
)

func ExampleRegisterGlobalPrimOp() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}

	err = gonix.RegisterGlobalPrimOp(
		ctx,
		"summer",
		[]string{"arg0", "arg1", "arg2"},
		"docs",
		func(ctx *gonix.Context, state *gonix.State, args []*gonix.Value, ret *gonix.Value) error {
			var sum int64 = 0
			for _, val := range args {
				i, _ := val.GetInt()
				sum += i
			}
			return ret.SetInt(sum)
		})
	if err != nil {
		panic(fmt.Errorf("failed to register: %v", err))
	}

	state := store.NewState([]string{os.Getenv("NIX_PATH")})

	res, err := state.EvalExpr(`builtins.summer 1 2 3`, ".")
	if err != nil {
		panic(fmt.Errorf("failed to eval: %v", err))
	}

	i, err := res.GetInt()
	if err != nil {
		panic(fmt.Errorf("failed to convert the value to int: %v", err))
	}
	if i != 6 {
		panic(fmt.Errorf("expected 6, got %d", i))
	}
}

func TestRegisterGlobalPrimOpWithError(t *testing.T) {
	ctx := gonix.NewContext()
	store, _ := gonix.NewStore(ctx, "dummy", nil)

	err := gonix.RegisterGlobalPrimOp(
		ctx,
		"goError",
		[]string{"arg0", "arg1", "arg2"},
		"docs",
		func(ctx *gonix.Context, state *gonix.State, args []*gonix.Value, ret *gonix.Value) error {
			return fmt.Errorf("oops! An error - args were: %v", args)
		})
	if err != nil {
		t.Fatalf("failed to register: %v", err)
	}

	state := store.NewState(nil)

	res, err := state.EvalExpr(`builtins.goError 1 2 3`, ".")
	if err == nil {
		t.Fatalf("expected an error, got %v", res)
	}

	if !strings.Contains(err.Error(), "oops! An error - args were:") {
		t.Fatalf("expected an error with the args, got %v", err)
	}
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
		func(ctx *gonix.Context, state *gonix.State, args []*gonix.Value, ret *gonix.Value) error {
			err := ret.SetString(fmt.Sprintf("hello, %s", args[0]))
			return err
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
