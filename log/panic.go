// log/panic.go
//
// Panic suppression control for error logging while testing.

package log

import (
	"sync"
)

// ============================================================================
// Internal variables
// ============================================================================

var panicSuppressed = 0
var panicMutex sync.Mutex

// ============================================================================
// Exported functions
// ============================================================================

// Set or unset panic suppression and return the previous state.
func PanicSuppression(suppress bool) (wasSuppressed bool) {
    panicMutex.Lock()
    defer panicMutex.Unlock()
    wasSuppressed = PanicSuppressed()
    if suppress {
        panicSuppressed++
    } else if wasSuppressed {
        panicSuppressed--
    }
    return
}

// Set panic suppression.
func SuppressPanic() bool {
    return PanicSuppression(true)
}

// Restore panic suppression to the previous state.
// If panic suppression is already unset then this function has no effect.
func RestorePanic() bool {
    return PanicSuppression(false)
}

// Indicate whether RequestErrorIn() may panic during testing.
func PanicSuppressed() bool {
    return panicSuppressed > 0
}
