// util/type.go

package util

import (
	"reflect"
	"slices"
	"strconv"

	"golang.org/x/exp/constraints"
)

// ============================================================================
// Internal variables
// ============================================================================

var aryKinds = []reflect.Kind{
    reflect.Array,
    reflect.Slice,
}

var intKinds = []reflect.Kind{
    reflect.Int,
    reflect.Int8,
    reflect.Int16,
    reflect.Int32,
    reflect.Int64,
    reflect.Uint,
    reflect.Uint8,
    reflect.Uint16,
    reflect.Uint32,
    reflect.Uint64,
}

var ptrKinds = []reflect.Kind{
    reflect.Ptr,
    reflect.Interface,
    reflect.Slice,
    reflect.Map,
    reflect.Chan,
    reflect.Func,
}

// ============================================================================
// Exported functions
// ============================================================================

// Indicate whether a value is a map.
func IsMap(arg any) bool {
    return Kind(arg) == reflect.Map
}

// Indicate whether a value behaves like an array.
func IsArray(arg any) bool {
    return slices.Contains(aryKinds, Kind(arg))
}

// Indicate whether a value is a kind of integer.
func IsInteger(arg any) bool {
    return slices.Contains(intKinds, Kind(arg))
}

// Indicate whether a value is a nilable pointer.
func IsPtr(arg any) bool {
    return slices.Contains(ptrKinds, Kind(arg))
}

// Indicate whether the provided value is nil.
//  NOTE: returns false if `value` does not have a nilable type.
func IsNil(arg any) bool {
    val := Value(arg)
    return IsPtr(val) && val.IsNil()
}

// Indicate whether the item is nil or has no elements.
//  NOTE: returns true if `value` is nil.
//  NOTE: returns false if `value` is neither a slice nor a map.
func IsEmpty(arg any) bool {
    val := Value(arg)
    switch {
        case IsMap(val):    return len(val.MapKeys()) == 0
        case IsArray(val):  return val.Len() == 0
        default:            return IsNil(val)
    }
}

// ============================================================================
// Exported functions
// ============================================================================

// Render a number in the indicated base (2 to 36, inclusive).
func Base[T ~string|constraints.Integer](arg T, base int) string {
    if val := Value(arg); val.Kind() != reflect.String {
        return strconv.FormatInt(val.Int(), base)
    } else if v, e := strconv.Atoi(val.String()); e == nil {
        return strconv.FormatInt(int64(v), base)
    } else {
        return "0"
    }
}

// ============================================================================
// Exported functions
// ============================================================================

// Returns the reflect.Kind for the original item.
func Kind(arg any) reflect.Kind {
    switch v := arg.(type) {
        case reflect.Type:  return v.Kind()
        case reflect.Value: return v.Kind()
        default:            return reflect.TypeOf(arg).Kind()
    }
}

// Returns the reflect.Value for the original item.
func Value(arg any) reflect.Value {
    switch v := arg.(type) {
        case reflect.Type:  panic("requires a value")
        case reflect.Value: return v
        default:            return reflect.ValueOf(arg)
    }
}
