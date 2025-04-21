// re/chars.go
//
// Regular expression support definitions.

package re

import (
	"unicode/utf8"
)

// ============================================================================
// Exported types
// ============================================================================

// A string that is intended to be handled as a regular expression.
//  NOTE: literals should be provided between backticks.
type Pattern = string

// ============================================================================
// Exported constants
// ============================================================================

// Regexp symbol.
const (
    BLANK = ""
    START = "^"
    END   = "$"
    SPACE = " "
)

// Regexp pattern.
const (
    SLASH = "/"
    WHITE_SPACE = Pattern(`\s+`)
)

// ============================================================================
// Exported functions
// ============================================================================

// Indicate whether `s` should be converted to a regular expression.
// If all regular expression characters are escaped then the string is not
// considered a candiate for use with regular expressions.
func IsPattern(s string) bool {
    if len(s) <= 1                          { return false }
    if HasRegexFlags(s)                     { return true }
    first, last := s[0:1], s[(len(s)-1):]
    if (first == START) || (last == END)    { return true }
    if (first == SLASH) && (last == SLASH)  { return true }
    return HasRegexChars(s)
}

// Indicate whether `s` is a simple matcher, requiring no regular expression
// handling.
// If all regular expression characters are escaped then the string is
// considered simple.
func IsSimple(s string) bool {
    return (len(s) <= 1) || !IsPattern(s)
}

// Indicate whether `s` has any unescaped regular expression character.
func HasRegexChars(s string) bool {
    esc := false
    for i := range len(s) {
        if esc {
            esc = false
        } else if s[i] == '\\' {
            esc = true
        } else if special(s[i]) {
            return true
        }
    }
    return false
}

// ============================================================================
// Internal functions
// ============================================================================

var specialBytes [16]byte

// Indicate whether the byte matches a regular expression special character.
func special(b byte) bool {
	return b < utf8.RuneSelf && specialBytes[b%16]&(1<<(b/16)) != 0
}

// ============================================================================
// Module initialization
// ============================================================================

// Called by the system to initialize this module.
func init() {
	for _, b := range []byte(`\.+*?()|[]{}^$`) {
		specialBytes[b%16] |= 1 << (b / 16)
	}
}
