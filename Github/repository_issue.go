// Github/repository_issue.go

package Github

import (
	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Internal functions
// ============================================================================

// Generate a GitHub options object for ListByRepo().
//
//  Milestone   "none", "*", milestone number
//  State       "open" (def), "closed", "all"
//  Assignee    "none", "*", user name
//  Creator     Filter based on creator.
//  Mentioned   Filter based on issues mentioned by a specific user.
//  Labels      Filter based on label(s).
//  Sort        "created" (default), "updated", "comments"
//  Direction   "asc", "desc" (def)
//  Since       Filter issues by time.
//  Page        The page of results to receive.
//  PerPage     Results to include per page.
//
func issueListByRepoOptions() *github.IssueListByRepoOptions {
    opt := github.IssueListByRepoOptions{}
    opt.Page    = 1
    opt.PerPage = MAX_PER_PAGE
    return &opt
}

// Get all GitHub Issue objects for the indicated repository.
//  NOTE: if an error was encountered a partial list may be returned
func getRepoIssues(client *github.Client, owner, repo string) ([]*github.Issue, error) {
    res := []*github.Issue{}
    opt := issueListByRepoOptions()
    fn  := util.FuncName()
    for opt.Page > 0 {
        list, rsp, err := client.Issues.ListByRepo(ctx, owner, repo, opt)
        if err != nil {
            return res, log.ErrorValueIn(fn, err)
        }
        res = append(res, list...)
        opt.Page = rsp.NextPage
    }
    return res, nil
}
