// Jira/project.go
//
// Application Project type backed by a jira.Project object.

package Jira

import (
	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

type Project struct {
    ptr    *jira.Project
    client *Client
}

// ============================================================================
// Exported functions
// ============================================================================

// Generate a Project type instance.
//  NOTE: never returns nil
func NewProjectType(client *Client, project *jira.Project) *Project {
    if client  == nil { panic(ERR_NIL_CLIENT) }
    if project == nil { panic(ERR_NIL_PROJECT) }
    return &Project{ptr: project, client: client}
}

// Get all projects for the Jira referenced by the client.
//  NOTE: all returned elements are non-nil
func GetProjects(client *Client) []*Project {
    items  := getProjects(client.ptr)
    result := make([]*Project, 0, len(items))
    for _, project := range items {
        result = append(result, NewProjectType(client, project))
    }
    return result
}

// Get the project with the given project key.
//  NOTE: returns nil on error
func GetProjectByKey(client *Client, key ProjKey) *Project {
    if project := getProjectByKey(client.ptr, key); project == nil {
        return nil
    } else {
        return NewProjectType(client, project)
    }
}

// ============================================================================
// Exported methods - properties
// ============================================================================

// Return the underlying Key value or the empty string.
func (p *Project) Key() string {
    if noProject(p) { return "" }
    return p.ptr.Key
}

// Return the underlying Name value or the empty string.
func (p *Project) Name() string {
    if noProject(p) { return "" }
    return p.ptr.Name
}

// ============================================================================
// Internal functions
// ============================================================================

// Indicate whether the argument is missing a Project object.
func noProject(p *Project) bool {
    return (p == nil) || (p.ptr == nil)
}

// ============================================================================
// Exported methods - issues
// ============================================================================

// Get all issues for the project.
//  NOTE: may return partial results on error
func (p *Project) Issues() []Issue {
    items  := getIssues(p.client.ptr, p.ptr.Key)
    result := make([]Issue, 0, len(items))
    for _, issue := range items {
        result = append(result, *NewIssueType(p.client, &issue))
    }
    return result
}

// Get the issue with the given issue key.
//  NOTE: returns nil on error
func (p *Project) GetIssue(key IssueKey) *Issue {
    return GetIssueByKey(p.client, key)
}

// ============================================================================
// Exported methods - rendering
// ============================================================================

// Render details about the instance.
func (p *Project) Details() string {
    return projectDetails(p.ptr)
}
