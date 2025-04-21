// util/meta.go
//
// Functions supporting metaprogramming.

package util

import (
	"runtime"
)

// ============================================================================
// Exported functions
// ============================================================================

// Returns the file where the calling function is defined.
func FileName() string {
    _, file, _, _ := runtime.Caller(1)
    return file
}

// Returns the name of the calling function.
func FuncName() string {
    pc, _, _, _ := runtime.Caller(1)
    return funcForPC(pc)
}

// Returns the name of the calling function's caller.
func CallerName() string {
    pc, _, _, _ := runtime.Caller(2)
    return funcForPC(pc)
}

// ============================================================================
// Internal functions
// ============================================================================

// Returns the name of the function at the given program counter.
func funcForPC(pc uintptr) string {
    if fn := runtime.FuncForPC(pc); fn == nil {
        return ""
    } else {
        return fn.Name()
    }
}
