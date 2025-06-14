// Github/global.go
//
// Values used throughout the package.

package Github

import (
	"context"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported constants
// ============================================================================

// Apparently the most GitHub will return per page.
const MAX_PER_PAGE = 100

// All target repos begin with "https://github.com/${ORG}".
const ORG = "uvalib"

// Whether the target platform is GitHub Enterprise.
//  NOTE: for uvalib this is false.
const ENTERPRISE = false

// Whether the owner can be set when creating an issue or comment.
//  NOTE: for regular GitHub this is false.
const CAN_SET_USER = ENTERPRISE

// An output marker indicating a missing value.
const MISSING = "-"

// The prefix for lines leading a GitHub issue body which convey Jira issue
// information that cannot be mapped on to the GitHub issue.
const ISSUE_ANNOTATION_TAG = "ORIGINAL JIRA ISSUE"

// The prefix for lines leading a GitHub comment body which convey Jira comment
// information that cannot be mapped on to the GitHub comment.
const COMMENT_ANNOTATION_TAG = "ORIGINAL JIRA COMMENT"

// ============================================================================
// Internal variables
// ============================================================================

var ctx = context.Background()

// Before the first use of fake items, remove all detected temporary GitHub items.
var preCleanAll = true

// In test teardown, remove temporary GitHub items generated by test execution.
var postCleanFakes = true

// In test teardown, remove all detected temporary GitHub items.
var postCleanAll = false

// ============================================================================
// Exported functions - GitHub
// ============================================================================

// Indicate whether the given error is actually a placeholder indicating an
// HTTP 202 Accepted was received, indicating that the requested action was
// scheduled to be performed later.
func IsScheduled(err error) bool {
    _, ok := err.(*github.AcceptedError)
    return ok
}

// ============================================================================
// Exported functions
// ============================================================================

// Prepare for interaction with GitHub.
func Initialize() bool {
    setupClient()
    setupRepository()
    setupIssue()
    setupComment()
    return true
}

// ============================================================================
// Module initialization
// ============================================================================

// Called by the system to initialize this module.
func init() {
	Initialize()
}
