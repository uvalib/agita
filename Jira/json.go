// Jira/json.go
//
// Addresses apparent bugs in encoding/json.
//
// [1] According to the documentation, the field tag `json:"-"` should always
// omit the Field from the result but it does not.  The only way to avoid
// seeing them in the JSON result is to avoid having the field in the struct.
//
// [2] Not a bug per se, but "omitempty" doesn't omit empty structs.  Because
// the Jira-defined objects use structs rather than pointers to structs, an
// empty literal structure always shows up in the results.  Redefining the
// matching field as a pointer allows this to be avoided.

package Jira

import (
	"fmt"
)

// ============================================================================
// Internal functions
// ============================================================================

// Returns a value only if the argument references a non-blank string.
//  NOTE: panics unless `arg` is a string or string pointer.
func asStringMarshal(arg any) *string {
    var str string
    switch v := arg.(type) {
        case string:    str = v
        case *string:   if (v != nil) { str = *v }
        default:        panic(fmt.Errorf("unexpected: %v", v))
    }
    if str == "" {
        return nil
    } else {
        return &str
    }
}

// Returns a value only if the argument references a non-zero integer.
//  NOTE: panics unless `arg` is an integer or integer pointer.
func asIntMarshal(arg any) *int {
    var i int
    switch v := arg.(type) {
        case int:    i = v
        case *int:   if (v != nil) { i = *v }
        default:     panic(fmt.Errorf("unexpected: %v", v))
    }
    if i == 0 {
        return nil
    } else {
        return &i
    }
}

// Returns a value only if the argument references a valid time.
//  NOTE: panics unless `arg` is a Time or Time pointer.
func asTimeMarshal(arg any) *Time {
    var t *Time
    switch v := arg.(type) {
        case Time:  t = &v
        case *Time: t = v
        default:    panic(fmt.Errorf("unexpected: %v", v))
    }
    if BogusTime(t) {
        return nil
    } else {
        return t
    }
}

// Returns a value only if the argument references a valid date.
//  NOTE: panics unless `arg` is a Date or Date pointer.
func asDateMarshal(arg any) *Date {
    var d *Date
    switch v := arg.(type) {
        case Date:  d = &v
        case *Date: d = v
        default:    panic(fmt.Errorf("unexpected: %v", v))
    }
    if BogusDate(d) {
        return nil
    } else {
        return d
    }
}
