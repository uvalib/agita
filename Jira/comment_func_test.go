// Jira/comment_func_test.go

package Jira

import (
	"strconv"
	"testing"

	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/test"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Internal functions
// ============================================================================

// Generate a minimal Comment object.
func testJiraComment(body string) *jira.Comment {
	return &jira.Comment{Body: body}
}

// Generate a minimal Comment object.
func testComment[T int|string](arg T) *Comment {
    var id string
    switch v := any(arg).(type) {
        case int:    id = strconv.Itoa(v)
        case string: id = v
    }
    return &Comment{ptr: &jira.Comment{ID: id}}
}

// A list of comment requirements suitable for comparison against a list of
// real comment objects.
func testWantComments(ids ...int) []Comment {
    result := make([]Comment, 0, len(ids))
    for _, id := range ids {
        result = append(result, *testComment(id))
    }
    return result
}

// ============================================================================
// Internal functions - verification
// ============================================================================

// Fail if the GitHub comment doesn't match the provided criteria.
func testVerifyComment(fn string, got, want *Comment, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if w := CommentNumber(want); w != "" {
        if g := CommentNumber(got); g != w {
            t.Errorf("%s().ID = %s, want %s", fn, g, w)
        }
    }
    if w := CommentBody(want); w != "" {
        if g, simple := CommentBody(got), re.IsSimple(w); simple && (g != w) {
            t.Errorf("%s().Body = %q, want %q", fn, g, w)
        } else if !simple && !re.New(w).Match(g) {
            t.Errorf("%s().Body = %q, want match of `%s`", fn, g, w)
        }
    }
}

// Fail if the GitHub comments don't match the provided list.
//  NOTE: extra `got` items are ignored unless `want` is empty
func testVerifyComments(fn string, got, want []Comment, t *testing.T) {
    if test.CheckForNils(fn, got, want, t) {
        return
    }
    if len1, len2 := len(got), len(want); len1 < len2 {
        t.Errorf("%s() = %v comments, want %v comments", fn, len1, len2)
        for _, comment := range got {
            t.Errorf("%s() got %v %q", fn, comment.ID(), comment.Body())
        }
        for _, comment := range want {
            t.Errorf("%s() want %v %q", fn, comment.ID(), comment.Body())
        }
    } else if len2 == 0 {
        if len1 > len2 {
            t.Errorf("%s() = %v comments, want no comments", fn, len1)
        }
    } else {
        g := map[string]*Comment{}
        for _, comment := range got {
            if id := CommentNumber(comment); id != "" {
                g[id] = &comment
            }
        }
        for _, comment := range want {
            if id := CommentNumber(comment); id != "" {
                if g[id] == nil {
                    t.Errorf("%s() missing comment %q", fn, id)
                }
            }
        }
    }
}
