// Github/import.go
//
// Mechanism for creating a complete GitHub issue with comments.

package Github

import (
	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

type IssueImportRequest struct {
    github.IssueImportRequest
}

// ============================================================================
// Exported functions
// ============================================================================

// Generate a wrapper for a Github issue import request object.
//  NOTE: never returns nil
func NewIssueImportRequest(imp github.IssueImport, comments ...*github.Comment) *IssueImportRequest {
    req := github.IssueImportRequest{IssueImport: imp, Comments: comments}
    return &IssueImportRequest{req}
}
