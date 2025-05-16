// Jira/issue.go
//
// Application Issue type backed by a jira.Issue object.

package Jira

import (
	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

type Issue struct {
    ptr    *jira.Issue
    client *Client
}

// ============================================================================
// Exported functions
// ============================================================================

// Generate an Issue type instance.
//  NOTE: never returns nil
func NewIssueType(client *Client, issue *jira.Issue) *Issue {
    if client == nil { panic(ERR_NIL_CLIENT) }
    if issue  == nil { panic(ERR_NIL_ISSUE) }
    return &Issue{ptr: issue, client: client}
}

// Get the issue with the given issue key.
//  NOTE: returns nil on error
func GetIssueByKey(client *Client, key IssueKey) *Issue {
    if issue := getIssueByKey(client.ptr, key); issue == nil {
        return nil
    } else {
        return NewIssueType(client, issue)
    }
}

// ============================================================================
// Exported methods - properties
// ============================================================================

// Return the underlying Key value or an empty string.
func (i *Issue) Key() IssueKey {
    if noIssue(i) { return "" }
    return i.ptr.Key
}

// Return the underlying Type name or an empty string.
func (i *Issue) Type() string {
    if noFields(i) { return "" }
    return i.ptr.Fields.Type.Name
}

// Return the underlying Priority name or an empty string.
func (i *Issue) Priority() string {
    if noFields(i) || (i.ptr.Fields.Priority == nil) { return "" }
    return i.ptr.Fields.Priority.Name
}

// Return the underlying Summary value or an empty string.
func (i *Issue) Summary() string {
    if noFields(i) { return "" }
    return i.ptr.Fields.Summary
}

// Return the underlying Creator value or an empty string.
func (i *Issue) Creator() string {
    if noFields(i) { return "" }
    return Account(i.ptr.Fields.Creator)
}

// Return the underlying Reporter value or an empty string.
func (i *Issue) Reporter() string {
    if noFields(i) { return "" }
    return Account(i.ptr.Fields.Reporter)
}

// Return the underlying Description value or an empty string.
func (i *Issue) Description() string {
    if noFields(i) { return "" }
    return i.ptr.Fields.Description
}

// Return the underlying Created value or nilTime.
func (i *Issue) Created() Time {
    if noFields(i) { return nilTime }
    return i.ptr.Fields.Created
}

// Return the underlying Updated value or nilTime.
func (i *Issue) Updated() Time {
    if noFields(i) { return nilTime }
    return i.ptr.Fields.Updated
}

// Return the underlying Resolutiondate value or nilTime.
func (i *Issue) Resolutiondate() Time {
    if noFields(i) { return nilTime }
    return i.ptr.Fields.Resolutiondate
}

// Return the underlying Resolution name or an empty string.
func (i *Issue) Resolution() string {
    if noFields(i) || (i.ptr.Fields.Resolution == nil) { return "" }
    return i.ptr.Fields.Resolution.Name
}

// Return the underlying Status name or an empty string.
func (i *Issue) Status() string {
    if noFields(i) || (i.ptr.Fields.Status == nil) { return "" }
    return i.ptr.Fields.Status.Name
}

// Return the underlying Assignee value or an empty string.
func (i *Issue) Assignee() string {
    if noFields(i) { return "" }
    return Account(i.ptr.Fields.Assignee)
}

// Return the underlying Labels.
func (i *Issue) Labels() []string {
    if noFields(i) { return []string{} }
    return i.ptr.Fields.Labels
}

// Return the underlying Attachments.
func (i *Issue) Attachments() []*jira.Attachment {
    if noFields(i) { return []*jira.Attachment{} }
    return i.ptr.Fields.Attachments
}

// ============================================================================
// Internal functions
// ============================================================================

// Indicate whether the argument is missing an Issue object.
func noIssue(i *Issue) bool {
    return (i == nil) || (i.ptr == nil)
}

// Indicate whether the argument is missing Issue.Fields.
func noFields(i *Issue) bool {
    return (i == nil) || (i.ptr == nil) || (i.ptr.Fields == nil)
}

// ============================================================================
// Exported methods - comments
// ============================================================================

// Get all comments for the issue.
//  NOTE: may return partial results on error
func (i *Issue) Comments() []Comment {
    if (i.client == nil) || (i.client.ptr == nil) { panic(ERR_NIL_CLIENT) }
    if (i.ptr    == nil) || (i.ptr.Key    == "")  { panic(ERR_NO_ISSUE) }
    items  := getComments(i.client.ptr, i.ptr.Key)
    result := make([]Comment, 0, len(items))
    for _, comment := range items {
        result = append(result, *NewCommentType(i.client, &comment))
    }
    return result
}

// Get the comment with the given comment ID.
//  NOTE: returns nil on error
func (i *Issue) Comment(id CommentId) *Comment {
    return GetCommentById(i.client, i.ptr.Key, id)
}

// ============================================================================
// Exported methods - rendering
// ============================================================================

// Render details about the instance.
func (i *Issue) Details() string {
    return issueDetails(i.ptr)
}
