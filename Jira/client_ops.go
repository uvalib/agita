// Jira/client_ops.go

package Jira

import (
	"fmt"

	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Internal functions
// ============================================================================

// The Jira token for authorization.
func authToken() string {
    name  := "JIRA_TOKEN"
    value := util.Getenv(name)
    if value == "" { panic(fmt.Errorf("%s not in environment", name)) }
    return value
}

// ============================================================================
// Module initialization
// ============================================================================

// Initialize variables related to Jira clients.
func setupClient() {
    // no-op
}
