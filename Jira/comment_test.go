// Jira/comment.go
//
// Application Comment type backed by a jira.Comment object.

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

func TestNewCommentType(t *testing.T) {
    const fn = "NewCommentType"

	type args struct {
		client  *Client
		comment *jira.Comment
	}
    type testCase struct {
		name string
		args args
		want *Comment
        err  string
	}

    client := TestClient
    Case   := func(idx int, text, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client: client}
        tc.err  = err
        if err == "" {
            body := test.Unique(text, tc.name)
            com  := testJiraComment(body)
            tc.args.comment = com
            tc.want = &Comment{ptr: com, client: client}
        }
        return
    }

    body  := FAKE_COMMENT_BODY
	tests := []testCase{
        Case(0, "",   ERR_NIL_COMMENT),
        Case(1, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := NewCommentType(tt.args.client, tt.args.comment)
            testVerifyComment(fn, got, tt.want, t)
		})
	}
}

func TestGetCommentById(t *testing.T) {
    const fn = "GetCommentById"

	type args struct {
		client *Client
		issue  IssueKey
		id     CommentId
	}
    type testCase struct {
		name string
		args args
		want *Comment
        err  string
	}

    client := TestClient
    issue  := SAMPLE_ISSUE
    Case   := func(idx int, commentId int, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, issue, commentId}
        tc.err  = err
        if commentId != 0 {
            tc.want = testComment(commentId)
        }
        return
    }

    com1  := SAMPLE_COMMENTS[0]
    com2  := SAMPLE_COMMENTS[1]
    tests := []testCase{
        Case(0, 0,    ERR_REQUEST_FAILED),
        Case(1, com1, ""),
        Case(2, com2, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := GetCommentById(tt.args.client, tt.args.issue, tt.args.id)
            testVerifyComment(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported methods - rendering
// ============================================================================

func TestComment_Details(t *testing.T) {
    const fn = "Comment.Details"

	type fields struct {
		ptr    *jira.Comment
		client *Client
	}
    type testCase struct{
		name   string
		fields fields
		want   string
        err    string
    }

	client := TestClient
	Case   := func(idx int, text, err string) (tc testCase) {
		tc.name   = test.CaseName(fn, idx)
		tc.fields = fields{nil, client}
		tc.err    = err
		if text == "" {
			tc.fields.ptr = &jira.Comment{}
		} else {
			body := test.Unique(text, tc.name)
			tc.fields.ptr = testJiraComment(body)
			tc.want = fmt.Sprintf("Body         %s", body)
		}
		return
	}

	body  := FAKE_COMMENT_BODY
    tests := []testCase {
		Case(0, "",    ""),
		Case(1, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			c := &Comment{
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
			if got, want := c.Details(), tt.want; got != want {
				t.Errorf("%s() = %q, want %q", fn, got, want)
			}
		})
	}
}

// ============================================================================
// Test setup
// ============================================================================

// Initialize variables related to testing Jira comments.
func testSetup_Comment() {
    // no op
}

// Clean up variables related to testing Jira comments.
func testTeardown_Comment() {
    // no op
}
