// Github/issue_test.go

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

func TestNewIssueType(t *testing.T) {
    const fn = "NewIssueType"

    type args struct {
		client *Client
		owner  string
		repo   string
		issue  *github.Issue
	}
    type testCase struct {
		name string
		args args
		want *Issue
        err  string
	}

    client := TestClient
    owner  := SAMPLE_ORG
    Case   := func(idx int, repo, title, body, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, owner, repo, nil}
        tc.err  = err
        if err == "" {
            body = test.Unique(body, tc.name)
            tc.want = testIssue(body, 0)
            tc.args.issue = testGithubIssue(body, 0)
        }
        return
    }

    repo  := SAMPLE_REPO
    title := FAKE_ISSUE_TITLE
    body  := FAKE_ISSUE_BODY
    tests := []testCase{
        Case(0, "",   "",    "",   ERR_NO_REPO_NAME),
        Case(1, repo, "",    "",   ERR_NIL_ISSUE),
        Case(2, repo, title, "",   ""),
        Case(3, repo, "",    body, ""),
        Case(4, repo, title, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := NewIssueType(tt.args.client, tt.args.owner, tt.args.repo, tt.args.issue)
            testVerifyIssue(fn, got, tt.want, t)
		})
	}
}

func TestImportIssue(t *testing.T) {
    const fn = "ImportIssue"
    if test.Passive(fn, t) { return }

    type args struct {
		client   *Client
		owner    string
		repo     string
		imp      *IssueImport
		comments []*CommentImport
	}
    type testCase struct {
		name string
		args args
        err  string
	}

    client := TestClient
    Case   := func(idx int, repo *Repository, title, body, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client: client, comments: []*CommentImport{}}
        tc.err  = err
        if repo != nil {
            title = test.Unique(title, tc.name)
            body  = test.Unique(body,  tc.name)
            tc.args.owner = repo.Owner
            tc.args.repo  = repo.Name
            tc.args.imp   = testIssueImport(title, body)
        }
        return
    }

    repo  := GetFakeRepo(client)
    title := FAKE_ISSUE_TITLE
    body  := FAKE_ISSUE_BODY
    tests := []testCase{
        Case(0, nil,  "",    "",   ERR_NO_ISSUE_IMPORT),
        Case(1, repo, "",    "",   ERR_INVALID_IMPORT),
        Case(2, repo, title, "",   ERR_INVALID_CREATE),
        Case(3, repo, "",    body, ERR_INVALID_CREATE),
        Case(4, repo, title, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            cli   := tt.args.client
            owner := tt.args.owner
            repo  := tt.args.repo
            if id := ImportIssue(cli, owner, repo, tt.args.imp, tt.args.comments...); id != 0 {
                r := &Repository{Owner: owner, Name: repo, client: cli}
                testVerifyIssueImport(fn, r, id, t)
            }
		})
	}
}

func TestCreateIssue(t *testing.T) {
    const fn = "CreateIssue"
    if test.Passive(fn, t) { return }

    type args struct {
		client *Client
		owner  string
		repo   string
		req    *github.IssueRequest
	}
    type testCase struct {
		name string
		args args
		want *Issue
        err  string
	}

    client := TestClient
    Case   := func(idx int, repo *Repository, title, body, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client: client}
        tc.err  = err
        if repo != nil {
            title = test.Unique(title, tc.name)
            body  = test.Unique(body,  tc.name)
            req  := testIssueRequest(title, body)
            iss  := &github.Issue{Title: req.Title, Body: req.Body}
            tc.args.owner = repo.Owner
            tc.args.repo  = repo.Name
            tc.args.req   = req
            tc.want       = &Issue{ptr: iss, client: client}
        }
        return
    }

    repo  := GetFakeRepo(client)
    title := FAKE_ISSUE_TITLE
    body  := FAKE_ISSUE_BODY
    tests := []testCase{
        Case(0, nil,  "",    "",   ERR_NO_ISSUE_REQUEST),
        Case(1, repo, "",    "",   ERR_INVALID_CREATE),
        Case(2, repo, title, "",   ""),
        Case(3, repo, "",    body, ERR_INVALID_CREATE),
        Case(4, repo, title, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := CreateIssue(tt.args.client, tt.args.owner, tt.args.repo, tt.args.req)
            testVerifyIssue(fn, got, tt.want, t)
		})
	}
}

func TestGetIssue(t *testing.T) {
    const fn = "GetIssue"

    type args struct {
		client *Client
		owner  string
		repo   string
		number int
	}
    type testCase struct {
		name string
		args args
		want *Issue
        err  string
	}

    client := TestClient
    owner  := SAMPLE_ORG
    Case   := func(idx int, repo string, number int, match, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, owner, repo, number}
        tc.want = testIssue(match, number)
        tc.err  = err
        return
    }

    repo  := SAMPLE_REPO
    issue := SAMPLE_ISSUE
    match := SAMPLE_ISSUE_MAP[issue].Body
    tests := []testCase{
        Case(0, "",   0,      "",     ERR_NOT_FOUND),
        Case(1, repo, 0,      "",     ERR_NOT_FOUND),
        Case(2, repo, issue, match, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := GetIssue(tt.args.client, tt.args.owner, tt.args.repo, tt.args.number)
            testVerifyIssue(fn, got, tt.want, t)
		})
	}
}

func TestGetIssues(t *testing.T) {
    const fn = "GetIssues"

    type args struct {
		client *Client
		owner  string
		repo   string
	}
    type testCase struct {
		name string
		args args
		want []*Issue
        err  string
	}

    client := TestClient
    owner  := SAMPLE_ORG
    Case   := func(idx int, repo string, want []*Issue, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, owner, repo}
        tc.want = want
        tc.err  = err
        return
    }

    repo   := SAMPLE_REPO
    issue0 := testWantIssues()
    issue1 := testWantIssues(SAMPLE_ISSUES...)
    tests  := []testCase{
        Case(0, "",   issue0, ERR_NOT_FOUND),
        Case(1, repo, issue1, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got  := GetIssues(tt.args.client, tt.args.owner, tt.args.repo)
            testVerifyIssues(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported members - comments
// ============================================================================

func TestIssue_Comments(t *testing.T) {
    const fn = "Issue.Comments"

    type fields struct {
		Number int
		ptr    *github.Issue
		repo   *Repository
		client *Client
	}
    type testCase struct {
		name   string
		fields fields
		want   []*Comment
        err    string
	}

    client := TestClient
    repo   := SampleRepository(client)
    Case   := func(idx int, issue int, want []*Comment, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{Number: issue, repo: repo, client: client}
        tc.want   = want
        tc.err    = err
        return
    }

    issue := SAMPLE_ISSUE
    com0  := testWantComments()
    com1  := testWantComments(SAMPLE_COMMENTS...)
    tests := []testCase{
        Case(0, -1,    com0, ERR_NOT_FOUND), // NOTE: id 0 give unexpected results
        Case(1, issue, com1, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			issue := &Issue{
				Number: tt.fields.Number,
				ptr:    tt.fields.ptr,
				repo:   tt.fields.repo,
				client: tt.fields.client,
			}
            got := issue.Comments()
            testVerifyComments(fn, got, tt.want, t)
		})
	}
}

func TestIssue_GetComment(t *testing.T) {
    const fn = "Issue.GetComment"

    type fields struct {
		Number int
		ptr    *github.Issue
		repo   *Repository
		client *Client
	}
	type args struct {
		id int64
	}
    type testCase struct {
		name   string
		fields fields
		args   args
		want   *Comment
        err    string
	}

    client := TestClient
    repo   := SampleRepository(client)
    Case   := func(idx int, issue int, comment int64, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{Number: issue, repo: repo, client: client}
        tc.args   = args{comment}
        tc.err    = err
        if comment != 0 {
            tc.want = testComment(comment)
        }
        return
    }

    issue := SAMPLE_ISSUE
    com   := SAMPLE_COMMENT
    tests := []testCase{
        Case(0, 0,     0,   ERR_NOT_FOUND),
        Case(1, issue, com, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			issue := &Issue{
				Number: tt.fields.Number,
				ptr:    tt.fields.ptr,
				repo:   tt.fields.repo,
				client: tt.fields.client,
			}
            got := issue.GetComment(tt.args.id)
            testVerifyComment(fn, got, tt.want, t)
		})
	}
}

func TestIssue_CreateCommentFrom(t *testing.T) {
    const fn = "Issue.CreateCommentFrom"
    if test.Passive(fn, t) { return }

    type fields struct {
		Number int
		ptr    *github.Issue
		repo   *Repository
		client *Client
	}
	type args struct {
		text string
	}
    type testCase struct {
		name   string
		fields fields
		args   args
		want   *Comment
        err    string
	}

    client := TestClient
    repo   := GetFakeRepo(client)
    Case   := func(idx int, iss *Issue, body, err string) (tc testCase) {
        body      = test.Unique(body, tc.name)
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{repo: repo, client: client}
        tc.args   = args{body}
        tc.err    = err
        if iss != nil {
            tc.fields.Number = iss.Number
            tc.want = &Comment{ptr: testGithubComment(body)}
        }
        return
    }

    issue := GetFakeIssue(repo)
    body  := FAKE_COMMENT_BODY
    tests := []testCase{
        Case(0, nil,   "",   ERR_NOT_FOUND),
        Case(1, issue, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            issue := &Issue{
				Number: tt.fields.Number,
				ptr:    tt.fields.ptr,
				repo:   tt.fields.repo,
				client: tt.fields.client,
			}
            got := issue.CreateCommentFrom(tt.args.text)
            testVerifyComment(fn, got, tt.want, t)
		})
	}
}

func TestIssue_DeleteComment(t *testing.T) {
	const fn = "Issue.DeleteComment"
	if test.Passive(fn, t) { return }

    type fields struct {
		Number int
		ptr    *github.Issue
		repo   *Repository
		client *Client
	}
	type args struct {
		commentId int64
	}
	type testCase struct {
		name   string
        fields fields
		args   args
		err    string
	}

	client := TestClient
	repo   := GetFakeRepo(client)
	Case   := func(idx int, comment *Comment, err string) (tc testCase) {
		tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{repo: repo, client: client}
		tc.args   = args{0}
		tc.err    = err
        if comment != nil {
            tc.args.commentId = comment.ID()
        }
		return
	}

	issue := GetFakeIssue(repo)
    com   := GetFakeComment(issue)
    tests := []testCase{
        Case(0, nil, ERR_NOT_FOUND),
        Case(0, com, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer test.EvaluatePanic(tt.name, tt.err, t)
            issue := &Issue{
				Number: tt.fields.Number,
				ptr:    tt.fields.ptr,
				repo:   tt.fields.repo,
				client: tt.fields.client,
			}
			issue.DeleteComment(tt.args.commentId)
		})
	}
}

// ============================================================================
// Tests - Exported members - rendering
// ============================================================================

func TestIssue_Details(t *testing.T) {
    const fn = "Issue.Details"

    type fields struct {
		Number int
		ptr    *github.Issue
		repo   *Repository
		client *Client
	}
    type testCase struct{
		name   string
		fields fields
		want   string
        err    string
    }

    client := TestClient
    repo   := SampleRepository(client)
    Case   := func(idx int, body, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{repo: repo, client: client}
        tc.err    = err
        if body == "" {
            tc.fields.ptr = &github.Issue{}
        } else {
            body = test.Unique(body, tc.name)
            tc.fields.ptr = &github.Issue{Body: github.Ptr(body)}
            tc.want = fmt.Sprintf("Body              %s", body)
        }
        return
    }

    body  := FAKE_COMMENT_BODY
    tests := []testCase {
        Case(0, "",   ""),
        Case(1, body, ""),
    }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			issue := &Issue{
				Number: tt.fields.Number,
				ptr:    tt.fields.ptr,
				repo:   tt.fields.repo,
				client: tt.fields.client,
			}
			if got, want := issue.Details(), tt.want; got != want {
				t.Errorf("%s() = %q, want %q", fn, got, want)
			}
		})
	}
}

// ============================================================================
// Test setup
// ============================================================================

// Initialize variables related to testing GitHub issues.
func testSetup_Issue() {
    // no-op
}

// Clean up variables related to testing GitHub issues.
func testTeardown_Issue() {
    // NOTE: FakeIssues are deleted by deleting FakeRepositories.
    PreserveFakes()
}
