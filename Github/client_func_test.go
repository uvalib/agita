// Github/client_func_test.go

package Github

import (
	"testing"

	"lib.virginia.edu/agita/test"
	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Internal functions - verification
// ============================================================================

// Fail if the GitHub client doesn't match the provided criteria.
func testVerifyClient(fn string, got, want *Client, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if g, w := got.BaseURL(), want.BaseURL(); !util.SameUrl(g, w) {
        t.Errorf("%s().BaseURL = %v, want %v", fn, g, w)
    }
}
