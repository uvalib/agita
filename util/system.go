// util/system.go
//
// System state information.

package util

import (
	"os"
	"path"
	"strings"
)

// ============================================================================
// Internal constants
// ============================================================================

const thisFile = "/util/system.go"

// ============================================================================
// Internal variables
// ============================================================================

var root string
var prog string

// ============================================================================
// Exported functions
// ============================================================================

// Returns the path of the project root.
func RootPath() string {
    if root == "" {
        root = strings.TrimSuffix(FileName(), thisFile)
    }
    return root
}

// The name of this program.
func Progname() string {
    if prog == "" {
        if InDebugger() {
            prog = path.Base(RootPath())
        } else {
            prog = path.Base(os.Args[0])
        }
    }
    return prog
}

// Indicate whether the program is executing in the debugger.
func InDebugger() bool {
    return strings.Contains(os.Args[0], "__debug")
}
