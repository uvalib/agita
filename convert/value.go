// convert/value.go
//
// Conversion of Jira field values to GitHub field values.

package convert

import (
	"lib.virginia.edu/agita/util"

	"lib.virginia.edu/agita/Github"
	"lib.virginia.edu/agita/Jira"
)

// ============================================================================
// Exported functions
// ============================================================================

// Return the (possibly) converted value and indication of whether it should be
// used as a field value.
func From(value any) (result any, useable bool) {
    use := false
    switch v := value.(type) {
        case int:        use = (v != 0)
        case string:     use = (v != "")
        case Jira.Time:  if Jira.ValidTime(v) { value, use = FromTime(v) }
        case *int:       use = (v != nil) && (*v != 0)
        case *string:    use = (v != nil) && (*v != "")
        case *Jira.Time: if Jira.ValidTime(v) { value, use = FromTime(*v) }
        default:         use = !util.IsEmpty(v)
    }
    return value, use
}

// Convert a Jira timestamp to a Github timestamp.
func FromTime(src Jira.Time) (result Github.Time, useable bool) {
    str := Jira.TimeString(src)
    res := Github.MakeTime(str)
    return res, !Github.NilTime(res)
}
