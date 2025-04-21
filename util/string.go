// util/string.go
//
// Functions supporting string manipulation.

package util

import (
	"strings"
)

// ============================================================================
// Exported functions
// ============================================================================

// Return the number of displayable characters in the string.
func CharCount(s string) int {
    return len([]rune(s))
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
