// Github/client.go
//
// Application Client type backed by a github.Client object.

package Github

import (
	"net/url"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// Application object referencing a GitHub client object.
type Client struct {
    ptr *github.Client
}

// ============================================================================
// Internal variables
// ============================================================================

var mainClient *Client

// ============================================================================
// Exported functions
// ============================================================================

// The default client used for application objects which do not specify one.
func MainClient() *Client {
    if mainClient == nil {
        mainClient = NewClient()
    }
    return mainClient
}

// Get a new authorized Client instance.
//  NOTE: never returns nil
func NewClient() *Client {
    token  := authToken()
    client := github.NewClient(nil).WithAuthToken(token)
    return &Client{ptr: client}
}

// ============================================================================
// Exported members - properties
// ============================================================================

// Return the underlying BaseURL field.
func (c *Client) BaseURL() *url.URL {
    if noClient(c) { return nil }
    return c.ptr.BaseURL
}

// ============================================================================
// Internal functions
// ============================================================================

// Indicate whether the argument is missing a Client object.
func noClient(c *Client) bool {
    return (c == nil) || (c.ptr == nil)
}

// ============================================================================
// Exported members
// ============================================================================

// Get all GitHub Repository objects for ORG.
func (c *Client) GetRepos() []*Repository {
    return c.GetOrgRepos(ORG)
}

// Get all GitHub Repository objects for org.
//  NOTE: if `org` is blank it defaults to ORG
func (c *Client) GetOrgRepos(org string) []*Repository {
    items, _ := getOrgRepos(c.ptr, org)
    result := make([]*Repository, 0, len(items))
    for _, repo := range items {
        result = append(result, AsRepositoryType(c, repo))
    }
    return result
}

// Get a Repository instance by fetching a GitHub Repository object.
//  NOTE: panics if the repository does not match `owner` and `name`.
//  NOTE: returns nil if `name` is blank
//  NOTE: returns nil on error
func (c *Client) GetRepository(owner, repo string) *Repository {
    return GetRepository(c, owner, repo)
}

// Get a User instance by fetching a GitHub User object.
//  NOTE: panics if the user does not match `login`.
//  NOTE: returns nil on error
func (c *Client) GetUser(login string) *User {
    return GetUser(c, login)
}
