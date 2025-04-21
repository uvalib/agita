// util/struct.go
//
// Functions supporting struct manipulation.

package util

import (
	"fmt"
	"reflect"
)

// ============================================================================
// Exported functions
// ============================================================================

// Return a copy of a struct instance.
func CloneStruct[T any](arg T) T {
    src := structValue(arg)
    typ := src.Type()
    max := src.NumField()
    res := reflect.New(typ).Elem()
    for i := range max {
        name := typ.Field(i).Name
        if field := res.FieldByName(name); field.CanSet() {
            field.Set(src.Field(i))
        }
    }
    return res.Interface().(T)
}

// Return the exported field names of a struct.
func StructFields(arg any) []string {
    src := structValue(arg)
    typ := src.Type()
    max := src.NumField()
    res := make([]string, 0, max)
    for i := range max {
        if field := src.Field(i); field.CanInterface() {
            key := typ.Field(i).Name
            res = append(res, key)
        }
    }
    return res
}

// Transform a struct into a map of field names and values.
func StructMap(arg any) map[string]any {
    src := structValue(arg)
    typ := src.Type()
    max := src.NumField()
    res := make(map[string]any, max)
    for i := range max {
        if field := src.Field(i); field.CanInterface() {
            key := typ.Field(i).Name
            val := field.Interface()
            res[key] = val
        }
    }
    return res
}

// ============================================================================
// Internal functions
// ============================================================================

// Turns a struct or pointer into a reflect.Value for the struct object.
func structValue(arg any) reflect.Value {
    val := reflect.ValueOf(arg)
    ptr := (val.Kind() == reflect.Ptr)
    if ptr {
        if val.IsNil() { panic("nil pointer") }
        val = val.Elem()
    }
    if kind := val.Kind(); kind != reflect.Struct {
        argType := kind.String()
        if ptr { argType = "pointer to " + argType }
        panic(fmt.Errorf("invalid %s arg = %v", argType, arg))
    }
    return val
}
