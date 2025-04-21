// util/time.go
//
// Functions supporting time and date.

package util

import (
	"fmt"
	"time"
)

// ============================================================================
// Exported types
// ============================================================================

type TimeArg interface { time.Time | *time.Time }

// ============================================================================
// Exported functions
// ============================================================================

// Get UTC time.  E.g.:
//
//  (2018)                    -> 2018-01-01 00:00:00 +0000 UTC
//  (2018,4)                  -> 2018-04-01 00:00:00 +0000 UTC
//  (2018,4,18)               -> 2018-04-18 00:00:00 +0000 UTC
//  (2018,4,18,12)            -> 2018-04-18 12:00:00 +0000 UTC
//  (2018,4,18,12,13)         -> 2018-04-18 12:13:00 +0000 UTC
//  (2018,4,18,12,13,14)      -> 2018-04-18 12:13:14 +0000 UTC
//  (2018,4,18,12,13,14,1500) -> 2018-04-18 12:13:14.0000015 +0000 UTC
//
func UtcTimeFor(args ...int) time.Time {
    return GetTimeFor(time.UTC, args...)
}

// Get local time.  E.g.:
//
//  (2018)                    -> 2018-01-01 00:00:00 -0500 EST
//  (2018,4)                  -> 2018-04-01 00:00:00 -0400 EDT
//  (2018,4,18)               -> 2018-04-18 00:00:00 -0400 EDT
//  (2018,4,18,12)            -> 2018-04-18 12:00:00 -0400 EDT
//  (2018,4,18,12,13)         -> 2018-04-18 12:13:00 -0400 EDT
//  (2018,4,18,12,13,14)      -> 2018-04-18 12:13:14 -0400 EDT
//  (2018,4,18,12,13,14,1500) -> 2018-04-18 12:13:14.0000015 -0400 EDT
//
func LocalTimeFor(args ...int) time.Time {
    return GetTimeFor(time.Local, args...)
}

// Get time.  E.g.:
//
//  (time.UTC,2018)                    -> 2018-01-01 00:00:00 +0000 UTC
//  (time.UTC,2018,4)                  -> 2018-04-01 00:00:00 +0000 UTC
//  (time.UTC,2018,4,18)               -> 2018-04-18 00:00:00 +0000 UTC
//  (time.UTC,2018,4,18,12)            -> 2018-04-18 12:00:00 +0000 UTC
//  (time.UTC,2018,4,18,12,13)         -> 2018-04-18 12:13:00 +0000 UTC
//  (time.UTC,2018,4,18,12,13,14)      -> 2018-04-18 12:13:14 +0000 UTC
//  (time.UTC,2018,4,18,12,13,14,1500) -> 2018-04-18 12:13:14.0000015 +0000 UTC
//
func GetTimeFor(loc *time.Location, args ...int) time.Time {
    var YY, MM, DD, hh, mm, ss, ns int
    size := len(args)
    if size > 0 { YY = args[0] }
    if size > 1 { MM = args[1] }
    if size > 2 { DD = args[2] }
    if size > 3 { hh = args[3] }
    if size > 4 { mm = args[4] }
    if size > 5 { ss = args[5] }
    if size > 6 { ns = args[6] }
    if MM == 0  { MM = 1 }
    if DD == 0  { DD = 1 }
    return time.Date(YY, time.Month(MM), DD, hh, mm, ss, ns, loc)
}

// Indicate whether two time values are equal.
func TimeEqual[T1, T2 TimeArg](arg1 T1, arg2 T2) bool {
    t1, t2 := timePtr(arg1), timePtr(arg2)
    if nil1, nil2 := (t1 == nil), (t2 == nil); nil1 || nil2 {
        return nil1 && nil2
    } else {
        return t1.Compare(*t2) == 0
    }
}

// ============================================================================
// Internal functions
// ============================================================================

// Normalize to time pointer.
func timePtr[T TimeArg](arg T) *time.Time {
    switch v := any(arg).(type) {
        case time.Time:  return &v
        case *time.Time: return v
        default:         panic(fmt.Errorf("unexpected: %v", arg))
    }
}
