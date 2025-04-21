// Github/issue_print.go
//
// Reporting on GitHub issues.

package Github

import (
	"fmt"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported members
// ============================================================================

// Show details about the instance.
func (i *Issue) Print() {
    number, title := 0, MISSING
    if (i.ptr != nil) && (i.ptr.Number != nil) { number = *i.ptr.Number }
    if (i.ptr != nil) && (i.ptr.Title  != nil) { title  = *i.ptr.Title }
    fmt.Printf("\n*** GITHUB Issue %d %q:\n", number, title)
    fmt.Println(i.Details())
}

// Show all comments for the issue.
func (i *Issue) PrintComments() {
    cli   := i.client
    owner := i.repo.Owner
    repo  := i.repo.Name
    issue := i.Number
    printIssueComments(cli.ptr, owner, repo, issue)
}

// ============================================================================
// Internal functions
// ============================================================================

// Show all comments associated with an issue.
func printIssueComments(client *github.Client, owner, repo string, issue int) {
    if items, err := getIssueComments(client, owner, repo, issue); err == nil {
        tag := fmt.Sprintf("\n*** %s/%s issue %d comments", owner, repo, issue)
        fmt.Printf("%s (%d)\n", tag, len(items))
        for _, comment := range items {
            id    := *comment.ID
            entry := commentDetails(comment)
            fmt.Printf("%s - comment %d:\n%s", tag, id, entry)
        }
    }
}
