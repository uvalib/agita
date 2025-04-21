// Github/repository_ops.go
//
// Operations on github.Repository objects.

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

// Transform a list of GitHub Repository objects into a table of repository
// names and descriptions.
func RepoMap(repos []*github.Repository) map[string]string {
    result := make(map[string]string, len(repos))
    for _, repo := range repos {
        name, desc, priv := "NONE", "NONE", false
        if repo.Name        != nil { name = *repo.Name }
        if repo.Description != nil { desc = *repo.Description }
        if repo.Private     != nil { priv = *repo.Private }
        if priv                    { desc = "[private] " + desc }
        result[name] = desc
    }
    return result
}

// ============================================================================
// Internal functions
// ============================================================================

// Get a GitHub Repository object for the indicated repository.
//  NOTE: if `owner` is blank it defaults to ORG
//  NOTE: panics if the repository does not match `owner` and `name`.
//  NOTE: returns nil if `name` is blank
//  NOTE: returns nil on error
func getRepository(client *github.Client, owner, name string) *github.Repository {
    if name == "" {
        return nil
    }
    owner = OrgOwner(owner)
    repo, _, err := client.Repositories.Get(ctx, owner, name)
    if (log.ErrorValue(err) != nil) || !validateRepo(repo, owner, name) {
        return nil
    }
    return repo
}

// Ensure that `repo` has the given `owner` and `name`.
func validateRepo(repo *github.Repository, owner, name string) bool {
    if repo == nil {
        return false
    }
    if owner != "" {
        if repo.Owner == nil {
            panic(fmt.Errorf("nil repo.Owner for %v", *repo))
        } else if rOwner := Account(repo.Owner); rOwner != owner {
            panic(fmt.Errorf("owner %q != repo.Owner %v", owner, rOwner))
        }
    }
    if name != "" {
        if repo.Name == nil {
            panic(fmt.Errorf("nil repo.Name for %v", *repo))
        } else if rName := *repo.Name; !strings.EqualFold(rName, name) {
            panic(fmt.Errorf("name %q != repo.Name %v", name, rName))
        }
    }
    return true
}

// ============================================================================
// Internal functions
// ============================================================================

type RepositoryRequest struct {
    github.Repository
}

// Create a new GitHub repository with the provided properties.
//  NOTE: panics if `data.Name` is not present
//  NOTE: defaults `data.Owner` to ORG
func createRepository(client *github.Client, data *RepositoryRequest) *github.Repository {
    if data == nil {
        panic(ERR_NO_DATA)
    }
    name := ""
    if ptr := data.Name; (ptr == nil) || (*ptr == "") {
        panic(fmt.Errorf("%s in %v", ERR_NO_REPO_GIVEN, *data))
    } else {
        name = *ptr
    }
    owner := ""
    if ptr := data.Owner; (ptr == nil) || (ptr.Login == nil) {
        owner = ORG
        data.Owner = getUser(client, owner)
    } else {
        owner = *ptr.Login
    }
    repo, _, err := client.Repositories.Create(ctx, owner, &data.Repository)
    if (log.ErrorValue(err) != nil) || !validateRepo(repo, owner, name) {
        return nil
    }
    return repo
}

// Remove an existing GitHub repository.
//  NOTE: panics if `data.Name` is not present
//  NOTE: defaults `data.Owner` to ORG
func deleteRepository(client *github.Client, owner, name string) {
    if name == "" {
        panic(ERR_NO_REPO_GIVEN)
    }
    owner = OrgOwner(owner)
    _, err := client.Repositories.Delete(ctx, owner, name)
    log.ErrorValue(err)
}

// ============================================================================
// Internal functions - rendering
// ============================================================================

const githubRepositoryFieldCount = 17 // 90

// Render a GitHub repository object as a multiline string.
func repoDetails(r *github.Repository) string {
    if r == nil { return "" }
    res := make([]string, 0, githubRepositoryFieldCount)
    wid := util.CharCount("OpenIssuesCount")
    add := func(key string, val any) {
        res = append(res, fmt.Sprintf("%-*s %v", wid, key, val))
    }

    if r.ID              != nil { add("ID",                 *r.ID) }
    if r.NodeID          != nil { add("NodeID",             *r.NodeID) }
    if r.Owner           != nil { add("Owner",               UserLabel(r.Owner)) }
    if r.Name            != nil { add("Name",               *r.Name) }
    if r.FullName        != nil { add("FullName",           *r.FullName) }
    if r.Description     != nil { add("Description",        *r.Description) }
    if r.Homepage        != nil { add("Homepage",           *r.Homepage) }
    if false                    { add("CodeOfConduct",      *r.CodeOfConduct)}
    if false                    { add("DefaultBranch",      *r.DefaultBranch)}
    if false                    { add("MasterBranch",       *r.MasterBranch)}
    if r.CreatedAt       != nil { add("CreatedAt",          *r.CreatedAt) }
    if r.PushedAt        != nil { add("PushedAt",           *r.PushedAt) }
    if r.UpdatedAt       != nil { add("UpdatedAt",          *r.UpdatedAt) }
    if r.HTMLURL         != nil { add("HTMLURL",            *r.HTMLURL) }
    if false                    { add("CloneURL",           *r.CloneURL)}
    if false                    { add("GitURL",             *r.GitURL)}
    if false                    { add("MirrorURL",          *r.MirrorURL)}
    if false                    { add("SSHURL",             *r.SSHURL)}
    if false                    { add("SVNURL",             *r.SVNURL)}
    if false                    { add("Language",           *r.Language)}
    if false                    { add("Fork",               *r.Fork)}
    if false                    { add("ForksCount",         *r.ForksCount)}
    if false                    { add("NetworkCount",       *r.NetworkCount)}
    if r.OpenIssuesCount != nil { add("OpenIssuesCount",    *r.OpenIssuesCount) }
    if r.OpenIssues      != nil { add("OpenIssues",         *r.OpenIssues) }
    if false                    { add("StargazersCount",    *r.StargazersCount)}
    if false                    { add("SubscribersCount",   *r.SubscribersCount)}
    if false                    { add("WatchersCount",      *r.WatchersCount)}
    if false                    { add("Watchers",           *r.Watchers)}
    if false                    { add("Size",               *r.Size)}
    if false                    { add("AutoInit",           *r.AutoInit)}
    if false                    { add("Parent",             *r.Parent)}
    if false                    { add("Source",             *r.Source)}
    if false                    { add("TemplateRepository", *r.TemplateRepository)}
    if false                    { add("Organization",       *r.Organization)}
    if r.Permissions     != nil { add("Permissions",         r.Permissions) }

    return strings.Join(res, "\n")
}

// ============================================================================
// Module initialization
// ============================================================================

// Initialize variables related to GitHub repositories.
func setupRepository() {
    // no-op
}
