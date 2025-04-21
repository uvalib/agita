// Github/issue_fake_test.go
//
// Fake issues for tests.

package Github

import (
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported functions
// ============================================================================

// Get an issue from a fake repository, creating one if necessary.
func GetFakeIssue(repo *Repository) *Issue {
    if repo == nil { panic(ERR_NIL_REPO) }
    issues := repo.GetIssues()
    if len(issues) == 0 {
        return CreateFakeIssue(repo)
    } else {
        return issues[0]
    }
}

// Create a new fake repository.
func CreateFakeIssue(repo *Repository) *Issue {
    if repo == nil { panic(ERR_NIL_REPO) }
    return repo.CreateIssue(fakeIssueRequest())
}

// ============================================================================
// Exported constants
// ============================================================================

// Temporary fake issue property.
const (
    FAKE_ISSUE_TITLE = "fake GitHub issue"
    FAKE_ISSUE_BODY  = "fake GitHub issue body"
)

// ============================================================================
// Internal constants
// ============================================================================

// IssueRequest.State value.
const (
    stateOpen   = "open"
    stateClosed = "closed"
)

// IssueRequest.StateReason value.
const (
    stateReasonCompleted  = "completed"
    stateReasonReopened   = "reopened"
    stateReasonNotPlanned = "not_planned"
)

// ============================================================================
// Internal functions
// ============================================================================

// Generate an IssueRequest for testing.
func fakeIssueRequest() (result *github.IssueRequest) {
    //assignees := []string{"UVADave"}
    assignees := []string{SAMPLE_USER}

    //state := stateOpen
    state := stateClosed // stateOpen

    //stateReason := stateReasonCompleted
    //stateReason := stateReasonReopened
    stateReason := stateReasonNotPlanned

    result = &github.IssueRequest{
        Title:          github.Ptr(util.Randomize(FAKE_ISSUE_TITLE)),
        Body:           github.Ptr(util.Randomize(FAKE_ISSUE_BODY)),
        Labels:         github.Ptr([]string{"label1","label2"}),
        State:          github.Ptr(state),
        StateReason:    github.Ptr(stateReason),
        //Milestone:    github.Ptr(1),
    }
    switch len(assignees) {
        case 0:  // skip
        case 1:  result.Assignee  = github.Ptr(assignees[0])
        default: result.Assignees = github.Ptr(assignees)
    }
    return
}

