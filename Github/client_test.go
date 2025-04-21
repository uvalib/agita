// Github/client_test.go

package Github

import (
	"testing"

	"lib.virginia.edu/agita/test"
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
// Tests - Exported members
// ============================================================================

func TestClient_GetOrgRepos(t *testing.T) {
    const fn = "Client.GetOrgRepos"

    type args struct {
		org string
	}
    type testCase struct {
		name string
		args args
		want []*Repository
        err  string
	}

    Case := func(idx int, org string, want []*Repository, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{org}
        tc.want = want
        tc.err  = err
        return
    }

    client := TestClient
    repos  := client.GetOrgRepos(SAMPLE_ORG) // TODO: improve
    tests  := []testCase{
        Case(0, "", repos, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            repos := client.GetOrgRepos(tt.args.org)
            testVerifyRepos(fn, repos, tt.want, t)
		})
	}
}

func TestClient_GetRepository(t *testing.T) {
    const fn = "Client.GetRepository"

    type args struct {
		owner string
		repo  string
	}
    type testCase struct {
		name string
		args args
		want *Repository
        err  string
	}

    Case := func(idx int, owner, repo string, want *Repository, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{owner, repo}
        tc.want = want
        tc.err  = err
        return
    }

    client := TestClient
    repo   := SampleRepository(client)
    name   := repo.Name
    tests  := []testCase{
        Case(0, "", "",    nil,   ""),
        Case(1, "", name, repo, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := TestClient.GetRepository(tt.args.owner, tt.args.repo)
            testVerifyRepository(fn, got, tt.want, t)
		})
	}
}

func TestClient_GetUser(t *testing.T) {
    const fn = "Client.GetUser"

    type args struct {
		login string
	}
    type testCase struct {
		name string
		args args
		want *User
        err  string
	}

    Case := func(idx int, user string, want *User, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{user}
        tc.want = want
        tc.err  = err
        return
    }

    client := TestClient
    user   := SampleUser(client)
    login  := user.Login
	tests  := []testCase{
        Case(0, "",    nil,  ERR_NO_LOGIN),
        Case(1, login, user, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			got := client.GetUser(tt.args.login)
            testVerifyUser(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Test setup
// ============================================================================

// Initialize variables related to testing GitHub clients.
func testSetup_Client() {
    TestClient = NewClient()
}

// Clean up variables related to testing GitHub clients.
func testTeardown_Client() {
    // no op
}
