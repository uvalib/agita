// Jira/issue_func_test.go

package Jira

import (
	"testing"

	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/test"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Internal functions
// ============================================================================

// Generate a minimal Issue object.
func testJiraIssue(key IssueKey) *jira.Issue {
    return &jira.Issue{Key: key, Fields: &jira.IssueFields{}}
}

// Generate a minimal Issue object.
func testIssue(key IssueKey) *Issue {
    return &Issue{ptr: testJiraIssue(key)}
}

// A list of issue requirements suitable for comparison against a list of real
// issue objects.
func testWantIssues(keys ...IssueKey) []Issue {
    result := make([]Issue, 0, len(keys))
    for _, key := range keys {
        result = append(result, *testIssue(key))
    }
    return result
}

// ============================================================================
// Internal functions - verification
// ============================================================================

// Fail if the GitHub issue doesn't match the provided criteria.
func testVerifyIssue(fn string, got, want *Issue, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if w := IssueNumber(want); w != "" {
        if g := IssueNumber(got); g != w {
            t.Errorf("%s().Key = %q, want %q", fn, g, w)
        }
    }
    if w := IssueText(want); w != "" {
        if g, simple := IssueText(got), re.IsSimple(w); simple && (g != w) {
            t.Errorf("%s().Fields.Description = %q, want %q", fn, g, w)
        } else if !simple && !re.New(w).Match(g) {
            t.Errorf("%s().Fields.Description = %q, want match of `%s`", fn, g, w)
        }
    }
}

// Fail if the GitHub issues don't match the provided list.
//  NOTE: extra `got` items are ignored unless `want` is empty
func testVerifyIssues(fn string, got, want []Issue, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if len1, len2 := len(got), len(want); len1 < len2 {
        t.Errorf("%s() = %v issues, want %v issues", fn, len1, len2)
        for _, issue := range got {
            t.Errorf("%s() got %v %q", fn, IssueNumber(issue), IssueText(issue))
        }
        for _, issue := range want {
            t.Errorf("%s() want %v %q", fn, IssueNumber(issue), IssueText(issue))
        }
    } else if len2 == 0 {
        if len1 > len2 {
            t.Errorf("%s() = %v issues, want no issues", fn, len1)
        }
    } else {
        g := map[string]*Issue {}
        for _, issue := range got {
            if id := IssueNumber(issue); id != "" {
                g[id] = &issue
            }
        }
        for _, issue := range want {
            if id := IssueNumber(issue); id != "" {
                if g[id] == nil {
                    t.Errorf("%s() missing issue %q", fn, id)
                }
            }
        }
    }
}
