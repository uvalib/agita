// Github/repository_print.go
//
// Reporting on GitHub repository objects.

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
func (r *Repository) Print() {
    name := MISSING
    if (r.ptr != nil) && (r.ptr.Name != nil) { name = *r.ptr.Name }
    fmt.Printf("\n*** GITHUB Repo %s:\n", name)
    fmt.Println(r.Details())
}

// Show all issues for the repository.
func (r *Repository) PrintIssues() {
    printRepoIssues(r.client.ptr, r.Owner, r.Name)
}

// Show all comments for all issues for the repository.
func (r *Repository) PrintComments() {
    printRepoComments(r.client.ptr, r.Owner, r.Name)
}

// ============================================================================
// Internal functions
// ============================================================================

// Show all issues for the indicated repository.
func printRepoIssues(client *github.Client, owner, repo string) {
    if items, err := getRepoIssues(client, owner, repo); err == nil {
        fmt.Printf("\n*** issues for %s/%s (%d)\n", owner, repo, len(items))
        util.PrintSortedMap(issueMap(items))
    }
}

// Show all comments for all issues for the indicated repository.
func printRepoComments(client *github.Client, owner, repo string) {
    if items, err := getRepoIssues(client, owner, repo); err == nil {
        tag := fmt.Sprintf("\n*** %s/%s issue comments", owner, repo)
        fmt.Printf("%s (%d issues)\n", tag, len(items))
        for _, issue := range items {
            number := *issue.Number
            if coms, err := getIssueComments(client, owner, repo, number); err == nil {
                innerTag := fmt.Sprintf("%s - issue %d", tag, number)
                for _, comment := range coms {
                    id    := *comment.ID
                    entry := commentDetails(comment)
                    fmt.Printf("%s - comment %d:\n%s\n", innerTag, id, entry)
                }
            }
        }
    }
}
