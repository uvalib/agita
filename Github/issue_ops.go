// Github/issue_ops.go
//
// Operations on github.Issue objects.

package Github

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Internal functions
// ============================================================================

// On GitHub, create a new issue on the indicated repository.
//  NOTE: this is not sufficient for Jira issues:
//      - No way to specify Issue.User (the user who submitted the issue)
//      - No way to specify Issue.Type
func createIssue(client *github.Client, owner, repo string, req *github.IssueRequest) *github.Issue {
    if req == nil { panic(ERR_NO_ISSUE_REQUEST) }
    result, rsp, err := client.Issues.Create(ctx, owner, repo, req)
    extractRateLimit(rsp)
    if log.ErrorValue(err) == nil {
        log.Info("\n*** create issue %q - rsp = %v\n", *req.Title, rsp)
    }
    return result
}

// From GitHub, retrieve the indicated repository issue.
func getIssue(client *github.Client, owner, repo string, number int) *github.Issue {
    result, rsp, err := client.Issues.Get(ctx, owner, repo, number)
    extractRateLimit(rsp)
    log.ErrorValue(err)
    return result
}

// Remove the indicated repository issue from GitHub.
func deleteIssue(client *github.Client, owner, repo string, number int) bool {
    if issue := getIssue(client, owner, repo, number); issue != nil {
        if nodeId := issue.NodeID; nodeId != nil {
            return GqlDeleteIssue(*nodeId)
        }
    }
    return false
}

// ============================================================================
// Internal functions - rendering
// ============================================================================

const githubIssueFieldCount = 30

// Render a GitHub issue object as a multiline string.
func issueDetails(i *github.Issue) string {
    if i == nil { return "" }
    res := make([]string, 0, githubIssueFieldCount)
    max := util.CharCount("AuthorAssociation")
    add := func(key string, val any) {
        res = append(res, fmt.Sprintf("%-*s %v", max, key, val))
    }

    if i.Title             != nil { add("Title",                *i.Title) }
    if i.ID                != nil { add("ID",                   *i.ID) }
    if i.Number            != nil { add("Number",               *i.Number) }
    if i.State             != nil { add("State",                *i.State) }
    if i.StateReason       != nil { add("StateReason",          *i.StateReason) }
    if i.Locked            != nil { add("Locked",               *i.Locked) }
    if i.AuthorAssociation != nil { add("AuthorAssociation",    *i.AuthorAssociation) }
    if i.User              != nil { add("User",                 UserLabel(i.User)) }
    if i.Labels            != nil { add("Labels",               LabelStrings(i.Labels)) }
    if i.Assignee          != nil { add("Assignee",             UserLabel(i.Assignee)) }
    if i.Comments          != nil { add("Comments",             *i.Comments) }
    if i.ClosedAt          != nil { add("ClosedAt",             *i.ClosedAt) }
    if i.CreatedAt         != nil { add("CreatedAt",            *i.CreatedAt) }
    if i.UpdatedAt         != nil { add("UpdatedAt",            *i.UpdatedAt) }
    if i.ClosedBy          != nil { add("ClosedBy",             *i.ClosedBy) }
    if i.URL               != nil { add("URL",                  *i.URL) }
    if i.HTMLURL           != nil { add("HTMLURL",              *i.HTMLURL) }
    if i.CommentsURL       != nil { add("CommentsURL",          *i.CommentsURL) }
    if i.EventsURL         != nil { add("EventsURL",            *i.EventsURL) }
    if i.LabelsURL         != nil { add("LabelsURL",            *i.LabelsURL) }
    if i.RepositoryURL     != nil { add("RepositoryURL",        *i.RepositoryURL) }
    if i.Milestone         != nil { add("Milestone",            *i.Milestone) }
    if i.PullRequestLinks  != nil { add("PullRequestLinks",     *i.PullRequestLinks) }
    if i.Repository        != nil { add("Repository",           *i.Repository) }
    if i.Reactions         != nil { add("Reactions",            ReactionString(i.Reactions)) }
    if i.Assignees         != nil { add("Assignees",            UserLabels(i.Assignees)) }
    if i.NodeID            != nil { add("NodeID",               *i.NodeID) }
    if i.Draft             != nil { add("Draft",                *i.Draft) }
    if i.Type              != nil { add("Type",                 *i.Type) }
    if i.Body              != nil { add("Body",                 *i.Body) }

    return strings.Join(res, "\n")
}

// Convert a list of Issue objects into a table of issue titles and details.
func issueMap(issues []*github.Issue) map[string]string {
    result := make(map[string]string, len(issues))
    for _, issue := range issues {
        name  := *issue.Title
        entry := issueDetails(issue)
        result[name] = "\n" + entry
    }
    return result
}

// ============================================================================
// Module initialization
// ============================================================================

// Initialize variables related to GitHub issues.
func setupIssue() {
    // no-op
}
