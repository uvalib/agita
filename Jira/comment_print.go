// Jira/comment_print.go
//
// Reporting on jira.Comment objects.

package Jira

import (
    "fmt"
)

// ============================================================================
// Exported methods - reporting
// ============================================================================

// Show details about the instance.
func (c *Comment) Print() {
    id := MISSING
    if (c.ptr != nil) && (c.ptr.ID != "") { id = c.ptr.ID }
    fmt.Printf("\n*** JIRA Comment %s:\n", id)
    fmt.Println(c.Details())
}
