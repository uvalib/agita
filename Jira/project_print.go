// Jira/project_print.go
//
// Reporting on jira.Project objects.

package Jira

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported methods
// ============================================================================

// Show details about the instance.
func (p *Project) Print() {
    key := MISSING
    if (p.ptr != nil) && (p.ptr.Key != "") { key = p.ptr.Key }
    fmt.Printf("\n*** JIRA Project %s:\n", key)
    fmt.Println(p.Details())
}

// ============================================================================
// Internal methods
// ============================================================================

// Show projects for the organization associated with the client.
func printProjects(client *jira.Client) {
    items := getProjects(client)
    keys  := make([]ProjKey, 0, len(items))
    for _, project := range items {
        keys = append(keys, project.Key)
    }
    fmt.Printf("\n*** Projects: %v\n", keys)
    for _, project := range items {
        fmt.Printf("\n*** Project %q:\n", project.Key)
        fmt.Println(projectDetails(project))
    }
}
