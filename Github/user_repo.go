// Github/user_repo.go

package Github

import (
	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Internal functions
// ============================================================================

// Generate a GitHub options object for ListByUser().
//
//  Type        "all" , "owner" (def), "member"
//  Sort        "created", "updated", "pushed", "full_name" (def)
//  Direction   "asc" (def when using full_name), "desc" (def otherwise)
//  Page        The page of results to receive.
//  PerPage     Results to include per page.
//
func repoListByUserOptions() *github.RepositoryListByUserOptions {
    opt := github.RepositoryListByUserOptions{}
    opt.Type    = "all"
    opt.Page    = 1
    opt.PerPage = MAX_PER_PAGE
    return &opt
}

// Get all GitHub Repository objects for user.
//  NOTE: Unlike getOrgRepos() this does *not* include private repos.
//  NOTE: may return partial results on error
func getUserRepos(client *github.Client, user string) ([]*github.Repository, error) {
    res := []*github.Repository{}
    opt := repoListByUserOptions()
    fn  := util.FuncName()
    for opt.Page > 0 {
        list, rsp, err := client.Repositories.ListByUser(ctx, user, opt)
        if err != nil {
            return res, log.ErrorValueIn(fn, err)
        }
        res = append(res, list...)
        opt.Page = rsp.NextPage
    }
    return res, nil
}
