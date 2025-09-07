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

	f, err := state.EvalExpr("cvalue: cvalue * 2", ".")
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

	fmt.Println(res)
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

	fmt.Println(val)
	// Output: 42
}

func ExampleState_NewString() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	s, err := state.NewString("hello")
	if err != nil {
		panic(fmt.Errorf("failed to create a string: %v", err))
	}

	fmt.Println(s)
	// Output:
	// hello
}

func ExampleState_NewList() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	newString := func(str string) *gonix.Value {
		s, err := state.NewString(str)
		if err != nil {
			panic(fmt.Errorf("failed to create a string: %v", err))
		}
		return s
	}

	newInt := func(i int64) *gonix.Value {
		s, err := state.NewInt(i)
		if err != nil {
			panic(fmt.Errorf("failed to create a number: %v", err))
		}
		return s
	}

	l, err := state.NewList([]*gonix.Value{
		newString("hello"),
		newString("world"),
		newInt(11),
	})
	if err != nil {
		panic(fmt.Errorf("failed to get the string cvalue: %v", err))
	}

	fmt.Println(l)
	// Output:
	// [hello world 11]
}

func ExampleState_NewInt() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	v, err := state.NewInt(15)
	if err != nil {
		panic(fmt.Errorf("failed to create an int: %v", err))
	}

	fmt.Println(v)
	// Output: 15
}

func ExampleState_NewFloat() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	v, err := state.NewFloat(3.50)
	if err != nil {
		panic(fmt.Errorf("failed to create a float: %v", err))
	}

	fmt.Println(v)
	// Output: 3.5
}

func ExampleState_NewBool() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	v, err := state.NewBool(true)
	if err != nil {
		panic(fmt.Errorf("failed to create a float: %v", err))
	}

	fmt.Println(v)
	// Output: true
}

func ExampleState_NewPath() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	s, err := state.NewPath("/foo/bar")
	if err != nil {
		panic(fmt.Errorf("failed to create a path: %v", err))
	}

	fmt.Println(s)
	// Output:
	// /foo/bar
}

func ExampleState_NewNull() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	s, err := state.NewNull()
	if err != nil {
		panic(fmt.Errorf("failed to create a null: %v", err))
	}

	fmt.Println(s)
	// Output:
	// <nil>
}

func ExampleState_NewAttrs() {
	ctx := gonix.NewContext()
	store, err := gonix.NewStore(ctx, "dummy", nil)
	if err != nil {
		panic(fmt.Errorf("failed to create a store: %v", err))
	}
	state := store.NewState(nil)

	newString := func(str string) *gonix.Value {
		s, err := state.NewString(str)
		if err != nil {
			panic(fmt.Errorf("failed to create a string: %v", err))
		}
		return s
	}

	newInt := func(i int64) *gonix.Value {
		s, err := state.NewInt(i)
		if err != nil {
			panic(fmt.Errorf("failed to create a number: %v", err))
		}
		return s
	}

	as, err := state.NewAttrs(map[string]*gonix.Value{
		"hello":  newString("world"),
		"number": newInt(11),
	})
	if err != nil {
		panic(fmt.Errorf("failed to build an attrset: %v", err))
	}

	fmt.Println(as)
	// Output:
	// map[hello:world number:11]
}
