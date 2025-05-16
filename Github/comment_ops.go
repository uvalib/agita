// Github/comment_ops.go
//
// Operations on github.IssueComment objects.

package Github

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Internal functions
// ============================================================================

// From GitHub, get the indicated issue comment object.
//  NOTE: returns nil on error
func getComment(client *github.Client, owner, repo string, id int64) *github.IssueComment {
    res, rsp, err := client.Issues.GetComment(ctx, owner, repo, id)
    extractRateLimit(rsp)
    log.ErrorValue(err)
    return res
}

// On GitHub, create a comment object associated with the indicated issue.
//  NOTE: returns nil on error
func createComment(client *github.Client, owner, repo string, issue int, src *github.IssueComment) *github.IssueComment {
    res, rsp, err := client.Issues.CreateComment(ctx, owner, repo, issue, src)
    extractRateLimit(rsp)
    log.ErrorValue(err)
    return res
}

// On GitHub, delete an issue comment object.
func deleteComment(client *github.Client, owner, repo string, commentId int64) {
    rsp, err := client.Issues.DeleteComment(ctx, owner, repo, commentId)
    extractRateLimit(rsp)
    log.ErrorValue(err)
}

// ============================================================================
// Internal functions - rendering
// ============================================================================

const githubCommentFieldCount = 11

// Render a GitHub issue comment object as a multiline string.
func commentDetails(c *github.IssueComment) string {
    if c == nil { return "" }
    res := make([]string, 0, githubCommentFieldCount)
    max := util.CharCount("AuthorAssociation")
    add := func(key string, val any) {
        res = append(res, fmt.Sprintf("%-*s %v", max, key, val))
    }

    if c.ID                != nil { add("ID",                   *c.ID) }
    if c.NodeID            != nil { add("NodeID",               *c.NodeID) }
    if c.User              != nil { add("User",                 UserLabel(c.User)) }
    if c.Reactions         != nil { add("Reactions",            ReactionString(c.Reactions)) }
    if c.CreatedAt         != nil { add("CreatedAt",            *c.CreatedAt) }
    if c.UpdatedAt         != nil { add("UpdatedAt",            *c.UpdatedAt) }
    if c.AuthorAssociation != nil { add("AuthorAssociation",    *c.AuthorAssociation) }
    if c.URL               != nil { add("URL",                  *c.URL) }
    if c.HTMLURL           != nil { add("HTMLURL",              *c.HTMLURL) }
    if c.IssueURL          != nil { add("IssueURL",             *c.IssueURL) }
    if c.Body              != nil { add("Body",                 *c.Body) }

    return strings.Join(res, "\n")
}

// ============================================================================
// Module initialization
// ============================================================================

// Initialize variables related to GitHub issue comments.
func setupComment() {
    // no-op
}
