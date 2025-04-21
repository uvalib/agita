// Github/user_print.go
//
// Reporting on GitHub user objects.

package Github

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported members
// ============================================================================

// Show details about the instance.
func (u *User) Print() {
    login := MISSING
    if (u.ptr != nil) && (u.ptr.Login != nil) { login = *u.ptr.Login }
    fmt.Printf("\n*** GITHUB User %s:\n", login)
    fmt.Println(u.Details())
}

// Show all organizations for user.
func (u *User) PrintOrgs() {
    printUserOrgs(u.client.ptr, u.Login)
}

// Show all repositories for the user.
func (u *User) PrintRepos() {
    printUserRepos(u.client.ptr, u.Login)
}

// ============================================================================
// Internal functions
// ============================================================================

// Show all organizations for user.
func printUserOrgs(client *github.Client, user string) {
    if items, err := getUserOrgs(client, user); err == nil {
        fmt.Printf("\n*** orgs for user %q (%d)\n", user, len(items))
        for _, org := range items {
            fmt.Printf("\t%v\n", org)
        }
    }
}

// Show all repositories for user.
func printUserRepos(client *github.Client, user string) {
    if items, err := getUserRepos(client, user); err == nil {
        fmt.Printf("\n*** all repos for user %q (%d)\n", user, len(items))
        util.PrintSortedMap(RepoMap(items))
    }
}
