// Github/user_test.go

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

func TestNewUserType(t *testing.T) {
    const fn = "NewUserType"

    type args struct {
		client *Client
		login  string
	}
    type testCase struct {
		name string
		args args
		want *User
        err  string
	}

    Case := func(idx int, client *Client, login string, want *User, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, login}
        tc.want = want
        tc.err  = err
        return
    }

    client := TestClient
    user   := SampleUser(client)
    name   := user.Login
    tests  := []testCase{
        Case(0, nil,    "",   nil,  ERR_NIL_CLIENT),
        Case(1, client, "",   nil,  ERR_NO_LOGIN),
        Case(2, client, name, user, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := NewUserType(tt.args.client, tt.args.login)
            testVerifyUser(fn, got, tt.want, t)
		})
	}
}

func TestGetUser(t *testing.T) {
    const fn = "GetUser"

    type args struct {
		client *Client
		login  string
	}
    type testCase struct {
		name string
		args args
		want *User
        err  string
	}

    Case := func(idx int, client *Client, login string, want *User, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.args = args{client, login}
        tc.want = want
        tc.err  = err
        return
    }

    client := TestClient
    user   := SampleUser(client)
    name   := user.Login
    tests  := []testCase{
        Case(0, nil,    "",   nil,  ERR_NIL_CLIENT),
        Case(1, client, "",   nil,  ERR_NO_LOGIN),
        Case(2, client, name, user, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
            got := GetUser(tt.args.client, tt.args.login)
            testVerifyUser(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported members
// ============================================================================

func TestUser_Orgs(t *testing.T) {
    const fn = "User.Orgs"

    type fields struct {
		Login  string
		ptr    *github.User
		client *Client
	}
    type testCase struct {
		name   string
		fields fields
		want   []*github.Organization
        err    string
	}

    client := TestClient
    Case   := func(idx int, login string, want []*github.Organization, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.fields = fields{login, nil, client}
        tc.want = want
        tc.err  = err
        return
    }

    user  := SAMPLE_USER
    orgs0 := testWantOrgs("samvera", "uvalib")
    orgs1 := testWantOrgs("uvalib")
    tests := []testCase{
        Case(0, "",   orgs0, ""), // All orgs for the current user.
        Case(0, user, orgs1, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			u := &User{
				Login:  tt.fields.Login,
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            got := u.Orgs()
            testVerifyOrgs(fn, got, tt.want, t)
		})
	}
}

func TestUser_Repos(t *testing.T) {
    const fn = "User.Repos"

    type fields struct {
		Login  string
		ptr    *github.User
		client *Client
	}
    type testCase struct {
		name   string
		fields fields
		want   []*Repository
        err    string
	}

    client := TestClient
    Case   := func(idx int, login string, want []*Repository, err string) (tc testCase) {
        tc.name = test.CaseName(fn, idx)
        tc.fields = fields{login, nil, client}
        tc.want = want
        tc.err  = err
        return
    }

    user  := SAMPLE_USER
    reps0 := testWantRepos()
    reps1 := testWantRepos("IAS3API", "junk", "tools", "virgo4")
    tests := []testCase{
        Case(0, "",   reps0, ERR_NOT_FOUND),
        Case(0, user, reps1, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            defer test.EvaluatePanic(tt.name, tt.err, t)
			u := &User{
				Login:  tt.fields.Login,
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
            got := u.Repos()
            testVerifyRepos(fn, got, tt.want, t)
		})
	}
}

// ============================================================================
// Tests - Exported members - rendering
// ============================================================================

func TestUser_Details(t *testing.T) {
    const fn = "User.Details"

	type fields struct {
		Login  string
		ptr    *github.User
		client *Client
	}
    type testCase struct {
		name   string
		fields fields
		want   string
        err    string
	}

    client := TestClient
    Case   := func(idx int, login, err string) (tc testCase) {
        tc.name   = test.CaseName(fn, idx)
        tc.fields = fields{login, nil, client}
        tc.err    = err
		if login == "" {
			tc.fields.ptr = &github.User{}
		} else {
			tc.fields.ptr = &github.User{Login: github.Ptr(login)}
			tc.want = fmt.Sprintf("Login                   %s", login)
		}
        return
    }

    user  := SAMPLE_USER
    tests := []testCase{
		Case(0, "",   ""),
		Case(1, user, ""),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer test.EvaluatePanic(tt.name, tt.err, t)
			u := &User{
				Login:  tt.fields.Login,
				ptr:    tt.fields.ptr,
				client: tt.fields.client,
			}
			if got, want := u.Details(), tt.want; got != want {
                t.Errorf("%s() = %q, want %q", fn, got, want)
			}
		})
	}
}

// ============================================================================
// Test setup
// ============================================================================

// Initialize variables related to testing GitHub users.
func testSetup_User() {
    // no op
}

// Clean up variables related to testing GitHub users.
func testTeardown_User() {
    // no op
}
