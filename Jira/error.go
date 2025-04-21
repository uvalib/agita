// Jira/error.go

package Jira

import (
    "lib.virginia.edu/agita/re"
)

// ============================================================================
// Exported constants
// ============================================================================

// Applicaton error
const (
    ERR_NIL_CLIENT      = "nil Client"
    ERR_NIL_COMMENT     = "nil Comment"
    ERR_NIL_ISSUE       = "nil Issue"
    ERR_NIL_PROJECT     = "nil Project"
    ERR_NO_ISSUE        = "no Issue provided"
)

// Jira error
const (
    ERR_NO_PROJECT      = re.Pattern(`^No project `)
    ERR_REQUEST_FAILED  = re.Pattern(`^request failed\.`)
)
