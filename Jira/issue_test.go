// Jira/issue_test.go

package Jira

import (
	"fmt"
	"testing"

	"lib.virginia.edu/agita/test"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Tests - Exported functions
// ============================================================================

func TestNewIssueType(t *testing.T) {
    const fn = "NewIssueType"

	type args struct {
		client *Client
		issue  *jira.Issue
	}
    type testCase struct {
		name string
		args args
		want *Issue
        err  string
	}

    client := TestClient
    Case   := func(idx int, proj, key, title, body, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, nil}
        tc.err  = err
        if proj != "" {
            title = test.Unique(title, tc.name)
            body  = test.Unique(body,  tc.name)
            iss  := &jira.Issue{Key: key, Fields: &jira.IssueFields{Summary: title, Description: body}}
            tc.want = &Issue{iss, client}
            tc.args.issue = iss
        }
        return
    }

    proj  := SAMPLE_PROJ
    issue := SAMPLE_ISSUE
    title := FAKE_ISSUE_TITLE
    body  := FAKE_ISSUE_BODY
	tests := []testCase{
        Case(0, "",   "",    "",    "",   ERR_NIL_ISSUE),
        Case(1, proj, "",    "",    "",   ""),
        Case(2, proj, issue, title, "",   ""),
        Case(3, proj, issue, "",    body, ""),
        Case(4, proj, issue, title, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := NewIssueType(tt.args.client, tt.args.issue)
            testVerifyIssue(fn, got, tt.want, t)
		})
	}
}

func TestGetIssueByKey(t *testing.T) {
    const fn = "GetIssueByKey"

	type args struct {
		client *Client
		key    IssueKey
	}
    type testCase struct {
		name string
		args args
		want *Issue
        err  string
	}

    client := TestClient
    Case   := func(idx int, issueKey IssueKey, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, issueKey}
        tc.err  = err
        if issueKey != "" {
            tc.want = testIssue(issueKey)
        }
        return
    }

    tests := []testCase{
        Case(0, "", ERR_REQUEST_FAILED),
	}
    for i, key := range SAMPLE_ISSUES {
        tests = append(tests, Case(i+1, key, ""))
    }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := GetIssueByKey(tt.args.client, tt.args.key)
            testVerifyIssue(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported members - comments
// ============================================================================

func TestIssue_Comments(t *testing.T) {
    const fn = "Issue.Comments"

	type fields struct {
		ptr    *jira.Issue
		client *Client
	}
    type testCase struct {
		name   string
		fields fields
		want   []Comment
        err    string
	}

    client := TestClient
    Case   := func(idx int, issueKey string, want []Comment, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{client: client}
        tc.want   = want
        tc.err    = err
        if issueKey == "" {
            tc.fields.ptr = &jira.Issue{}
        } else {
            tc.fields.ptr = &jira.Issue{Key: issueKey}
        }
        return
    }

    issue := SAMPLE_ISSUE
    coms0 := testWantComments()
    coms1 := testWantComments(SAMPLE_COMMENTS...)
    tests := []testCase{
        Case(0, "",    coms0, ERR_NO_ISSUE),
        Case(1, issue, coms1, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			i := &Issue{
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            got := i.Comments()
            testVerifyComments(fn, got, tt.want, t)
		})
	}
}

func TestIssue_Comment(t *testing.T) {
    const fn = "Issue.Comment"

	type fields struct {
		ptr    *jira.Issue
		client *Client
	}
	type args struct {
		id CommentId
	}
    type testCase struct {
		name   string
		fields fields
		args   args
		want   *Comment
        err    string
	}

    client := TestClient
    issue  := GetIssueByKey(client, SAMPLE_ISSUE)
    Case   := func(idx int, id CommentId, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{issue.ptr, client}
        tc.args   = args{id}
        tc.err    = err
        if id != 0 {
            tc.want = testComment(id)
        }
        return
    }

    tests := []testCase{
        Case(0, 0, ERR_REQUEST_FAILED),
	}
    for i, id := range SAMPLE_COMMENTS {
        tests = append(tests, Case(i+1, id, ""))
    }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			i := &Issue{
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            got := i.Comment(tt.args.id)
            testVerifyComment(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported members - rendering
// ============================================================================

func TestIssue_Details(t *testing.T) {
    const fn = "Issue.Details"

	type fields struct {
		ptr    *jira.Issue
		client *Client
	}
    type testCase struct{
		name   string
		fields fields
		want   string
        err    string
    }

    client := TestClient
    issue  := GetIssueByKey(client, SAMPLE_ISSUE)
    Case   := func(idx int, body, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{issue.ptr, client}
        tc.err    = err
        if body == "" {
            tc.fields.ptr = &jira.Issue{}
        } else {
            body = test.Unique(body, tc.name)
            tc.fields.ptr = &jira.Issue{Fields: &jira.IssueFields{Summary: body}}
            tc.want = fmt.Sprintf("Fields.Summary        %s", body)
        }
        return
    }

    body  := FAKE_COMMENT_BODY
    tests := []testCase{
        Case(0, "",   ""),
        Case(1, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			i := &Issue{
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
			if got, want := i.Details(), tt.want; got != want {
				t.Errorf("%s() = %q, want %q", fn, got, want)
			}
		})
	}
}

// ============================================================================
// Test setup
// ============================================================================

// Initialize variables related to testing Jira issues.
func testSetup_Issue() {
    // no op
}

// Clean up variables related to testing Jira issues.
func testTeardown_Issue() {
    // no op
}
