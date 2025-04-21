// Github/user_func_test.go

package Github

import (
	"testing"

	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/test"
)

// ============================================================================
// Internal functions - verification
// ============================================================================

// Fail if the GitHub user doesn't match the provided criteria.
func testVerifyUser(fn string, got, want *User, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if w := Account(want); w != "" {
        if g, simple := Account(got), re.IsSimple(w); simple && (g != w) {
            t.Errorf("%s().Login = %q, want %q", fn, g, w)
        } else if !simple && !re.New(w).Match(g) {
            t.Errorf("%s().Login = %q, want match of `%s`", fn, g, w)
        }
    }
}
