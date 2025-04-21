// Github/comment_util.go

package Github

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// A generic comment reference.
type CommentArg interface {
    Comment | github.IssueComment | *Comment | *github.IssueComment
}

// ============================================================================
// Exported functions
// ============================================================================

// Transform a generic comment reference.
func CommentNormalize[T CommentArg](arg T) *github.IssueComment {
    switch v := any(arg).(type) {
        case Comment:              return v.ptr
        case github.IssueComment:  return &v
        case *Comment:             if v != nil { return v.ptr }
        case *github.IssueComment: return v
        default:                   panic(fmt.Errorf("unexpected: %v", v))
    }
    return nil
}

// Return the comment identifier for the given comment.
func CommentNumber[T CommentArg](arg T) int64 {
    if com := CommentNormalize(arg); (com == nil) || (com.ID == nil) {
        return 0
    } else {
        return *com.ID
    }
}

// Return the text for the given comment.
func CommentBody[T CommentArg](arg T) string {
    if com := CommentNormalize(arg); (com == nil) || (com.Body == nil) {
        return ""
    } else {
        return *com.Body
    }
}

// Indicate whether the two user objects refer to the same comment.
func SameComment[T1, T2 CommentArg](arg1 T1, arg2 T2) bool {
    if nil1, nil2 := util.IsNil(arg1), util.IsNil(arg2); nil1 || nil2 {
        return nil1 && nil2
    } else if num1, num2 := CommentNumber(arg1), CommentNumber(arg2); (num1 != 0) || (num2 != 0) {
        return num1 == num2
    } else {
        return CommentBody(arg1) == CommentBody(arg2)
    }
}
