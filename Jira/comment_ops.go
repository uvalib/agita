// Jira/comment_ops.go
//
// Operations on jira.Comment objects.

package Jira

import (
	"fmt"
	"strings"
	"testing"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

type CommentId = int

// ============================================================================
// Internal functions
// ============================================================================

// Get all comments for the indicated issue.
//  NOTE: may return partial results on error
func getComments(client *jira.Client, issue IssueKey) []jira.Comment {
    result   := []jira.Comment{}
    urlStr   := fmt.Sprintf("rest/api/2/issue/%s/comment", issue)
    req, err := client.NewRequest("GET", urlStr, nil)
    if log.ErrorValue(err) == nil {
        buffer := jira.Comments{}
        _, err = client.Do(req, &buffer)
        if log.ErrorValue(err) == nil {
            for _, comment := range buffer.Comments {
                result = append(result, *comment)
            }
        }
    }
    return result
}

// Get the comment with the given comment ID.
//  NOTE: returns nil on error
func getCommentById(client *jira.Client, issue IssueKey, id CommentId) (result *jira.Comment) {
    urlStr := fmt.Sprintf("rest/api/2/issue/%s/comment/%d", issue, id)
    req, err := client.NewRequest("GET", urlStr, nil)
    if log.ErrorValue(err) == nil {
        buffer := jira.Comment{}
        _, err = client.Do(req, &buffer)
        if log.ErrorValue(err) == nil {
            result = &buffer
        }
    }
    return
}

// ============================================================================
// Internal functions - rendering
// ============================================================================

const jiraCommentFieldCount = 9

// Render a Jira issue comment object as a multiline string.
//  @see COMMENT_MARSHAL
func commentDetails(c *jira.Comment) string {
    if c == nil { return "" }
    res := make([]string, 0, jiraCommentFieldCount)
    wid := util.CharCount("UpdateAuthor")
    add := func(key string, val any) {
        res = append(res, fmt.Sprintf("%-*s %v", wid, key, val))
    }

    var Author, UpdateAuthor string
    if testing.Testing() {
        Author       = Account(c.Author)
        UpdateAuthor = Account(c.UpdateAuthor)
    } else {
        Author       = UserLabel(&c.Author)
        UpdateAuthor = UserLabel(&c.UpdateAuthor)
    }

    if c.ID              != "" { add("ID",              c.ID) }
    if c.Self            != "" { add("Self",            c.Self) }
    if c.Name            != "" { add("Name",            c.Name) }
    if Author            != "" { add("Author",          Author) }
    if UpdateAuthor      != "" { add("UpdateAuthor",    UpdateAuthor) }
    if c.Updated         != "" { add("Updated",         c.Updated) }
    if c.Created         != "" { add("Created",         c.Created) }
    if c.Visibility.Type != "" { add("Visibility",      c.Visibility) }
    if c.Body            != "" { add("Body",            c.Body) }

    return strings.Join(res, "\n")
}

// ============================================================================
// Module initialization
// ============================================================================

// Initialize variables related to Jira comments.
func setupComment() {
    // no-op
}
