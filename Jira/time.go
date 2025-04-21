// Jira/time.go

package Jira

import (
	"fmt"
	"time"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

type Time = jira.Time
type Date = jira.Date

type TimeArg interface { Time | *Time }
type DateArg interface { Date | *Date }

// ============================================================================
// Exported constants
// ============================================================================

// A rendered date field for a missing date value.
const NO_TIME = "0001-01-01T00:00:00.000+0000"

// A rendered date field for a missing date value.
const NO_DATE = "0001-01-01"

// ============================================================================
// Internal variables
// ============================================================================

var nilTime = Time{}
var nilDate = Date{}

// ============================================================================
// Exported functions
// ============================================================================

// Cast a Jira timestamp as a time.Time.
func AsTime[T jira.Time|jira.Date](arg T) time.Time {
    return time.Time(arg)
}

// ============================================================================
// Exported functions - time
// ============================================================================

// String representation of a Jira time.
func TimeString(t jira.Time) string {
    return AsTime(t).Format("2006-01-02T15:04:05.000-0700")
}

// Indicate whether the argument is a missing or empty Jira time object.
func NilTime[T TimeArg](arg T) (result bool) {
    switch v := any(arg).(type) {
        case Time:  result = (v == nilTime)
        case *Time: result = (v == nil) || (*v == nilTime)
        default:    panic(fmt.Errorf("unexpected: %v", v))
    }
    return
}

// Indicate whether the argument represents a valid Jira time value, which is
// non-blank and does not resolve to NO_TIME.
func ValidTime[T string|*string|TimeArg](arg T) bool {
    s := NO_TIME
    switch v := any(arg).(type) {
        case string:  s = v
        case *string: if v != nil { s = *v }
        case Time:    if !NilTime(v) { s = TimeString(v) }
        case *Time:   if !NilTime(v) { s = TimeString(*v) }
        default:      panic(fmt.Errorf("unexpected: %v", v))
    }
    return s != NO_TIME
}

// Indicate whether the argument represents an invalid Jira time value, either
// because the argument is blank or because it resolves to NO_TIME.
func BogusTime[T string|*string|TimeArg](arg T) bool {
    return !ValidTime(arg)
}

// ============================================================================
// Exported functions - date
// ============================================================================

// String representation of a Jira date.
func DateString(d jira.Date) string {
    return AsTime(d).Format("2006-01-02")
}

// Indicate whether the argument is a missing or empty Jira date object.
func NilDate[T DateArg](arg T) (result bool) {
    switch v := any(arg).(type) {
        case Date:  result = (v == nilDate)
        case *Date: result = (v == nil) || (*v == nilDate)
        default:    panic(fmt.Errorf("unexpected: %v", v))
    }
    return
}

// Indicate whether the argument represents a valid Jira date value, which is
// non-blank and does not resolve to NO_DATE.
func ValidDate[T string|*string|DateArg](arg T) bool {
    s := NO_DATE
    switch v := any(arg).(type) {
        case string:  s = v
        case *string: if v != nil { s = *v }
        case Date:    if !NilDate(v) { s = DateString(v) }
        case *Date:   if !NilDate(v) { s = DateString(*v) }
        default:      panic(fmt.Errorf("unexpected: %v", v))
    }
    return s != NO_DATE
}

// Indicate whether the argument represents an invalid Jira date value, either
// because the argument is blank or because it resolves to NO_DATE.
func BogusDate[T string|*string|DateArg](arg T) bool {
    return !ValidDate(arg)
}
