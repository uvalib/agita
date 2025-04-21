// Jira/project_func_test.go

package Jira

import (
	"testing"

	"lib.virginia.edu/agita/test"
)

// ============================================================================
// Internal functions - verification
// ============================================================================

// Fail if the Jira project doesn't match the provided criteria.
func testVerifyProject(fn string, got, want *Project, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if w := ProjectKey(want); w != "" {
        if g := ProjectKey(got); g != w {
            t.Errorf("%s().Key = %q, want %q", fn, g, w)
        }
    }
}

// Fail if the Jira projects don't match the provided list.
//  NOTE: extra `got` items are ignored unless `want` is empty
func testVerifyProjects(fn string, got, want []*Project, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if len1, len2 := len(got), len(want); len1 < len2 {
        t.Errorf("%s() = %v projects, want %v projects", fn, len1, len2)
        t.Errorf("%s() got  %q", fn, projKeys(got))
        t.Errorf("%s() want %q", fn, projKeys(want))
    } else if len2 == 0 {
        if len1 > len2 {
            t.Errorf("%s() = %v projects, want no projects", fn, len1)
        }
    } else {
        g := map[string]*Project{}
        for _, proj := range got {
            if proj != nil {
                g[proj.Name()] = proj
            }
        }
        for _, proj := range got {
            if (proj != nil) && (g[proj.Name()] == nil) {
                t.Errorf("%s() missing project %q", fn, proj.Name())
            }
        }
    }
}

// Render projects as a sorted list.
func projKeys(projects []*Project) string {
    keys := []string{}
    for _, proj := range projects {
        if proj != nil {
            if key := proj.Key(); key != "" {
                keys = append(keys, key)
            }
        }
    }
    return test.SortedListing(keys)
}
