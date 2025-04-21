// Jira/global_test.go
//
// Test values used throughout the package.

package Jira

import (
	"os"
	"testing"
)

// ============================================================================
// Variables
// ============================================================================

var TestClient *Client

// ============================================================================
// Internal variables - samples
// ============================================================================

var sampleProject *Project

// ============================================================================
// Functions - samples
// ============================================================================

// An existing project object acquired from Jira.
func SampleProject(client *Client) *Project {
    if sampleProject == nil {
        if client == nil {
            client = TestClient
        }
        sampleProject = client.GetProjectByKey(SAMPLE_PROJ)
    }
    return sampleProject
}

// ============================================================================
// Test setup
// ============================================================================

// Invoked once by "go test" for any/all tests in this package.
func TestMain(m *testing.M) {
    testSetup()
    code := m.Run()
    testTeardown()
    os.Exit(code)
}

// Setup functions run before tests begin.
func testSetup() {
    testSetup_Client()
    testSetup_Project()
    testSetup_Issue()
    testSetup_Comment()
}

// Teardown functions run after tests finish.
func testTeardown() {
    testTeardown_Comment()
    testTeardown_Issue()
    testTeardown_Project()
    testTeardown_Client()
}
