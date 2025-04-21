// Jira/client_print.go
//
// Reporting on client objects.

package Jira

import (
)

// ============================================================================
// Exported methods
// ============================================================================

// Show projects for the organization associated with the client.
func (c *Client) PrintProjects() {
    printProjects(c.ptr)
}
