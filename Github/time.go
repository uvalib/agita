// Github/time.go

package Github

import (
	"fmt"
	"time"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

type Time = github.Timestamp

type TimeArg interface { Time | *Time }

// ============================================================================
// Internal variables
// ============================================================================

var nilTime = Time{}

// ============================================================================
// Exported functions
// ============================================================================

// Indicate whether the argument is a missing or empty GitHub time object.
func NilTime[T TimeArg](arg T) (result bool) {
    switch v := any(arg).(type) {
        case Time:  result = (v == nilTime)
        case *Time: result = (v == nil) || (*v == nilTime)
        default:    panic(fmt.Errorf("unexpected: %v", v))
    }
    return
}

// Create a Github timestamp.
//  NOTE: an invalid conversion results in an empty time value.
func MakeTime(src string) Time {
    // @see https://dev.to/luthfisauqi17/golangs-unique-way-to-parse-string-to-time-2jmk
    t, err := time.Parse("2006-01-02T15:04:05.999Z0700", src)
    if log.ErrorValue(err) == nil {
        return AsTimestamp(t)
    } else {
        return nilTime
    }
}

// Cast a time.Time as a GitHub timestamp.
func AsTimestamp(t time.Time) Time {
    return Time{Time: t}
}

// Create a Github timestamp from date/time parts per util.GetTimeFor.
func TimeFor(args ...int) Time {
    return LocalTimeFor(args...)
}

// Create a Github timestamp from date/time parts per util.GetTimeFor.
func UtcTimeFor(args ...int) Time {
    return AsTimestamp(util.UtcTimeFor(args...))
}

// Create a Github timestamp from date/time parts per util.GetTimeFor.
func LocalTimeFor(args ...int) Time {
    return AsTimestamp(util.LocalTimeFor(args...))
}
