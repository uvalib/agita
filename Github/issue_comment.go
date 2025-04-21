// Github/issue_comment.go

package Github

import (
	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Internal functions
// ============================================================================

// Generate an options object for ListComments.
//
//  Sort        "created", "updated"
//  Direction   "asc", "desc"
//  Since       Filter issues by time.
//  Page        The page of results to receive.
//  PerPage     Results to include per page.
//
func issueListCommentsOptions() *github.IssueListCommentsOptions {
    opt := github.IssueListCommentsOptions{}
    opt.Page    = 1
    opt.PerPage = MAX_PER_PAGE
    return &opt
}

// Get all comments associated with an issue.
//  NOTE: an error may result in partial results
func getIssueComments(client *github.Client, owner, repo string, issue int) ([]*github.IssueComment, error) {
    res := []*github.IssueComment{}
    opt := issueListCommentsOptions()
    fn  := util.FuncName()
    for opt.Page > 0 {
        list, rsp, err := client.Issues.ListComments(ctx, owner, repo, issue, opt)
        if err != nil {
            return res, log.ErrorValueIn(fn, err)
        }
        res = append(res, list...)
        opt.Page = rsp.NextPage
    }
    return res, nil
}
