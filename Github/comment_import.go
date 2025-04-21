// Github/comment_import.go

package Github

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// Application wrapper for the GitHub object used to import comments.
type CommentImport struct {
    github.Comment
}

// ============================================================================
// Exported functions
// ============================================================================

// Generate a wrapper for a Github comment import object from a table of field
// names and values
//  NOTE: never returns nil
func NewCommentImport(fields map[string]any) *CommentImport {
    com := github.Comment{}
    for key, val := range fields {
        var s string
        var t Time
        switch v := val.(type) {
            case string:    s = v
            case Time:      t = v
            case *string:   s = *v
            case *Time:     t = *v
        }
        switch key {
            case "CreatedAt":   com.CreatedAt   = github.Ptr(t)
            case "Body":        com.Body        = s
        }
    }
    return &CommentImport{Comment: com}
}

// Get the native GitHub comment import objects.
func Comments(comments []*CommentImport) []*github.Comment {
    result := []*github.Comment{}
    for _, comment := range comments {
        result = append(result, &comment.Comment)
    }
    return result
}

// ============================================================================
// Exported members - rendering
// ============================================================================

const githubCommentImportFieldCount = 2

// Render details about the instance.
func (c *CommentImport) Details() string {
    res := make([]string, 0, githubCommentImportFieldCount)
    max := util.CharCount("CreatedAt")
    add := func(key string, val any) {
        res = append(res, fmt.Sprintf("%-*s %v", max, key, val))
    }

    if c.CreatedAt != nil   { add("CreatedAt",  *c.CreatedAt) }
    if c.Body      != ""    { add("Body",        c.Body) }

    return strings.Join(res, "\n")
}
