// Github/issue.go
//
// Application Issue type backed by a github.Issue object.

package Github

import (
	"fmt"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// Application object referencing a GitHub issue object.
type Issue struct {
    Number int
    ptr    *github.Issue
    repo   *Repository
    client *Client
}

// ============================================================================
// Exported variables
// ============================================================================

// A stand-in representing a nil Issue.
var NilIssue = &Issue{}

// ============================================================================
// Exported functions
// ============================================================================

// Generate an Issue type instance.
//  NOTE: never returns nil
func NewIssueType(client *Client, owner, repo string, issue *github.Issue) *Issue {
    if repo == ""          { panic(ERR_NO_REPO_NAME) }
    if issue == nil        { panic(ERR_NIL_ISSUE) }
    if issue.Number == nil { panic(fmt.Errorf("%s in %v", ERR_NO_ISSUE_NUMBER, *issue)) }
    return &Issue{
        Number: *issue.Number,
        ptr:    issue,
        client: client,
        repo:   NewRepositoryType(client, owner, repo),
    }
}

// On GitHub, create a new issue on the indicated repository.
//  NOTE: returns 0 if finished; returns the import request ID otherwise.
func ImportIssue(client *Client, owner, repo string, imp *IssueImport, comments ...*CommentImport) int {
    cli := client.ptr
    com := Comments(comments)
    var iss *github.IssueImport
    if imp != nil { iss = &imp.IssueImport }
    return importIssue(cli, owner, repo, iss, com...)
}

// On GitHub, create a new issue on the indicated repository.
//  NOTE: returns nil on error
func CreateIssue(client *Client, owner, repo string, req *github.IssueRequest) *Issue {
    if issue := createIssue(client.ptr, owner, repo, req); issue == nil {
        return nil
    } else {
        return NewIssueType(client, owner, repo, issue)
    }
}

// Create a new Issue type instance from an issue on GitHub.
//  NOTE: returns nil on error
func GetIssue(client *Client, owner, repo string, number int) *Issue {
    if issue := getIssue(client.ptr, owner, repo, number); issue == nil {
        return nil
    } else {
        return NewIssueType(client, owner, repo, issue)
    }
}

// Fetch all issues from GitHub for the indicated repository.
//  NOTE: if an error was encountered a partial list may be returned
func GetIssues(client *Client, owner, repo string) []*Issue {
    result := []*Issue{}
    items, _ := getRepoIssues(client.ptr, owner, repo)
    if size := len(items); size > 0 {
        result = make([]*Issue, 0, size)
        for _, issue := range items {
            obj := NewIssueType(client, owner, repo, issue)
            result = append(result, obj)
        }
    }
    return result
}

// Remove the indicated repository issue from GitHub.
func DeleteIssue(client *Client, owner, repo string, number int) bool {
    return deleteIssue(client.ptr, owner, repo, number)
}

// Remove all issues from a GitHub repository.
func DeleteIssues(client *Client, owner, repo string) (count int) {
    items, _ := getRepoIssues(client.ptr, owner, repo)
    for _, issue := range items {
        if nodeId := issue.NodeID; nodeId != nil {
            if GqlDeleteIssue(*nodeId) {
                count++
            }
        }
    }
    return
}

// ============================================================================
// Exported members - properties
// ============================================================================

// Return the underlying Title value or an empty string.
func (i *Issue) Title() string {
    if noIssue(i) || (i.ptr.Title == nil) { return "" }
    return *i.ptr.Title
}

// Return the underlying Body value or an empty string.
func (i *Issue) Body() string {
    if noIssue(i) || (i.ptr.Body == nil) { return "" }
    return *i.ptr.Body
}

// Return the underlying State value or an empty string.
func (i *Issue) State() string {
    if noIssue(i) || (i.ptr.State == nil) { return "" }
    return *i.ptr.State
}

// Return the underlying StateReason value or an empty string.
func (i *Issue) StateReason() string {
    if noIssue(i) || (i.ptr.StateReason == nil) { return "" }
    return *i.ptr.StateReason
}

// Return the underlying Labels.
func (i *Issue) Labels() []string {
    if noIssue(i) || (i.ptr.Labels == nil) { return []string{} }
    return LabelStrings(i.ptr.Labels)
}

// Return the underlying Assignee value or an empty string.
func (i *Issue) Assignee() string {
    if noIssue(i) || (i.ptr.Assignee == nil) { return "" }
    return Account(i.ptr.Assignee)
}

// Return the underlying Assignees.
func (i *Issue) Assignees() []string {
    if noIssue(i) || (i.ptr.Assignees == nil) { return []string{} }
    return Accounts(i.ptr.Assignees)
}

// Return the underlying Milestone value or 0.
func (i *Issue) Milestone() int {
    if noIssue(i) || (i.ptr.Milestone == nil) { return 0 }
    return *i.ptr.Milestone.Number
}

// Return the underlying ClosedAt value or nilTime.
func (i *Issue) ClosedAt() Time {
    if noIssue(i) || NilTime(i.ptr.ClosedAt) { return nilTime }
    return *i.ptr.ClosedAt
}

// Return the underlying CreatedAt value or nilTime.
func (i *Issue) CreatedAt() Time {
    if noIssue(i) || NilTime(i.ptr.CreatedAt) { return nilTime }
    return *i.ptr.CreatedAt
}

// Return the underlying UpdatedAt value or nilTime.
func (i *Issue) UpdatedAt() Time {
    if noIssue(i) || NilTime(i.ptr.UpdatedAt) { return nilTime }
    return *i.ptr.UpdatedAt
}

// ============================================================================
// Internal functions
// ============================================================================

// Indicate whether the argument is missing an Issue object.
func noIssue(i *Issue) bool {
    return (i == nil) || (i.ptr == nil)
}

// ============================================================================
// Exported members - comments
// ============================================================================

// Get all comments associated with the issue.
func (i *Issue) Comments() []*Comment {
    res   := []*Comment{}
    cli   := i.client
    owner := i.repo.Owner
    repo  := i.repo.Name
    issue := i.Number
    if com, err := getIssueComments(cli.ptr, owner, repo, issue); err == nil {
        res = make([]*Comment, 0, len(com))
        for _, comment := range com {
            repository := NewRepositoryType(cli, owner, repo)
            res = append(res, NewCommentType(cli, repository, comment))
        }
    }
    return res
}

// From GitHub, create a Comment type instance from an IssueComment object.
func (i *Issue) GetComment(id int64) *Comment {
    cli   := i.client
    owner := i.repo.Owner
    repo  := i.repo.Name
    return GetComment(cli, owner, repo, id)
}

// On GitHub, create a comment associated with the issue.
//  NOTE: setting User in github.IssueComment has no effect!
func (i *Issue) CreateCommentFrom(text string) *Comment {
    src := &github.IssueComment{Body: &text}
    return i.CreateComment(src)
}

// On GitHub, create a comment associated with the issue.
//  NOTE: setting User in github.IssueComment has no effect!
func (i *Issue) CreateComment(src *github.IssueComment) *Comment {
    cli   := i.client
    owner := i.repo.Owner
    repo  := i.repo.Name
    issue := i.Number
    return CreateComment(cli, owner, repo, issue, src)
}

// On GitHub, delete an issue comment object.
func (i *Issue) DeleteComment(id int64) {
    cli   := i.client
    owner := i.repo.Owner
    repo  := i.repo.Name
    DeleteComment(cli, owner, repo, id)
}

// ============================================================================
// Exported members - rendering
// ============================================================================

// Render details about the instance.
func (i *Issue) Details() string {
    return issueDetails(i.ptr)
}
