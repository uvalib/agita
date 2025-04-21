// Jira/client_test.go

package Jira

import (
	"testing"

	"lib.virginia.edu/agita/test"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Tests - Exported functions
// ============================================================================

func TestNewClient(t *testing.T) {
    const fn = "NewClient"

    type testCase struct {
		name string
		want *Client
        err  string
	}

    Case := func(idx int, want *Client) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.want = want
        return
    }

    client := TestClient
    tests  := []testCase{
        Case(0, client),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := NewClient()
            testVerifyClient(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported methods - projects
// ============================================================================

func TestClient_GetProjects(t *testing.T) {
    const fn = "Client.GetProjects"

	type fields struct {
		ptr *jira.Client
	}
    type testCase struct {
		name   string
		fields fields
		want   []*Project
        err    string
	}

    client := TestClient
    Case   := func(idx int, org string, want []*Project, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{client.ptr}
        tc.want   = want
        tc.err    = err
        return
    }

    projs := client.GetProjects() // TODO: improve
    tests := []testCase{
        Case(0, "", projs, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			c := &Client{
				ptr: tt.fields.ptr,
			}
            got := c.GetProjects()
            testVerifyProjects(fn, got, tt.want, t)
		})
	}
}

func TestClient_GetProjectByKey(t *testing.T) {
    const fn = "Client.GetProjectByKey"

	type fields struct {
		ptr *jira.Client
	}
	type args struct {
		key ProjKey
	}
    type testCase struct {
		name   string
		fields fields
		args   args
		want   *Project
        err    string
	}

	client := TestClient
	Case   := func(idx int, key string, err string) (tc testCase) {
		tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{client.ptr}
		tc.args   = args{key}
		tc.err    = err
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
			c := &Client{
				ptr: tt.fields.ptr,
			}
            got := c.GetProjectByKey(tt.args.key)
            testVerifyProject(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Test setup
// ============================================================================

// Initialize variables related to testing Jira clients.
func testSetup_Client() {
    TestClient = NewClient()
}

// Clean up variables related to testing Jira clients.
func testTeardown_Client() {
    // no op
}
