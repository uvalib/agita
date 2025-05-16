// Github/repository.go
//
// Application Repository type backed by a github.Repository object.

package Github

import (
	"testing"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// Application object referencing a GitHub repository object.
type Repository struct {
    Owner  string
    Name   string
    ptr    *github.Repository
    client *Client
}

// ============================================================================
// Exported functions
// ============================================================================

// Generate a Repository type instance.
//  NOTE: if `owner` is blank it defaults to ORG
//  NOTE: panics if `name` is blank
//  NOTE: never returns nil
func NewRepositoryType(client *Client, owner, name string) *Repository {
    if client == nil { panic(ERR_NIL_CLIENT) }
    if name   == ""  { panic(ERR_NO_REPO_NAME) }
    return &Repository{Owner: OrgOwner(owner), Name: name, client: client}
}

// Generate a Repository type instance from a GitHub Repository object.
//  NOTE: panics if `repo` is nil
//  NOTE: never returns nil
func AsRepositoryType(client *Client, repo *github.Repository) *Repository {
    if repo == nil { panic(ERR_NO_REPO) }
    owner := Account(repo.Owner)
    name  := *repo.Name
    return NewRepositoryType(client, owner, name).initRepo(repo)
}

// Get a Repository instance by fetching a GitHub Repository object.
//  NOTE: panics if the repository does not match `owner` and `name`.
//  NOTE: returns nil if `name` is blank
//  NOTE: returns nil on error
func GetRepository(client *Client, owner, name string, silent bool) *Repository {
    if repo := getRepository(client.ptr, owner, name, silent); repo == nil {
        return nil
    } else {
        return NewRepositoryType(client, owner, name).initRepo(repo)
    }
}

// ============================================================================
// Exported functions
// ============================================================================

// Create a new GitHub repository with the provided properties.
//  NOTE: panics if `data.Name` is not present
//  NOTE: defaults `data.Owner` to ORG
func CreateRepository(client *Client, data *RepositoryRequest) *Repository {
    if repo := createRepository(client.ptr, data); repo == nil {
        return nil
    } else {
        owner := Account(repo.Owner)
        name  := ""
        if repo.Name != nil {
            name = *repo.Name
        }
        return NewRepositoryType(client, owner, name).initRepo(repo)
    }
}

// Remove an existing GitHub repository.
//  NOTE: this is only supported during testing
//  NOTE: panics if `data.Name` is not present
//  NOTE: defaults `data.Owner` to ORG
func DeleteRepository(client *Client, owner, name string) {
    if !testing.Testing() { panic("can only delete repositories in test") }
    deleteRepository(client.ptr, owner, name)
}

// ============================================================================
// Internal members
// ============================================================================

// Set `r.ptr` with an already-validated `repo`.
func (r *Repository) initRepo(repo *github.Repository) *Repository {
    r.ptr = repo
    return r
}

// ============================================================================
// Exported members - properties
// ============================================================================

// Return the underlying repository object, fetching it if necessary.
func (r *Repository) Repo() *github.Repository {
    if r.ptr == nil {
        r.ptr = getRepository(r.client.ptr, r.Owner, r.Name, false)
    }
    return r.ptr
}

// ============================================================================
// Exported members - issues
// ============================================================================

// On GitHub, create a new issue for the repository.
//  NOTE: returns 0 if finished; returns the import request ID otherwise.
func (r *Repository) ImportIssue(imp *IssueImport, comments ...*CommentImport) int {
    return ImportIssue(r.client, r.Owner, r.Name, imp, comments...)
}

// On GitHub, create a new issue for the repository.
//  NOTE: returns nil on error
func (r *Repository) CreateIssue(req *github.IssueRequest) *Issue {
    return CreateIssue(r.client, r.Owner, r.Name, req)
}

// Fetch an Issue from GitHub.
//  NOTE: returns nil on error
func (r *Repository) GetIssue(number int) *Issue {
    return GetIssue(r.client, r.Owner, r.Name, number)
}

// Fetch all issues from GitHub.
//  NOTE: if an error was encountered a partial list may be returned
func (r *Repository) GetIssues() []*Issue {
    return GetIssues(r.client, r.Owner, r.Name)
}

// ============================================================================
// Exported members - rendering
// ============================================================================

// Render details about the instance.
func (r *Repository) Details() string {
    return repoDetails(r.ptr)
}
