// Github/issue_import.go
//
// Full information for an issue includes the issue itself and all of its
// comments.

package Github

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// Application wrapper for the GitHub object used to import issues.
type IssueImport struct {
    github.IssueImport
}

// ============================================================================
// Exported functions
// ============================================================================

// Generate a wrapper for a Github issue import object from a table of field
// names and values.
//  NOTE: never returns nil
func NewIssueImport(fields map[string]any) *IssueImport {
    imp := github.IssueImport{}
    for key, val := range fields {
        var s string
        var i int
        var b bool
        var t Time
        var a []string
        switch v := val.(type) {
            case string:    s = v
            case int:       i = v
            case bool:      b = v
            case Time:      t = v
            case *string:   s = *v
            case *int:      i = *v
            case *bool:     b = *v
            case *Time:     t = *v
            case []string:  a = v
        }
        switch key {
            case "Title":       imp.Title       = s
            case "Body":        imp.Body        = s
            case "CreatedAt":   imp.CreatedAt   = github.Ptr(t)
            case "ClosedAt":    imp.ClosedAt    = github.Ptr(t)
            case "UpdatedAt":   imp.UpdatedAt   = github.Ptr(t)
            case "Assignee":    imp.Assignee    = github.Ptr(s)
            case "Milestone":   imp.Milestone   = github.Ptr(i)
            case "Closed":      imp.Closed      = github.Ptr(b)
            case "Labels":      imp.Labels      = a
        }
    }
    return &IssueImport{IssueImport: imp}
}

// ============================================================================
// Exported members - rendering
// ============================================================================

const githubIssueImportFieldCount = 9

// Render details about the instance.
func (i *IssueImport) Details() string {
    res := make([]string, 0, githubIssueImportFieldCount)
    max := util.CharCount("AuthorAssociation")
    add := func(key string, val any) {
        res = append(res, fmt.Sprintf("%-*s %v", max, key, val))
    }

    if i.Title     != ""    { add("Title",       i.Title) }
    if i.CreatedAt != nil   { add("CreatedAt",  *i.CreatedAt) }
    if i.ClosedAt  != nil   { add("ClosedAt",   *i.ClosedAt) }
    if i.UpdatedAt != nil   { add("UpdatedAt",  *i.UpdatedAt) }
    if i.Assignee  != nil   { add("Assignee",   *i.Assignee) }
    if i.Milestone != nil   { add("Milestone",  *i.Milestone) }
    if i.Closed    != nil   { add("Closed",     *i.Closed) }
    if len(i.Labels) > 0    { add("Labels",      i.Labels) }
    if i.Body      != ""    { add("Body",        i.Body) }

    return strings.Join(res, "\n")
}

// ============================================================================
// Internal functions
// ============================================================================

// On GitHub, create an issue and its comments on the indicated repository.
//  NOTE: returns 0 if finished; returns the import request ID otherwise.
func importIssue(client *github.Client, owner, repo string, imp *github.IssueImport, comments ...*github.Comment) int {
    if imp == nil { panic(ERR_NO_ISSUE_IMPORT) }
    req := NewIssueImportRequest(*imp, comments...).IssueImportRequest
    impRsp, rsp, err := client.IssueImport.Create(ctx, owner, repo, &req)
    pending := IsScheduled(err)
    if pending || (log.ErrorValue(err) == nil) {
        log.Info("\n*** import issue %q - rsp = %v\n", imp.Title, rsp)
    }
    if pending && (impRsp != nil) && (impRsp.ID != nil) {
        return *impRsp.ID
    } else {
        return 0
    }
}

// Determine whether the import occurred and with what status.
func checkImportIssue(client *github.Client, owner, repo string, importID int) (done bool, status string) {
    impRsp, _, err := client.IssueImport.CheckStatus(ctx, owner, repo, int64(importID))
    if !IsScheduled(err) {
        log.ErrorValue(err)
    }
    if (impRsp != nil) && (impRsp.Status != nil) {
        status = *impRsp.Status
    }
    done = (status == "imported")
    return
}
