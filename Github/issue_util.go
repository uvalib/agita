// Github/issue_util.go

package Github

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// A generic issue reference.
type IssueArg interface {
    Issue | github.Issue | *Issue | *github.Issue
}

// ============================================================================
// Exported functions
// ============================================================================

// Transform a generic issue reference.
func NormalizeIssue[T IssueArg](arg T) *github.Issue {
    switch v := any(arg).(type) {
        case Issue:         return v.ptr
        case github.Issue:  return &v
        case *Issue:        if v != nil { return v.ptr }
        case *github.Issue: return v
        default:            panic(fmt.Errorf("unexpected: %v", v))
    }
    return nil
}

// Return the number for the given issue.
func IssueNumber[T IssueArg](arg T) int {
    switch v := any(arg).(type) {
        case Issue:  if v.Number != 0 { return v.Number }
        case *Issue: if (v != nil) && (v.Number != 0) { return v.Number }
    }
    if iss := NormalizeIssue(arg); (iss == nil) || (iss.Number == nil) {
        return 0
    } else {
        return *iss.Number
    }
}

// Return the name for the given issue.
func IssueTitle[T IssueArg](arg T) string {
    if iss := NormalizeIssue(arg); (iss == nil) || (iss.Title == nil) {
        return ""
    } else {
        return *iss.Title
    }
}

// Return the text for the given issue.
func IssueBody[T IssueArg](arg T) string {
    if iss := NormalizeIssue(arg); (iss == nil) || (iss.Body == nil) {
        return ""
    } else {
        return *iss.Body
    }
}

// Indicate whether the two issue objects refer to the same issue.
func SameIssue[T1, T2 IssueArg](arg1 T1, arg2 T2) bool {
    if nil1, nil2 := util.IsNil(arg1), util.IsNil(arg2); nil1 || nil2 {
        return nil1 && nil2
    } else if num1, num2 := IssueNumber(arg1), IssueNumber(arg2); (num1 != 0) || (num2 != 0) {
        return num1 == num2
    } else if tit1, tit2 := IssueTitle(arg1), IssueTitle(arg2); (tit1 != "") || (tit2 != "") {
        return tit1 == tit2
    } else {
        return IssueBody(arg1) == IssueBody(arg2)
    }
}
