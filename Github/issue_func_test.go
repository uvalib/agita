// Github/issue_func_test.go

package Github

import (
	"testing"
	"time"

	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/test"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Internal functions
// ============================================================================

// Generate a minimal issue object.
func testIssue(text string, number int) *Issue {
    ptr := testGithubIssue(text, number)
    return &Issue{Number: *ptr.Number, ptr: ptr}
}

// Generate a minimal issue object.
func testGithubIssue(text string, number int) *github.Issue {
    return &github.Issue{Number: github.Ptr(number), Body: github.Ptr(text)}
}

// Generate a minimal issue import object.
func testIssueImport(title, body string) *IssueImport {
    return &IssueImport{github.IssueImport{Title: title, Body: body}}
}

// Generate a minimal issue request object.
func testIssueRequest(title, body string) *github.IssueRequest {
    pTitle := github.Ptr(title)
    pBody  := github.Ptr(body)
    return &github.IssueRequest{Title: pTitle, Body: pBody}
}

// A list of issue requirements suitable for comparison against a list of real
// issue objects.
func testWantIssues(ids ...int) []*Issue {
    result := make([]*Issue, 0, len(ids))
    for _, id := range ids {
        result = append(result, &Issue{Number: id})
    }
    return result
}

// ============================================================================
// Internal functions - verification
// ============================================================================

// Fail if the GitHub issue import didn't complete successfully.
func testVerifyIssueImport(fn string, repo *Repository, id int, t *testing.T) {
    if id == 0 { return }
    done, status := false, MISSING
    for i := 0; !done && (i < 3); i++ {
        time.Sleep(time.Second)
        done, status = checkImportIssue(repo.client.ptr, repo.Owner, repo.Name, id)
    }
    if !done {
        t.Errorf("%s() did not complete; status = %q", fn, status)
    } else if status != "imported" {
        t.Errorf("%s() completed with status = %q", fn, status)
    }
}

// Fail if the GitHub issue doesn't match the provided criteria.
func testVerifyIssue(fn string, got, want *Issue, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if w := IssueNumber(want); w != 0 {
        if g := IssueNumber(got); g != w {
            t.Errorf("%s().Number = %d, want %d", fn, g, w)
        }
    }
    if w := IssueTitle(want); w != "" {
        if g, simple := IssueTitle(got), re.IsSimple(w); simple && (g != w) {
            t.Errorf("%s().Title = %q, want %q", fn, g, w)
        } else if !simple && !re.New(w).Match(g) {
            t.Errorf("%s().Title = %q, want match of `%s`", fn, g, w)
        }
    }
    if w := IssueBody(want); w != "" {
        if g, simple := IssueBody(got), re.IsSimple(w); simple && (g != w) {
            t.Errorf("%s().Body = %q, want %q", fn, g, w)
        } else if !simple && !re.New(w).Match(g) {
            t.Errorf("%s().Body = %q, want match of `%s`", fn, g, w)
        }
    }
}

// Fail if the GitHub issues don't match the provided list.
//  NOTE: extra `got` items are ignored unless `want` is empty
func testVerifyIssues(fn string, got, want []*Issue, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if len1, len2 := len(got), len(want); len1 < len2 {
        t.Errorf("%s() = %v issues, want %v issues", fn, len1, len2)
        for _, issue := range got {
            t.Errorf("%s() got %v %q", fn, issue.Number, issue.Body())
        }
        for _, issue := range want {
            t.Errorf("%s() want %v %q", fn, issue.Number, issue.Body())
        }
    } else if len2 == 0 {
        if len1 > len2 {
            t.Errorf("%s() = %v issues, want no issues", fn, len1)
        }
    } else {
        g := map[int]*Issue {}
        for _, issue := range got {
            if id := IssueNumber(issue); id != 0 {
                g[id] = issue
            }
        }
        for _, issue := range want {
            if id := IssueNumber(issue); id != 0 {
                if g[id] == nil {
                    t.Errorf("%s() missing issue %d", fn, id)
                }
            }
        }
    }
}
