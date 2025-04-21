// test/test.go
//
// Test case support definitions.

// Methods supporting testing.
package test

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Constants
// ============================================================================

// Indicate whether to avoid tests which create/update temporary GitHub items.
const PASSIVE = false

// ============================================================================
// Types
// ============================================================================

type ErrorArg interface { string | re.Regex }

// ============================================================================
// Test support functions
// ============================================================================

// For use at the top of a test which requires writing to GitHub.
func Passive(label string, t *testing.T) bool {
    if PASSIVE {
        Output("%s: skipped - test.PASSIVE is true", label)
    }
    return PASSIVE
}

// Output during tests.
// Unlike `t.Logf()` the output line is not prefixed with the line number.
func Output(format string, args ...any) {
    if !strings.HasSuffix(format, "\n") {
        format += "\n"
    }
    fmt.Fprintf(os.Stdout, format, args...)
}

// Log an appropriate error message if there was a panic in the calling
// function that was not expected.  An expected panic is silently recovered
// here; an unexpected panic is propagated.
//  NOTE: must be called via `defer`
func EvaluatePanic[T ErrorArg](label string, expected T, t *testing.T) {
    var regex *re.Regex
    switch v := any(expected).(type) {
        case re.Regex: regex = &v
        case string:   if re.IsPattern(v) { regex = re.New(v) }
    }
    var expect, exFmt string
    if regex == nil {
        switch v := any(expected).(type) { case string: expect = v }
        exFmt  = "%q"
    } else {
        expect = regex.String()
        exFmt  = "`%s`"
    }
    if err := recover(); err == nil {
        if expect != "" {
            t.Errorf("%s: failed to panic with " + exFmt, label, expect)
        }
    } else {
        msg := log.PanicMessage(err)
        if expect == "" {
            t.Errorf("%s: unexpected panic: %q", label, msg)
            panic(err)
        } else if (regex != nil) && !regex.Match(msg) {
            t.Errorf("%s: panic: got %q, want: " + exFmt, label, msg, expect)
            panic(err)
        } else if (regex == nil) && (expect != msg) {
            t.Errorf("%s: panic: got %q, want: " + exFmt, label, msg, expect)
            panic(err)
        }
    }
}

// ============================================================================
// Test support functions
// ============================================================================

// Generate a test case name.
func CaseName(fn string, idx int) string {
    return fmt.Sprintf("%s-%d", fn, idx)
}

// Transform a test value into a unique string.
//  NOTE: if text is blank then the result will be blank.
func Unique(text, context string) string {
    if text != "" {
        if context != "" {
            text += " - " + context
        }
        text = util.Randomize(text)
    }
    return text
}

// Render a slice of strings as a sorted list.
func SortedListing(array []string) string {
    list := slices.Clone(array)
    slices.Sort(list)
    return strings.Join(list, "; ")
}

// Fail unless `got` and `want` are both nil or both non-nil.
//  NOTE: return false if neither are nil
func CheckForNils(fn string, got, want any, t *testing.T) bool {
    if gotNil, wantNil := (got == nil), (want == nil); gotNil || wantNil {
        if !gotNil || !wantNil {
            t.Errorf("%s() = %v, want %v", fn, got, want)
        }
        return true
    }
    return false
}
