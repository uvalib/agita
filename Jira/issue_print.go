// Jira/issue_print.go
//
// Reporting on jira.Issue objects.

package Jira

import (
    "fmt"
)

// ============================================================================
// Exported methods
// ============================================================================

// Show details about the instance.
func (i *Issue) Print() {
    key := MISSING
    if (i.ptr != nil) && (i.ptr.Key != "") { key = i.ptr.Key }
    fmt.Printf("\n*** JIRA Issue %s:\n", key)
    fmt.Println(i.Details())
}
