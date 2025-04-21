// Jira/client.go

package Jira

import (
	"net/http"
	"net/url"

	"lib.virginia.edu/agita/log"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

// Application object referencing a Jira client object.
type Client struct {
    ptr *jira.Client
}

// ============================================================================
// Internal variables
// ============================================================================

var mainClient *Client

// ============================================================================
// Exported functions
// ============================================================================

// The default client used for application objects which do not specify one.
func MainClient() *Client {
    if mainClient == nil {
        mainClient = NewClient()
    }
    return mainClient
}

// Get a new authorized Client instance for accessing all projects within the
// Jira at BASE_URL.
func NewClient() (result *Client) {
    token := authToken()
    httpClient := &http.Client{Transport: &jira.PATAuthTransport{Token: token}}
    client, err := jira.NewClient(httpClient, BASE_URL)
    if log.ErrorValue(err) == nil {
        result = &Client{ptr: client}
    }
    return
}

// ============================================================================
// Exported members - properties
// ============================================================================

// Return the underlying BaseURL value.
func (c *Client) BaseURL() url.URL {
    if noClient(c) { return url.URL{} }
    return c.ptr.GetBaseURL()
}

// ============================================================================
// Internal functions
// ============================================================================

// Indicate whether the argument is missing a Client object.
func noClient(c *Client) bool {
    return (c == nil) || (c.ptr == nil)
}

// ============================================================================
// Exported methods - projects
// ============================================================================

// Get all projects for the Jira associated with the client.
func (c *Client) GetProjects() []*Project {
    return GetProjects(c)
}

// Get the project with the given project key.
func (c *Client) GetProjectByKey(key ProjKey) *Project {
    return GetProjectByKey(c, key)
}
