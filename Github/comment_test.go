// Github/comment_test.go

package Github

import (
	"fmt"
	"testing"

	"lib.virginia.edu/agita/test"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Tests - Exported functions
// ============================================================================

func TestNewCommentType(t *testing.T) {
	const fn = "NewCommentType"

	type args struct {
		client  *Client
		repo    *Repository
		comment *github.IssueComment
	}
	type testCase struct {
		name string
		args args
		want *Comment
		err  string
	}

	client := TestClient
	repo   := SampleRepository(client)
	Case   := func(idx int, text, err string) (tc testCase) {
		tc.name = test.CaseName(fn, idx)
		tc.args = args{client, repo, nil}
		tc.err  = err
		if err == "" {
			body := test.Unique(text, tc.name)
			com  := testGithubComment(body)
			tc.want = &Comment{ptr: com, client: client}
			tc.args.comment = com
		}
		return
	}

	body  := FAKE_COMMENT_BODY
	tests := []testCase{
		Case(0, "",   ERR_NIL_ISSUE_COMMENT),
		Case(1, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer test.EvaluatePanic(tt.name, tt.err, t)
			got := NewCommentType(tt.args.client, tt.args.repo, tt.args.comment)
            testVerifyComment(fn, got, tt.want, t)
		})
	}
}

func TestGetComment(t *testing.T) {
	const fn = "GetComment"

	type args struct {
		client *Client
		owner  string
		repo   string
		id     int64
	}
	type testCase struct {
		name string
		args args
		want *Comment
		err  string
	}

	client := TestClient
    owner  := SAMPLE_ORG
    repo   := SAMPLE_REPO
	Case   := func(idx int, id int64, text, err string) (tc testCase) {
		tc.name = test.CaseName(fn, idx)
		tc.args = args{client, owner, repo, id}
		tc.err  = err
		if id != 0 {
			com := testGithubComment(text)
			tc.want = &Comment{ptr: com, client: client}
		}
		return
	}

    com   := SAMPLE_COMMENT
	body  := SAMPLE_COMMENT_MAP[com].Body
	tests := []testCase{
		Case(0, 0,   "",   ERR_NOT_FOUND),
		Case(1, com, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer test.EvaluatePanic(tt.name, tt.err, t)
			got := GetComment(tt.args.client, tt.args.owner, tt.args.repo, tt.args.id)
            testVerifyComment(fn, got, tt.want, t)
		})
	}
}

func TestCreateComment(t *testing.T) {
	const fn = "CreateComment"
	if test.Passive(fn, t) { return }

	type args struct {
		client *Client
		owner  string
		repo   string
		issue  int
		src    *github.IssueComment
	}
	type testCase struct {
		name string
		args args
		want *Comment
		err  string
	}

	client := TestClient
	repo   := GetFakeRepo(client)
	issue  := GetFakeIssue(repo)
	Case   := func(idx int, text, err string) (tc testCase) {
		tc.name = test.CaseName(fn, idx)
		tc.args = args{client, repo.Owner, repo.Name, issue.Number, nil}
		tc.err  = err
		if err == "" {
			body := test.Unique(text, tc.name)
			com  := testGithubComment(body)
			tc.want = &Comment{ptr: com, client: client}
			tc.args.src = com
		}
		return
	}

	body  := FAKE_COMMENT_BODY
	tests := []testCase{
		Case(0, "",   ERR_NO_ISSUE_COMMENT),
		Case(1, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer test.EvaluatePanic(tt.name, tt.err, t)
			got := CreateComment(tt.args.client, tt.args.owner, tt.args.repo, tt.args.issue, tt.args.src)
            testVerifyComment(fn, got, tt.want, t)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	const fn = "DeleteComment"
	if test.Passive(fn, t) { return }

	type args struct {
		client    *Client
		owner     string
		repo      string
		commentId int64
	}
	type testCase struct {
		name string
		args args
		err  string
	}

	client := TestClient
	repo   := GetFakeRepo(client)
	Case   := func(idx int, comment *Comment, err string) (tc testCase) {
		tc.name = test.CaseName(fn, idx)
		tc.args = args{client, repo.Owner, repo.Name, 0}
		tc.err  = err
        if comment != nil {
            tc.args.commentId = comment.ID()
        }
		return
	}

	issue  := GetFakeIssue(repo)
    com    := GetFakeComment(issue)
    tests  := []testCase{
        Case(0, nil, ERR_NOT_FOUND),
        Case(0, com, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer test.EvaluatePanic(tt.name, tt.err, t)
			DeleteComment(tt.args.client, tt.args.owner, tt.args.repo, tt.args.commentId)
		})
	}
}

// ============================================================================
// Tests - Exported methods - rendering
// ============================================================================

func TestComment_Details(t *testing.T) {
	const fn = "Comment.Details"

	type fields struct {
		ptr    *github.IssueComment
		repo   *Repository
		client *Client
	}
	type testCase struct {
		name   string
		fields fields
		want   string
		err    string
	}

	client := TestClient
	repo   := SampleRepository(client)
	Case   := func(idx int, text, err string) (tc testCase) {
		tc.name   = test.CaseName(fn, idx)
		tc.fields = fields{nil, repo, client}
		tc.err    = err
		if text == "" {
			tc.fields.ptr = &github.IssueComment{}
		} else {
			body := test.Unique(text, tc.name)
			tc.fields.ptr = testGithubComment(body)
			tc.want = fmt.Sprintf("Body              %s", body)
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
			com := NewCommentType(tt.fields.client, tt.fields.repo, tt.fields.ptr)
			if got, want := com.Details(), tt.want; got != want {
				t.Errorf("%s() = %q, want %q", fn, got, want)
			}
		})
	}
}

// ============================================================================
// Test setup
// ============================================================================

// Initialize variables related to testing GitHub comments.
func testSetup_Comment() {
	// no-op
}

// Clean up variables related to testing GitHub comments.
func testTeardown_Comment() {
	// NOTE: FakeComments are deleted by deleting FakeRepositories.
	//PreserveFakes()
}
