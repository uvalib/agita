// Github/comment_fake_test.go
//
// Fake comments for tests.

package Github

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported constants
// ============================================================================

// Temporary fake issue property.
const (
    FAKE_COMMENT_USER = "rwl"
    FAKE_COMMENT_BODY = "fake comment"
)

// ============================================================================
// Exported functions - fakes
// ============================================================================

// Get a comment from a fake issue, creating one if necessary.
func GetFakeComment(issue *Issue) *Comment {
    if issue == nil { panic(ERR_NIL_ISSUE) }
    comments := issue.Comments()
    if len(comments) == 0 {
        return CreateFakeComment(issue)
    } else {
        return comments[0]
    }
}

// Create a new fake repository.
func CreateFakeComment(issue *Issue) *Comment {
    if issue == nil { panic(ERR_NIL_ISSUE) }
    return issue.CreateComment(fakeIssueComment())
}

// ============================================================================
// Internal functions
// ============================================================================

// Generate an IssueComment for testing.
func fakeIssueComment() *github.IssueComment {
    user := FAKE_COMMENT_USER
    body := util.Randomize(FAKE_COMMENT_BODY)
    if !CAN_SET_USER {
        tag := "ORIGINAL JIRA COMMENT" // see convert.commentAnnotations()
        body = fmt.Sprintf("%s Author = %q\n\n%s", tag, user, body)
    }
    react := github.Reactions{PlusOne: github.Ptr(18), Heart: github.Ptr(9)}
    return &github.IssueComment{
        Body:       &body,
        User:       &github.User{Login: github.Ptr(user)},
        Reactions:  &react, // NOTE: this doesn't get set
        CreatedAt:  nil,    // NOTE: can't set
        UpdatedAt:  nil,    // NOTE: can't set
    }
}
