// Github/client_print.go
//
// Reporting on client objects.

package Github

import (
)

// ============================================================================
// Exported members
// ============================================================================

// Show all repositories for org.
func (c *Client) PrintOrgRepos(org string) {
    printOrgRepos(c.ptr, org)
}
