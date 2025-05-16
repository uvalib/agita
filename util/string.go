// util/string.go
//
// Functions supporting string manipulation.

package util

import (
	"strings"
)

// ============================================================================
// Exported constants
// ============================================================================

const WHITE_SPACE = " \t\v\n\r"

// ============================================================================
// Exported functions
// ============================================================================

// Return the number of displayable characters in the string.
func CharCount(s string) int {
    return len([]rune(s))
}

// Return a copy of the string with leading and trailing whitespace removed.
func Strip(s string) string {
    return strings.TrimSpace(s)
}

// Return a copy of the string with leading whitespace removed.
func StripLeft(s string) string {
    return strings.TrimLeft(s, WHITE_SPACE)
}

// Return a copy of the string with trailing whitespace removed.
func StripRight(s string) string {
    return strings.TrimRight(s, WHITE_SPACE)
}

// Return the string with only the first letter modified.
func UpcaseFirst(str string) string {
    return strings.ToUpper(str[0:1]) + str[1:]
}

// Modify str by appending a note.
func AppendNote(str string, note string) string {
    sep := ""
    if strings.HasSuffix(str, "\n") {
        note += "\n"
    } else if strings.ContainsAny(str, " \t\n\r\f\v") {
        if sep = " - "; !strings.Contains(str, sep) { sep = " " }
    } else {
        if sep = "-";   !strings.Contains(str, sep) { sep = "_" }
    }
    return str + sep + note
}
