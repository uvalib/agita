// Github/org.go
//
// Functions supporting GitHub organizations.

package Github

import (
	"fmt"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// A generic repository reference.
type OrgArg interface {
    github.Organization | *github.Organization
}

// ============================================================================
// Exported functions
// ============================================================================

// For use in arguments as the repository owner, defaulting to ORG.
func OrgOwner(owner string) string {
    if owner == ""  {
        owner = ORG
        log.Warn("using default org %q", owner)
    }
    return owner
}

// ============================================================================
// Internal functions
// ============================================================================

// Generate a GitHub options object for ListByOrg().
//
//  Type        "all" (def), "public", "private", "forks", "sources", "member"
//  Sort        "created" (def), "updated", "pushed", "full_name"
//  Direction   "asc" (def when using full_name), "desc" (def otherwise)
//  Page        The page of results to receive.
//  PerPage     Results to include per page.
//
func repoListByOrgOptions() *github.RepositoryListByOrgOptions {
    opt := github.RepositoryListByOrgOptions{}
    opt.Sort    = "full_name"
    opt.Page    = 1
    opt.PerPage = MAX_PER_PAGE
    return &opt
}

// Get all GitHub Repository objects for org.
//  NOTE: Unlike getUserRepos() this *does* include private repos.
//  NOTE: if `org` is blank it defaults to ORG
func getOrgRepos(client *github.Client, org string) ([]*github.Repository, error) {
    res := []*github.Repository{}
    opt := repoListByOrgOptions()
    fn  := util.FuncName()
    org = OrgOwner(org)
    for opt.Page > 0 {
        list, rsp, err := client.Repositories.ListByOrg(ctx, org, opt)
        if err != nil {
            return res, log.ErrorValueIn(fn, err)
        }
        res = append(res, list...)
        opt.Page = rsp.NextPage
    }
    return res, nil
}

// ============================================================================
// Internal functions - reporting
// ============================================================================

// Show all repositories for org.
func printOrgRepos(client *github.Client, org string) {
    if repos, err := getOrgRepos(client, org); err == nil {
        fmt.Printf("\n*** all repos for org %q (%d)\n", org, len(repos))
        util.PrintSortedMap(RepoMap(repos))
    }
}
