// Github/org_func_test.go

package Github

import (
	"testing"

	"lib.virginia.edu/agita/test"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Internal functions
// ============================================================================

// A list of organization requirements suitable for comparison against a list
// of real organization objects.
func testWantOrgs(names ...string) []*github.Organization {
    result := make([]*github.Organization, 0, len(names))
    for _, name := range names {
        result = append(result, &github.Organization{Login: github.Ptr(name)})
    }
    return result
}

// ============================================================================
// Internal functions - verification
// ============================================================================

// Fail if the GitHub organizations don't match the provided list.
//  NOTE: extra `got` items are ignored unless `want` is empty
func testVerifyOrgs(fn string, got, want []*github.Organization, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if len1, len2 := len(got), len(want); len1 < len2 {
        t.Errorf("%s() = %v orgs, want %v orgs", fn, len1, len2)
        t.Errorf("%s() got  %q", fn, orgNames(got))
        t.Errorf("%s() want %q", fn, orgNames(want))
    } else if len2 == 0 {
        if len1 > len2 {
            t.Errorf("%s() = %v orgs, want no orgs", fn, len1)
        }
    } else {
        g := map[string]*github.Organization{}
        for _, org := range got {
            if (org != nil) && (org.Name != nil) {
                g[*org.Name] = org
            }
        }
        for _, org := range got {
            if (org != nil) && (org.Name != nil) && (g[*org.Name] == nil) {
                t.Errorf("%s() missing org %q", fn, *org.Name)
            }
        }
    }
}

// Render organizations as a sorted list.
func orgNames(orgs []*github.Organization) string {
    names := []string{}
    for _, org := range orgs {
        if org != nil {
            if name := org.Login; name != nil {
                names = append(names, *name)
            }
        }
    }
    return test.SortedListing(names)
}
