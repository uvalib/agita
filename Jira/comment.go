// Jira/comment.go
//
// Application Comment type backed by a jira.Comment object.

package Jira

import (
	"strconv"

	"lib.virginia.edu/agita/log"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

// Application object referencing a Jira comment object.
type Comment struct {
    ptr    *jira.Comment
    client *Client
}

// ============================================================================
// Exported functions
// ============================================================================

// Generate a Comment type instance.
//  NOTE: never returns nil
func NewCommentType(client *Client, comment *jira.Comment) *Comment {
    if client  == nil { panic(ERR_NIL_CLIENT) }
    if comment == nil { panic(ERR_NIL_COMMENT) }
    return &Comment{ptr: comment, client: client}
}

// Get the comment with the given comment ID.
//  NOTE: returns nil on error
func GetCommentById(client *Client, issue IssueKey, id CommentId) *Comment {
    if comment := getCommentById(client.ptr, issue, id); comment == nil {
        return nil
    } else {
        return NewCommentType(client, comment)
    }
}

// ============================================================================
// Exported methods - properties
// ============================================================================

// Return the underlying ID value or 0.
func (c *Comment) ID() (result CommentId) {
    if noComment(c) || (c.ptr.ID == "") { return }
    if id, err := strconv.Atoi(c.ptr.ID); log.ErrorValue(err) == nil {
        result = id
    }
    return
}

// Return the underlying Name value or an empty string.
func (c *Comment) Name() string {
    if noComment(c) { return "" }
    return c.ptr.Name
}

// Return the underlying Author value or an empty string.
func (c *Comment) Author() string {
    if noComment(c) { return "" }
    return Account(&c.ptr.Author)
}

// Return the underlying Body value or an empty string.
func (c *Comment) Body() string {
    if noComment(c) { return "" }
    return c.ptr.Body
}

// Return the underlying UpdateAuthor value or an empty string.
func (c *Comment) UpdateAuthor() string {
    if noComment(c) { return "" }
    return Account(&c.ptr.UpdateAuthor)
}

// Return the underlying Updated value or an empty string.
func (c *Comment) Updated() string {
    if noComment(c) { return "" }
    return c.ptr.Updated
}

// Return the underlying Created value or an empty string.
func (c *Comment) Created() string {
    if noComment(c) { return "" }
    return c.ptr.Created
}

// Return the underlying Visibility value or an empty string.
func (c *Comment) Visibility() string {
    if noComment(c) { return "" }
    return c.ptr.Visibility.Value
}

// ============================================================================
// Internal functions
// ============================================================================

// Indicate whether the argument is missing a Comment object.
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
