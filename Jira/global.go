// Jira/global.go
//
// Values used throughout the package.

package Jira

// ============================================================================
// Exported constants
// ============================================================================

// All Jira projects are rooted here.
const BASE_URL = "https://jira.admin.virginia.edu/"

// Apparently the most Jira will return per page.
const MAX_PER_PAGE = 1000

// An output marker indicating a missing value.
const MISSING = "-"

const ____ = false

// ============================================================================
// Exported type
// ============================================================================

type AppObject interface {
    Details() string                // Render as text listing.
    MarshalJSON() ([]byte, error)   // Render as JSON.
}

// ============================================================================
// Module initialization
// ============================================================================

// Prepare for interaction with Jira.
func Initialize() bool {
    setupClient()
    setupProject()
    setupIssue()
    setupComment()
    return true
}

// Called by the system to initialize this module.
func init() {
	Initialize()
}
