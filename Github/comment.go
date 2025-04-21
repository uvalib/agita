// Github/comment.go
//
// Application Comment type backed by a github.IssueComment object.

package Github

import (
	"testing"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// Application object referencing a GitHub issue comment object.
type Comment struct {
    ptr    *github.IssueComment
    repo   *Repository
    client *Client
}

// ============================================================================
// Exported functions
// ============================================================================

// Generate a Comment type instance.
//  NOTE: panics if `client` is `comment` is nil
//  NOTE: never returns nil
func NewCommentType(client *Client, repo *Repository, comment *github.IssueComment) *Comment {
    if client  == nil { panic(ERR_NIL_CLIENT) }
    if comment == nil { panic(ERR_NIL_ISSUE_COMMENT) }
    return &Comment{ptr: comment, repo: repo, client: client}
}

// From GitHub, create a Comment type instance from an IssueComment object.
//  NOTE: returns nil on error
func GetComment(client *Client, owner, repo string, id int64) *Comment {
    if com := getComment(client.ptr, owner, repo, id); com == nil {
        return nil
    } else {
        repo := NewRepositoryType(client, owner, repo)
        return NewCommentType(client, repo, com)
    }
}

// Create a Comment type instance associated with the indicated issue.
//  NOTE: returns nil on error
func CreateComment(client *Client, owner, repo string, issue int, src *github.IssueComment) *Comment {
    if src == nil { panic(ERR_NO_ISSUE_COMMENT) }
    if com := createComment(client.ptr, owner, repo, issue, src); com == nil {
        return nil
    } else {
        repo := NewRepositoryType(client, owner, repo)
        return NewCommentType(client, repo, com)
    }
}

// On GitHub, delete an issue comment object.
//  NOTE: this is only supported during testing
func DeleteComment(client *Client, owner, repo string, commentId int64) {
    if !testing.Testing() { panic("can only delete comments in test") }
    deleteComment(client.ptr, owner, repo, commentId)
}

// ============================================================================
// Exported methods - properties
// ============================================================================

// Return the text of the comment.
func (c *Comment) Body() string {
    if noComment(c) || (c.ptr.Body == nil) { return "" }
    return *c.ptr.Body
}

// Return the ID of the comment.
func (c *Comment) ID() int64 {
    if noComment(c) || (c.ptr.ID == nil) { return 0 }
    return *c.ptr.ID
}

// ============================================================================
// Internal functions
// ============================================================================

// Indicate whether the argument is missing an IssueComment.
func noComment(c *Comment) bool {
    return (c == nil) || (c.ptr == nil)
}

// ============================================================================
// Exported methods - rendering
// ============================================================================

// Render details about the instance.
func (c *Comment) Details() string {
    return commentDetails(c.ptr)
}
