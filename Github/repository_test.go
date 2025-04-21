// Github/repository_test.go

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

func TestNewRepositoryType(t *testing.T) {
    const fn = "NewRepositoryType"

    type args struct {
		client *Client
		owner  string
		name   string
	}
    type testCase struct {
		name string
		args args
		want *Repository
        err  string
	}

    client := TestClient
    owner  := SAMPLE_ORG
    Case   := func(idx int, repo string, want *Repository, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, owner, repo}
        tc.want = want
        tc.err  = err
        return
    }

    repo  := SampleRepository(client).setID(0)
    name  := repo.Name
    tests := []testCase{
        Case(0, "",   nil,  ERR_NO_REPO_NAME),
        Case(1, name, repo, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := NewRepositoryType(tt.args.client, tt.args.owner, tt.args.name)
            testVerifyRepository(fn, got, tt.want, t)
		})
	}
}

func TestAsRepositoryType(t *testing.T) {
    const fn = "AsRepositoryType"

    type args struct {
		client *Client
		repo   *github.Repository
	}
    type testCase struct {
		name string
		args args
		want *Repository
        err  string
	}

    client := TestClient
    Case   := func(idx int, repo *github.Repository, want *Repository, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, repo}
        tc.want = want
        tc.err  = err
        return
    }

    repo  := SampleRepository(client)
    tests := []testCase{
        Case(0, nil,       nil,   ERR_NO_REPO),
        Case(1, repo.ptr, repo, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := AsRepositoryType(tt.args.client, tt.args.repo)
            testVerifyRepository(fn, got, tt.want, t)
		})
	}
}

func TestGetRepository(t *testing.T) {
    const fn = "GetRepository"

    type args struct {
		client *Client
		owner  string
		name   string
	}
    type testCase struct {
		name string
		args args
		want *Repository
        err  string
	}

    client := TestClient
    owner  := SAMPLE_ORG
    Case   := func(idx int, repo, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, owner, repo}
        tc.err  = err
        if repo != "" {
            tc.want = &Repository{Owner: owner, Name: repo}
        }
        return
    }

    repo  := SAMPLE_REPO
    tests := []testCase{
        Case(0, "",   ""),
        Case(1, repo, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := GetRepository(tt.args.client, tt.args.owner, tt.args.name)
            testVerifyRepository(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported functions
// ============================================================================

func TestCreateRepository(t *testing.T) {
    const fn = "CreateRepository"
    if test.Passive(fn, t) { return }

    type args struct {
		client *Client
		data   *RepositoryRequest
	}
    type testCase struct {
		name string
		args args
        want *Repository
        err  string
	}

    client := TestClient
    Case   := func(idx int, data *RepositoryRequest, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, data}
        tc.err  = err
        var repo *github.Repository
        if data == nil {
            repo = TemplateRepoDataAsRepository()
            tc.args.data = &RepositoryRequest{Repository: *repo}
        } else {
            repo = &data.Repository
        }
        tc.want = &Repository{ptr: repo}
        return
    }
    InitFakeRepos(client)

    tests := []testCase{
        Case(0, nil, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            repo := CreateRepository(tt.args.client, tt.args.data)
            defer AddFakeRepo(tt.args.client, repo)
            testVerifyRepository(fn, repo, tt.want, t)
		})
	}
}

func TestDeleteRepository(t *testing.T) {
    const fn = "DeleteRepository"
    if test.Passive(fn, t) { return }

    type args struct {
		client *Client
		owner  string
        name   string
	}
    type testCase struct {
		name string
		args args
        err  string
	}

    client := TestClient
    owner  := SAMPLE_ORG
    Case   := func(idx int, name, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, owner, name}
        tc.err  = err
        return
    }

    repo  := GetFakeRepo(client)
    name  := repo.Name
    tests := []testCase{
        Case(0, "",   ERR_NO_REPO_GIVEN),
        Case(1, name, ""),
	}
    defer ClearFakeRepo(client, repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            DeleteRepository(tt.args.client, tt.args.owner, tt.args.name)
		})
	}
}

// ============================================================================
// Tests - Exported members - issues
// ============================================================================

func TestRepository_ImportIssue(t *testing.T) {
    const fn = "Repository.ImportIssue"
    if test.Passive(fn, t) { return }

	type fields struct {
		Owner  string
		Name   string
		ptr    *github.Repository
		client *Client
	}
	type args struct {
		imp      *IssueImport
		comments []*CommentImport
	}
    type testCase struct {
		name   string
		fields fields
		args   args
        err    string
	}

    client := TestClient
    repo   := GetFakeRepo(client)
    Case   := func(idx int, title, body, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{repo.Owner, repo.Name, repo.ptr, repo.client}
        tc.args   = args{comments: []*CommentImport{}}
        tc.err    = err
        if err == "" {
            title = test.Unique(title, tc.name)
            body  = test.Unique(body,  tc.name)
            tc.args.imp = testIssueImport(title, body)
        }
        return
    }

    title := FAKE_ISSUE_TITLE
    body  := FAKE_ISSUE_BODY
	tests := []testCase{
        Case(0, "",    "",   ERR_NO_ISSUE_IMPORT),
        Case(1, title, "",   ERR_NO_ISSUE_IMPORT),
        Case(2, "",    body, ERR_NO_ISSUE_IMPORT),
        Case(3, title, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			r := &Repository{
				Owner:  tt.fields.Owner,
				Name:   tt.fields.Name,
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            if id := r.ImportIssue(tt.args.imp, tt.args.comments...); id != 0 {
                testVerifyIssueImport(fn, r, id, t)
            }
		})
	}
}

func TestRepository_CreateIssue(t *testing.T) {
    const fn = "Repository.CreateIssue"
    if test.Passive(fn, t) { return }

	type fields struct {
		Owner  string
		Name   string
		ptr    *github.Repository
		client *Client
	}
	type args struct {
		req *github.IssueRequest
	}
    type testCase struct {
		name   string
		fields fields
		args   args
		want   *Issue
        err    string
	}

    client := TestClient
    repo   := GetFakeRepo(client)
    Case   := func(idx int, title, body, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{repo.Owner, repo.Name, repo.ptr, repo.client}
        tc.args   = args{nil}
        tc.err    = err
        if title != "" {
            title = test.Unique(title, tc.name)
            body  = test.Unique(body,  tc.name)
            tc.want = testIssue(body, 0)
            tc.args.req = testIssueRequest(title, body)
        }
        return
    }

    title := FAKE_ISSUE_TITLE
    body  := FAKE_ISSUE_BODY
    tests := []testCase{
        Case(0, "",    "",   ERR_NO_ISSUE_REQUEST),
        Case(1, title, "",   ""),
        Case(2, "",    body, ERR_NO_ISSUE_REQUEST),
        Case(3, title, body, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			r := &Repository{
				Owner:  tt.fields.Owner,
				Name:   tt.fields.Name,
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            got := r.CreateIssue(tt.args.req)
            testVerifyIssue(fn, got, tt.want, t)
		})
	}
}

func TestRepository_GetIssue(t *testing.T) {
    const fn = "Repository.GetIssue"

    type fields struct {
		Owner  string
		Name   string
		ptr    *github.Repository
		client *Client
	}
	type args struct {
		number int
	}
    type testCase struct {
		name   string
		fields fields
		args   args
		want   *Issue
        err    string
	}

    client := TestClient
    repo   := SampleRepository(client)
    Case   := func(idx int, number int, match, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{repo.Owner, repo.Name, repo.ptr, client}
        tc.args   = args{number}
        tc.want   = testIssue(match, number)
        tc.err    = err
        return
    }

    issue := SAMPLE_ISSUE
    match := SAMPLE_ISSUE_MAP[issue].Body
    tests := []testCase{
        Case(0, 0,     "",    ERR_NOT_FOUND),
        Case(1, issue, match, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			r := &Repository{
				Owner:  tt.fields.Owner,
				Name:   tt.fields.Name,
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            got := r.GetIssue(tt.args.number)
            testVerifyIssue(fn, got, tt.want, t)
		})
	}
}

func TestRepository_GetIssues(t *testing.T) {
    const fn = "Repository.GetIssues"

	type fields struct {
		Owner  string
		Name   string
		ptr    *github.Repository
		client *Client
	}
    type testCase struct {
		name   string
		fields fields
		want   []*Issue
        err    string
	}

    client := TestClient
    repo   := SampleRepository(client)
    Case   := func(idx int, want []*Issue, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{repo.Owner, repo.Name, repo.ptr, client}
        tc.want   = want
        tc.err    = err
        return
    }

    issues := testWantIssues(SAMPLE_ISSUES...)
	tests  := []testCase{
        Case(0, issues, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			r := &Repository{
				Owner:  tt.fields.Owner,
				Name:   tt.fields.Name,
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            got  := r.GetIssues()
            testVerifyIssues(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported members - rendering
// ============================================================================

func TestRepository_Details(t *testing.T) {
    const fn = "Repository.Details"

	type fields struct {
		Owner  string
		Name   string
		ptr    *github.Repository
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
    Case   := func(idx int, name, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{repo.Owner, repo.Name, repo.ptr, client}
        tc.err    = err
        if name == "" {
            tc.fields.ptr = &github.Repository{}
        } else {
            name = test.Unique(name, tc.name)
            tc.fields.ptr = &github.Repository{Name: github.Ptr(name)}
            tc.want = fmt.Sprintf("Name            %s", name)
        }
        return
    }

    name  := repo.Name
    tests := []testCase{
        Case(0, "",   ""),
        Case(1, name, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			r := &Repository{
				Owner:  tt.fields.Owner,
				Name:   tt.fields.Name,
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
			if got, want := r.Details(), tt.want; got != want {
				t.Errorf("%s() = %q, want %q", fn, got, want)
			}
		})
	}
}

// ============================================================================
// Test setup
// ============================================================================

// Initialize variables related to testing GitHub repositories.
func testSetup_Repository() {
    // no-op
}

// Clean up variables related to testing GitHub repositories.
func testTeardown_Repository() {
    switch {
        case postCleanAll:   DeleteAllTemporaryRepos(TestClient)
        case postCleanFakes: DeleteFakeRepos(TestClient)
    }
}
