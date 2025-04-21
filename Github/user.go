// Github/user.go
//
// Application User type backed by a github.User object.

package Github

import (
	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

type User struct {
    Login  string
    ptr    *github.User
    client *Client
}

// ============================================================================
// Exported functions
// ============================================================================

// Generate a User type instance.
//  NOTE: never returns nil
func NewUserType(client *Client, login string) *User {
    if client == nil { panic(ERR_NIL_CLIENT) }
    if login  == ""  { panic(ERR_NO_LOGIN) }
    return &User{Login: login, client: client}
}

// Get a User instance by fetching a GitHub User object.
//  NOTE: panics if the user does not match `login`.
//  NOTE: returns nil on error
func GetUser(client *Client, login string) *User {
    if (client == nil) || (client.ptr == nil) { panic(ERR_NIL_CLIENT) }
    if user := getUser(client.ptr, login); user == nil {
        return nil
    } else {
        return NewUserType(client, login).initUser(user)
    }
}

// ============================================================================
// Internal members
// ============================================================================

// Set `u.ptr` with an already-validated `user`.
func (u *User) initUser(user *github.User) *User {
    u.ptr = user
    return u
}

// ============================================================================
// Exported members
// ============================================================================

// Get all GitHub Organization objects for the user.
func (u *User) Orgs() []*github.Organization {
    items, _ := getUserOrgs(u.client.ptr, u.Login)
    return items
}

// Get all GitHub Repository objects for the user.
func (u *User) Repos() []*Repository {
    items, _ := getUserRepos(u.client.ptr, u.Login)
    result := make([]*Repository, 0, len(items))
    for _, repo := range items {
        result = append(result, AsRepositoryType(u.client, repo))
    }
    return result
}

// ============================================================================
// Exported members - rendering
// ============================================================================

// Render details about the instance.
func (u *User) Details() string {
    return userDetails(u.ptr)
}
