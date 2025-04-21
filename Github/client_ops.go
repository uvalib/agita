// Github/client_ops.go

package Github

import (
	"fmt"

	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Internal functions
// ============================================================================

// The Github token for authorization.
func authToken() string {
    name  := "GITHUB_TOKEN"
    value := util.Getenv(name)
    if value == "" { panic(fmt.Errorf("%s not in environment", name)) }
    return value
}

// ============================================================================
// Module initialization
// ============================================================================

// Initialize variables related to GitHub clients.
func setupClient() {
    // no-op
}
