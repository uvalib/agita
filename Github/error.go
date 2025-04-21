// Github/error.go

package Github

// ============================================================================
// Exported constants
// ============================================================================

// Applicaton error
const (
    ERR_NIL_CLIENT          = "nil Client"
    ERR_NIL_ISSUE           = "nil Issue"
    ERR_NIL_ISSUE_COMMENT   = "nil IssueComment"
    ERR_NIL_REPO            = "nil Repository"
    ERR_NO_DATA             = "no data"
    ERR_NO_ISSUE_REQUEST    = "no IssueRequest provided"
    ERR_NO_ISSUE_COMMENT    = "no IssueComment provided"
    ERR_NO_ISSUE_IMPORT     = "no IssueImport provided"
    ERR_NO_LOGIN            = "missing login name"
    ERR_NO_REPO             = "missing repo"
    ERR_NO_REPO_NAME        = "missing repo name"
    ERR_NO_REPO_GIVEN       = "no repository name given"
    ERR_NO_ISSUE_NUMBER     = "no issue number"
)

// GitHub error
const (
    ERR_NOT_FOUND           = "Not Found"
    ERR_INVALID_CREATE      = "Validation Failed: "
    ERR_INVALID_IMPORT      = "Validation Failed: ; "
)
