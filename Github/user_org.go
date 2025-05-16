// Github/user_org.go

package Github

import (
	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Internal functions
// ============================================================================

// Generate a GitHub options object for List().
//
//  Page    The page of results to receive.
//  PerPage Results to include per page.
//
func listOptions() *github.ListOptions {
    opt := github.ListOptions{}
    opt.Page    = 1
    opt.PerPage = MAX_PER_PAGE
    return &opt
}

// Get all GitHub Organization objects for user.
//  NOTE: may return partial results on error
func getUserOrgs(client *github.Client, user string) ([]*github.Organization, error) {
    res := []*github.Organization{}
    opt := listOptions()
    fn  := util.FuncName()
    for opt.Page > 0 {
        list, rsp, err := client.Organizations.List(ctx, user, opt)
        extractRateLimit(rsp)
        if err != nil {
            return res, log.ErrorValueIn(fn, err)
        }
        res = append(res, list...)
        opt.Page = rsp.NextPage
    }
    return res, nil
}
