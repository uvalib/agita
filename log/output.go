// log/output.go
//
// Support for report output.

package log

import (
	"fmt"
	"strings"
	"testing"

	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Exported functions
// ============================================================================

// An informational log entry prefixed with the current function.
func Info(msg string, args ...any) {
    InfoIn(util.CallerName(), msg, args...)
}

// A warning log entry prefixed with the current function.
func Warn(msg string, args ...any) {
    WarnIn(util.CallerName(), msg, args...)
}

// An error log entry prefixed with the current function.
func Error(msg string, args ...any) {
    ErrorIn(util.CallerName(), msg, args...)
}

// An informational log entry.
func InfoIn(fn string, msg string, args ...any) {
    format := fmt.Sprintf("%s: %s", fn, msg)
    args = append([]any{format}, args...)
    logWrite(args...)
}

// A warning log entry.
func WarnIn(fn, msg string, args ...any) {
    format := fmt.Sprintf("%s: WARNING: %s", fn, msg)
    args = append([]any{format}, args...)
    logWrite(args...)
}

// An error log entry.
func ErrorIn(fn, msg string, args ...any) {
    format := fmt.Sprintf("%s: ERROR: %s", fn, msg)
    args = append([]any{format}, args...)
    logWrite(args...)
}

// ============================================================================
// Internal functions
// ============================================================================

// Write to the log.
// If `args[0]` is a format string it is used to print the rest of the args.
// Otherwise, the args are printed with Println.
func logWrite(args ...any) {
    if testing.Testing() { return }
    if len(args) > 1 {
        switch format := args[0].(type) {
            case string:
                if strings.Contains(format, "%") {
                    if !strings.HasSuffix("\n", format) { format += "\n" }
                    fmt.Printf(format, args[1:]...)
                    return
                }
        }
    }
    fmt.Println(args...)
}
