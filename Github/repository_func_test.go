// Github/comment_func_test.go

package Github

import (
	"testing"

	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/test"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Repository members
// ============================================================================

// Create a copy of `ptr` in which the ID is set to the given value.
// The original github.Repository is untouched.
func (r *Repository) setID(id int64) *Repository {
    if r != nil {
        if r.ptr == nil {
            r.ptr = &github.Repository{}
        } else {
            r.ptr = github.Ptr(*r.ptr)
        }
        r.ptr.ID = github.Ptr(id)
    }
    return r
}

// ============================================================================
// Internal functions
// ============================================================================

// A list of repository requirements suitable for comparison against a list
// of real repository objects.
func testWantRepos(names ...string) []*Repository {
    result := make([]*Repository, 0, len(names))
    for _, name := range names {
        result = append(result, &Repository{Name: name})
    }
    return result
}

// ============================================================================
// Internal functions - verification
// ============================================================================

// Fail if the GitHub repository doesn't match the provided criteria.
func testVerifyRepository(fn string, got, want *Repository, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if w := RepoNumber(want); w != 0 {
        if g := RepoNumber(got); g != w {
            t.Errorf("%s().ID = %d, want %d", fn, g, w)
        }
    }
    if w := RepoOwner(want); w != "" {
        if g, simple := RepoOwner(got), re.IsSimple(w); simple && (g != w) {
            t.Errorf("%s().Owner = %q, want %q", fn, g, w)
        } else if !simple && !re.New(w).Match(g) {
            t.Errorf("%s().Owner = %q, want match of `%s`", fn, g, w)
        }
    }
    if w := RepoName(want); w != "" {
        if g, simple := RepoName(got), re.IsSimple(w); simple && (g != w) {
            t.Errorf("%s().Name = %q, want %q", fn, g, w)
        } else if !simple && !re.New(w).Match(g) {
            t.Errorf("%s().Name = %q, want match of `%s`", fn, g, w)
        }
    }
}

// Fail if the GitHub repositories don't match the provided list.
//  NOTE: extra `got` items are ignored unless `want` is empty
func testVerifyRepos(fn string, got, want []*Repository, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if len1, len2 := len(got), len(want); len1 < len2 {
        t.Errorf("%s() = %v repos, want %v repos", fn, len1, len2)
        t.Errorf("%s() got  %q", fn, repoNames(got))
        t.Errorf("%s() want %q", fn, repoNames(want))
    } else if len2 == 0 {
        if len1 > len2 {
            t.Errorf("%s() = %v repos, want no repos", fn, len1)
        }
    } else {
        g := map[string]*Repository{}
        for _, repo := range got {
            if (repo != nil) && (repo.Name != "") {
                g[repo.Name] = repo
            }
        }
        for _, repo := range got {
            if (repo != nil) && (repo.Name != "") && (g[repo.Name] == nil) {
                t.Errorf("%s() missing repo %q", fn, repo.Name)
            }
        }
    }
}

// Render repositories as a sorted list.
func repoNames(repos []*Repository) string {
    names := []string{}
    for _, repo := range repos {
        if repo != nil {
            if name := repo.Name; name != "" {
                names = append(names, name)
            }
        }
    }
    return test.SortedListing(names)
}
