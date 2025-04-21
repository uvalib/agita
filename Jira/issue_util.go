// Jira/issue_util.go

package Jira

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

// A generic issue reference.
type IssueArg interface {
    Issue | jira.Issue | *Issue | *jira.Issue
}

// ============================================================================
// Exported functions
// ============================================================================

// Transform a generic issue reference.
func NormalizeIssue[T IssueArg](arg T) *jira.Issue {
    switch v := any(arg).(type) {
        case Issue:       return v.ptr
        case jira.Issue:  return &v
        case *Issue:      if v != nil { return v.ptr }
        case *jira.Issue: return v
        default:          panic(fmt.Errorf("unexpected: %v", v))
    }
    return nil
}

// Return the key for the given issue.
func IssueNumber[T IssueArg](arg T) IssueKey {
    if i := NormalizeIssue(arg); i == nil {
        return ""
    } else {
        return i.Key
    }
}

// Return the text for the given issue.
func IssueText[T IssueArg](arg T) string {
    if i := NormalizeIssue(arg); (i == nil) || (i.Fields == nil) {
        return ""
    } else {
        return i.Fields.Description
    }
}

// Indicate whether the two issue objects refer to the same issue.
func SameIssue[T1, T2 IssueArg](arg1 T1, arg2 T2) bool {
    if nil1, nil2 := util.IsNil(arg1), util.IsNil(arg2); nil1 || nil2 {
        return nil1 && nil2
    } else if num1, num2 := IssueNumber(arg1), IssueNumber(arg2); (num1 != "") || (num2 != "") {
        return num1 == num2
    } else {
        return IssueText(arg1) == IssueText(arg2)
    }
}
