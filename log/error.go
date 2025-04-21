// log/panic.go
//
// Support for logging errors.

package log

import (
	"fmt"
	"strings"
	"testing"

	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported functions
// ============================================================================

// Report an API request error from the calling function.
func ErrorValue(err error) error {
    if err != nil {
        ErrorValueIn(util.CallerName(), err)
    }
    return err
}

// Report an API request error.
func ErrorValueIn(fn string, err error) error {
    if err != nil {
        msg := PanicMessage(err)
        if testing.Testing() && !PanicSuppressed() {
            panic(msg)
        } else {
            ErrorIn(fn, msg)
        }
    }
    return err
}

// Returns the message from an error.
func PanicMessage(err any) string {
    switch v := err.(type) {
        case github.ErrorResponse:  return GithubPanicMessage(v)
        case *github.ErrorResponse: return GithubPanicMessage(*v)
        default:                    return fmt.Sprintf("%v", err)
    }
}

// Returns blank if no message was given or could be inferred.
func GithubPanicMessage(err github.ErrorResponse) string {
    // if rsp := err.Response; rsp != nil {
    //     return rsp.Status
    // }
    msg := err.Message
    if len(err.Errors) > 0 {
        part := []string{}
        for _, e := range err.Errors {
            part = append(part, e.Message)
        }
        msg = strings.TrimRight(msg, ".,;:") + ": " + strings.Join(part, "; ")
    }
    return msg
}
