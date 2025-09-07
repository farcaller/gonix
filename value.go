package gonix

// #cgo pkg-config: nix-expr-c nix-main-c
// #include <stdlib.h>
// #include <nix_api_util.h>
// #include <nix_api_value.h>
// #include <nix_api_expr.h>
// #include <nix_api_main.h>
/*
typedef const char cchar_t;
void nixGetCallbackString_cgo(cchar_t * start, unsigned int n, char ** user_data);
*/
import "C"

import (
	"fmt"
	"runtime"
	"runtime/cgo"
	"unsafe"
)

// ValueType is the nix type contained in the [Value].
type ValueType int

const (
	NixTypeThunk ValueType = iota
	NixTypeInt
	NixTypeFloat
	NixTypeBool
	NixTypeString
	NixTypePath
	NixTypeNull
	NixTypeAttrs
	NixTypeList
	NixTypeFunction
	NixTypeExternal
)

// Value is a wrapper around a nix value.
type Value struct {
	state  *State
	cvalue *C.nix_value
}

func wrapValue(state *State, cvalue *C.nix_value) (*Value, error) {
	runtime.SetFinalizer(&cvalue, func(v **C.nix_value) {
		C.nix_gc_decref(state.context().ccontext, unsafe.Pointer(*v))
	})
	return &Value{state, cvalue}, nil
}

// newValue creates a new value.
//
// It seems it's an error to create a value and never write anything into it,
// so this API is private. Use the [State] to create a value from some value.
func newValue(state *State) (*Value, error) {
	cvalue := C.nix_alloc_value(state.context().ccontext, state.cstate)
	if cvalue == nil {
		return nil, state.context().lastError()
	}
	return wrapValue(state, cvalue)
}

func (v *Value) context() *Context {
	return v.state.ctx
}

// Force forces the evaluation of a nix value.
func (v *Value) Force() error {
	cerr := C.nix_value_force(v.context().ccontext, v.state.cstate, v.cvalue)
	return nixError(cerr, v.context())
}

// ForceDeep calls [Value.Force] recursively.
func (v *Value) ForceDeep() error {
	cerr := C.nix_value_force_deep(v.context().ccontext, v.state.cstate, v.cvalue)
	return nixError(cerr, v.context())
}

// Type returns the nix type stored in this value.
func (v *Value) Type() ValueType {
	t := C.nix_get_type(v.context().ccontext, v.cvalue)

	return (ValueType)(t)
}

func (v *Value) String() string {
	var val, err any
	switch v.Type() {
	case NixTypeThunk:
		// TODO: I can totally see this not ending well.
		v.ForceDeep()
		return v.String()
	case NixTypeInt:
		val, err = v.GetInt()

	case NixTypeFloat:
		val, err = v.GetFloat()
	case NixTypeBool:
		val, err = v.GetBool()
	case NixTypeString:
		val, err = v.GetString()
	case NixTypePath:
		val, err = v.GetPath()
	case NixTypeNull:
		val, err = nil, nil
	case NixTypeAttrs:
		val, err = v.GetAttrs()
	case NixTypeList:
		val, err = v.GetList()
	case NixTypeFunction:
		return "<function>"
	case NixTypeExternal:
		val, err = v.GetExternalValue()
	}
	if err != nil {
		return fmt.Sprintf("<error: %v", err)
	}
	return fmt.Sprint(val)
}

type stringCallback func(string)

//export nixGetCallbackString
func nixGetCallbackString(str *C.cchar_t, n C.int, userdata unsafe.Pointer) {
	// FIXME: Probably need some error handling here.
	h := cgo.Handle(userdata)
	v := h.Value().(*string)
	res := C.GoStringN(str, n)
	*v = res
}

// GetString returns the string value iff the value contains a string.
func (v *Value) GetString() (string, error) {
	cctx := v.context().ccontext

	typ := v.Type()
	if typ != C.NIX_TYPE_STRING {
		return "", fmt.Errorf("expected a string, got %v", typ)
	}

	var str *string = new(string)
	strh := cgo.NewHandle(str)
	defer strh.Delete()
	cerr := C.nix_get_string(cctx, v.cvalue, (*[0]byte)(C.nixGetCallbackString_cgo), unsafe.Pointer(strh))
	if cerr != C.NIX_OK {
		return "", fmt.Errorf("failed to get the string value: %v", v.context().lastError())
	}
	return *str, nil
}

// GetInt returns the int value iff the value contains an int.
func (v *Value) GetInt() (int64, error) {
	cctx := v.context().ccontext

	typ := v.Type()
	if typ != C.NIX_TYPE_INT {
		return 0, fmt.Errorf("expected an int, got %v", typ)
	}

	i := C.nix_get_int(cctx, v.cvalue)
	if err := v.context().lastError(); err != nil {
		return 0, err
	}
	return int64(i), nil
}

// GetFloat returns the float value iff the value contains a float.
func (v *Value) GetFloat() (float64, error) {
	cctx := v.context().ccontext

	typ := v.Type()
	if typ != C.NIX_TYPE_FLOAT {
		return 0, fmt.Errorf("expected a float, got %v", typ)
	}

	f := C.nix_get_float(cctx, v.cvalue)
	if err := v.context().lastError(); err != nil {
		return 0, err
	}
	return float64(f), nil
}

// GetBool returns the bool value iff the value contains a bool.
func (v *Value) GetBool() (bool, error) {
	cctx := v.context().ccontext

	typ := v.Type()
	if typ != C.NIX_TYPE_BOOL {
		return false, fmt.Errorf("expected a bool, got %v", typ)
	}

	b := C.nix_get_bool(cctx, v.cvalue)
	if err := v.context().lastError(); err != nil {
		return false, err
	}
	return bool(b), nil
}

// GetPath returns the path value iff the value contains a path.
func (v *Value) GetPath() (string, error) {
	cctx := v.context().ccontext

	typ := v.Type()
	if typ != C.NIX_TYPE_PATH {
		return "", fmt.Errorf("expected a string, got %v", typ)
	}

	s := C.nix_get_path_string(cctx, v.cvalue)
	if s == nil {
		return "", fmt.Errorf("fialed to get the path string value: %v", v.context().lastError())
	}
	str := C.GoString(s)
	return str, nil
}

// GetExternalValue returns the external value iff the value contains an external value.
func (v *Value) GetExternalValue() (*ExternalValue, error) {
	typ := v.Type()
	if typ != C.NIX_TYPE_EXTERNAL {
		return nil, fmt.Errorf("expected an external value, got %v", typ)
	}
	cev := C.nix_get_external(v.context().ccontext, v.cvalue)
	if err := v.context().lastError(); err != nil {
		return nil, err
	}

	return &ExternalValue{cev}, nil
}

// GetList returns the list value iff the value contains a list.
func (v *Value) GetList() ([]*Value, error) {
	typ := v.Type()
	if typ != C.NIX_TYPE_LIST {
		return nil, fmt.Errorf("expected a list, got %v", typ)
	}

	listLen := C.nix_get_list_size(v.context().ccontext, v.cvalue)
	if listLen == 0 {
		return nil, nil
	}
	list := make([]*Value, listLen)
	for idx := 0; idx < int(listLen); idx++ {
		cval := C.nix_get_list_byidx(v.context().ccontext, v.cvalue, v.state.cstate, C.uint(idx))
		if cval == nil {
			return nil, v.context().lastError()
		}
		val, err := wrapValue(v.state, cval)
		if err != nil {
			return nil, err
		}
		list[idx] = val
	}
	return list, nil
}

// GetAttrs returns the attrset value iff the value contains an attrset.
func (v *Value) GetAttrs() (map[string]*Value, error) {
	typ := v.Type()
	if typ != C.NIX_TYPE_ATTRS {
		return nil, fmt.Errorf("expected an attrset, got %v", typ)
	}

	attrs := make(map[string]*Value)

	attrsLen := C.nix_get_attrs_size(v.context().ccontext, v.cvalue)

	if attrsLen == 0 {
		return attrs, nil
	}
	for idx := 0; idx < int(attrsLen); idx++ {
		var namePtr *C.char
		cval := C.nix_get_attr_byidx(v.context().ccontext, v.cvalue, v.state.cstate, C.uint(idx), (**C.char)(&namePtr))
		if cval == nil {
			return nil, v.context().lastError()
		}
		val, err := wrapValue(v.state, cval)
		if err != nil {
			return nil, err
		}
		attrs[C.GoString(namePtr)] = val
	}
	return attrs, nil
}

// SetBool sets the value to the passed bool.
func (v *Value) SetBool(b bool) error {
	cerr := C.nix_init_bool(v.context().ccontext, v.cvalue, C.bool(b))
	return nixError(cerr, v.context())
}

// SetString sets the value to the passed string.
func (v *Value) SetString(s string) error {
	cerr := C.nix_init_string(v.context().ccontext, v.cvalue, C.CString(s))
	return nixError(cerr, v.context())
}

// SetPath sets the value to the passed string as a path.
func (v *Value) SetPath(ps string) error {
	cerr := C.nix_init_path_string(v.context().ccontext, v.state.cstate, v.cvalue, C.CString(ps))
	return nixError(cerr, v.context())
}

// SetFloat sets the value to the passed float.
func (v *Value) SetFloat(f float64) error {
	cerr := C.nix_init_float(v.context().ccontext, v.cvalue, C.double(f))
	return nixError(cerr, v.context())
}

// SetNull sets the value to null.
func (v *Value) SetNull() error {
	cerr := C.nix_init_null(v.context().ccontext, v.cvalue)
	return nixError(cerr, v.context())
}

// SetExternalValue sets the value to the passed external value.
func (v *Value) SetExternalValue(ev *ExternalValue) error {
	cerr := C.nix_init_external(v.context().ccontext, v.cvalue, ev.cev)
	err := nixError(cerr, v.context())
	if err != nil {
		return err
	}
	return nil
}

// SetList sets the value to the passed list.
func (v *Value) SetList(items []*Value) error {
	clb := C.nix_make_list_builder(v.context().ccontext, v.state.cstate, (C.ulong)(len(items)))
	if clb == nil {
		return v.context().lastError()
	}
	for idx, item := range items {
		cerr := C.nix_list_builder_insert(v.context().ccontext, clb, (C.uint)(idx), item.cvalue)
		err := nixError(cerr, v.context())
		if err != nil {
			return err
		}
	}
	C.nix_make_list(v.context().ccontext, clb, v.cvalue)
	C.nix_list_builder_free(clb)
	return nil
}

// SetAttrs sets the value to the passed attrset.
func (v *Value) SetAttrs(attrs map[string]*Value) error {
	cbb := C.nix_make_bindings_builder(v.context().ccontext, v.state.cstate, C.ulong(len(attrs)))
	if cbb == nil {
		return v.context().lastError()
	}
	defer C.nix_bindings_builder_free(cbb)

	for k, v := range attrs {
		cerr := C.nix_bindings_builder_insert(v.context().ccontext, cbb, C.CString(k), v.cvalue)
		err := nixError(cerr, v.context())
		if err != nil {
			return err
		}
	}
	cerr := C.nix_make_attrs(v.context().ccontext, v.cvalue, cbb)
	err := nixError(cerr, v.context())
	if err != nil {
		return err
	}
	return nil
}

// SetPrimOp sets the value to the passed PrimOp.
func (v *Value) SetPrimOp(op *PrimOp) error {
	cerr := C.nix_init_primop(v.context().ccontext, v.cvalue, op.cprimop)
	return nixError(cerr, v.context())
}
