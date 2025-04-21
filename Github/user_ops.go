// Github/user_ops.go
//
// Operations on github.User objects.

package Github

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported functions
// ============================================================================

// Render a GitHub User object as a user name label.
func UserLabel(user *github.User) string {
    switch {
        case user == nil:       return "[nil]"
        case user.Login == nil: return "[missing]"
        case *user.Login == "": return "[blank]"
        default:                return *user.Login
    }
}

// Render a slice of GitHub User objects as a slice of user name labels.
func UserLabels(users []*github.User) []string {
    result := make([]string, 0, len(users))
    n, m, b := 1, 1, 1
    for _, user := range users {
        login := UserLabel(user)
        switch login {
            case "[nil]":     login = fmt.Sprintf("[nil-%d]",     n); n++
            case "[missing]": login = fmt.Sprintf("[missing-%d]", m); m++
            case "[blank]":   login = fmt.Sprintf("[blank-%d]",   b); b++
        }
        result = append(result, login)
    }
    return result
}

// ============================================================================
// Internal functions
// ============================================================================

// Get a GitHub User object.
//  NOTE: panics if the user does not match `login`.
//  NOTE: returns nil on error
func getUser(client *github.Client, login string) *github.User {
    user, _, err := client.Users.Get(ctx, login)
    if (log.ErrorValue(err) != nil) || !validateUser(user, login) {
        user = nil
    }
    return user
}

// Ensure that `user` has the given `login`.
//  NOTE: either panics or returns true
func validateUser(user *github.User, login string) bool {
    if user == nil {
        return false
    }
    if login != "" {
        if user.Login == nil {
            panic(fmt.Errorf("nil user.Login for %v", *user))
        } else if !strings.EqualFold(*user.Login, login) {
            panic(fmt.Errorf("login %q != user.Login %v", login, *user.Login))
        }
    }
    return true
}

// ============================================================================
// Internal functions - rendering
// ============================================================================

const githubUserFieldCount = 41

// Render a GitHub user object as a multiline string.
func userDetails(u *github.User) string {
    if u == nil { return "" }
    res := make([]string, 0, githubUserFieldCount)
    wid := util.CharCount("TwoFactorAuthentication")
    add := func(key string, val any) {
        res = append(res, fmt.Sprintf("%-*s %v", wid, key, val))
    }

    if u.Login                      != nil { add("Login",                   *u.Login) }
    if u.ID                         != nil { add("ID",                      *u.ID) }
    if u.NodeID                     != nil { add("NodeID",                  *u.NodeID) }
    if u.AvatarURL                  != nil { add("AvatarURL",               *u.AvatarURL) }
    if u.HTMLURL                    != nil { add("HTMLURL",                 *u.HTMLURL) }
    if u.GravatarID                 != nil { add("GravatarID",              *u.GravatarID) }
    if u.Name                       != nil { add("Name",                    *u.Name) }
    if u.Company                    != nil { add("Company",                 *u.Company) }
    if u.Blog                       != nil { add("Blog",                    *u.Blog) }
    if u.Location                   != nil { add("Location",                *u.Location) }
    if u.Email                      != nil { add("Email",                   *u.Email) }
    if u.Hireable                   != nil { add("Hireable",                *u.Hireable) }
    if u.Bio                        != nil { add("Bio",                     *u.Bio) }
    if u.TwitterUsername            != nil { add("TwitterUsername",         *u.TwitterUsername) }
    if u.PublicRepos                != nil { add("PublicRepos",             *u.PublicRepos) }
    if u.PublicGists                != nil { add("PublicGists",             *u.PublicGists) }
    if u.Followers                  != nil { add("Followers",               *u.Followers) }
    if u.Following                  != nil { add("Following",               *u.Following) }
    if u.CreatedAt                  != nil { add("CreatedAt",               *u.CreatedAt) }
    if u.UpdatedAt                  != nil { add("UpdatedAt",               *u.UpdatedAt) }
    if u.SuspendedAt                != nil { add("SuspendedAt",             *u.SuspendedAt) }
    if u.Type                       != nil { add("Type",                    *u.Type) }
    if u.SiteAdmin                  != nil { add("SiteAdmin",               *u.SiteAdmin) }
    if u.TotalPrivateRepos          != nil { add("TotalPrivateRepos",       *u.TotalPrivateRepos) }
    if u.OwnedPrivateRepos          != nil { add("OwnedPrivateRepos",       *u.OwnedPrivateRepos) }
    if u.PrivateGists               != nil { add("PrivateGists",            *u.PrivateGists) }
    if u.DiskUsage                  != nil { add("DiskUsage",               *u.DiskUsage) }
    if u.Collaborators              != nil { add("Collaborators",           *u.Collaborators) }
    if u.TwoFactorAuthentication    != nil { add("TwoFactorAuthentication", *u.TwoFactorAuthentication) }
    if u.Plan                       != nil { add("Plan",                    *u.Plan) }
    if u.LdapDn                     != nil { add("LdapDn",                  *u.LdapDn) }
    if u.URL                        != nil { add("URL",                     *u.URL) }
    if u.EventsURL                  != nil { add("EventsURL",               *u.EventsURL) }
    if u.FollowingURL               != nil { add("FollowingURL",            *u.FollowingURL) }
    if u.FollowersURL               != nil { add("FollowersURL",            *u.FollowersURL) }
    if u.GistsURL                   != nil { add("GistsURL",                *u.GistsURL) }
    if u.OrganizationsURL           != nil { add("OrganizationsURL",        *u.OrganizationsURL) }
    if u.ReceivedEventsURL          != nil { add("ReceivedEventsURL",       *u.ReceivedEventsURL) }
    if u.ReposURL                   != nil { add("ReposURL",                *u.ReposURL) }
    if u.StarredURL                 != nil { add("StarredURL",              *u.StarredURL) }
    if u.SubscriptionsURL           != nil { add("SubscriptionsURL",        *u.SubscriptionsURL) }

    return strings.Join(res, "\n")
}
