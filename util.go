// util.go
//
// Miscellaneous support functions.

package main

import (
	"fmt"
	"os"
	"strings"

	"lib.virginia.edu/agita/convert"
	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Constants
// ============================================================================

const NORMAL_EXIT = 0
const ABORT_EXIT  = 2

// ============================================================================
// Functions
// ============================================================================

// Terminate the program with a message on stderr.
func Abort(msg string, arg ...any) {
    report := ShowError(msg, arg...)
    if util.InDebugger() {
        panic("aborted: " + report)
    } else {
        os.Exit(ABORT_EXIT)
    }
}

// An error message on stderr.
func ShowError(msg string, arg ...any) string {
    if last := string(msg[len(msg)-1]); !strings.ContainsAny(last, ".!?") {
        msg += "."
    }
    msg = "ERROR: " + msg + "\n\n"
    return Show(msg, arg...)
}

// Print a line on stderr.
func Show(msg string, arg ...any) string {
    if len(arg) > 0 {
        msg = fmt.Sprintf(msg, arg...)
    }
    if !strings.HasSuffix(msg, "\n") {
        msg += "\n"
    }
    fmt.Fprint(os.Stderr, msg)
    return msg
}

// Return a list of GitHub repository names, converting any JIRA project name
// to its matching GitHub repository.
func RepositoryNames(names ...string) []string {
    for i, name := range names {
        if repo := convert.ProjectToRepo[name]; repo != "" {
            names[i] = repo
        }
    }
    return names
}
