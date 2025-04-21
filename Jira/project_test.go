// Jira/project_test.go

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

func TestNewProjectType(t *testing.T) {
    const fn = "NewProjectType"

	type args struct {
		client  *Client
		project *jira.Project
	}
    type testCase struct {
		name string
		args args
		want *Project
        err  string
	}

    client := TestClient
    Case   := func(idx int, key, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, nil}
        tc.err  = err
        if err == "" {
            proj := &jira.Project{Key: key}
            tc.want = &Project{proj, client}
            tc.args.project = proj
        }
        return
    }

	tests := []testCase{
        Case(0, "",     ERR_NIL_PROJECT),
        Case(1, "FAKE", ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := NewProjectType(tt.args.client, tt.args.project)
            testVerifyProject(fn, got, tt.want, t)
		})
	}
}

func TestGetProjects(t *testing.T) {
    const fn = "GetProjects"

	type args struct {
		client *Client
	}
    type testCase struct {
		name string
		args args
		want []*Project
        err  string
	}

    client := TestClient
    Case   := func(idx int, org string, want []*Project, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client}
        tc.want = want
        tc.err  = err
        return
    }

    projs := client.GetProjects() // TODO: improve
    tests := []testCase{
        Case(0, "", projs, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := GetProjects(tt.args.client)
            testVerifyProjects(fn, got, tt.want, t)
		})
	}
}

func TestGetProjectByKey(t *testing.T) {
    const fn = "GetProjectByKey"

    type args struct {
		client *Client
		key    ProjKey
	}
    type testCase struct {
		name string
		args args
		want *Project
        err  string
	}

	client := TestClient
	Case   := func(idx int, key string, err string) (tc testCase) {
		tc.name = test.CaseName(fn, idx)
		tc.args = args{client, key}
		tc.err  = err
		if err == "" {
			tc.want = &Project{ptr: &jira.Project{Key: key}, client: client}
		}
		return
	}

    proj  := SAMPLE_PROJ
    tests := []testCase{
        Case(0, "",   ERR_NO_PROJECT),
        Case(1, proj, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := GetProjectByKey(tt.args.client, tt.args.key)
            testVerifyProject(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported members - issues
// ============================================================================

func TestProject_Issues(t *testing.T) {
    const fn = "Project.Issues"

	type fields struct {
		ptr    *jira.Project
		client *Client
	}
    type testCase struct {
		name   string
		fields fields
		want   []Issue
        err    string
	}

    client := TestClient
    proj   := SampleProject(client)
    Case   := func(idx int, want []Issue, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{proj.ptr, client}
        tc.want   = want
        tc.err    = err
        return
    }

    issues := testWantIssues(SAMPLE_ISSUES...)
    tests := []testCase{
        Case(0, issues, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			proj := &Project{
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            got := proj.Issues()
            testVerifyIssues(fn, got, tt.want, t)
		})
	}
}

func TestProject_GetIssue(t *testing.T) {
    const fn = "Project.GetIssue"

	type fields struct {
		ptr    *jira.Project
		client *Client
	}
	type args struct {
		key IssueKey
	}
    type testCase struct {
		name   string
		fields fields
		args   args
		want   *Issue
        err    string
	}

    client := TestClient
    proj   := SampleProject(client)
    Case   := func(idx int, key IssueKey, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{proj.ptr, client}
        tc.args   = args{key}
        tc.want   = testIssue(key)
        tc.err    = err
        if key != "" {
            tc.want.ptr.Fields.Summary = SAMPLE_ISSUE_MAP[key].Body
        }
        return
    }

    issue := SAMPLE_ISSUE
    tests := []testCase{
        Case(0, "",    ERR_REQUEST_FAILED),
        Case(1, issue, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			proj := &Project{
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            got := proj.GetIssue(tt.args.key)
            testVerifyIssue(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Exported members - rendering
// ============================================================================

func TestProject_Details(t *testing.T) {
    const fn = "Project.Details"

	type fields struct {
		ptr    *jira.Project
		client *Client
	}
    type testCase struct {
		name   string
		fields fields
		want   string
        err    string
	}

    client := TestClient
    Case   := func(idx int, key, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{nil, client}
        tc.err    = err
        if key == "" {
            tc.fields.ptr = &jira.Project{}
        } else {
            key = test.Unique(key, tc.name)
            tc.fields.ptr = &jira.Project{Key: key}
            tc.want = fmt.Sprintf("Key             %s", key)
        }
        return
    }

    tests := []testCase{
        Case(0, "",          ""),
        Case(1, SAMPLE_PROJ, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			proj := &Project{
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
			if got, want := proj.Details(), tt.want; got != want {
				t.Errorf("%s() = %q, want %q", fn, got, want)
			}
		})
	}
}

// ============================================================================
// Test setup
// ============================================================================

// Initialize variables related to testing Jira projects.
func testSetup_Project() {
    // no op
}

// Clean up variables related to testing Jira projects.
func testTeardown_Project() {
    // no op
}
