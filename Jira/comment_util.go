// Jira/comment_util.go

package Jira

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

// A generic comment reference.
type CommentArg interface {
    Comment | jira.Comment | *Comment | *jira.Comment
}

// ============================================================================
// Exported functions
// ============================================================================

// Transform a generic comment reference.
func CommentNormalize[T CommentArg](arg T) *jira.Comment {
    switch v := any(arg).(type) {
        case Comment:       return v.ptr
        case jira.Comment:  return &v
        case *Comment:      if v != nil { return v.ptr }
        case *jira.Comment: return v
        default:            panic(fmt.Errorf("unexpected: %v", v))
    }
    return nil
}

// Return the comment identifier for the given comment.
func CommentNumber[T CommentArg](arg T) string {
    if com := CommentNormalize(arg); com == nil {
        return ""
    } else {
        return com.ID
    }
}

// Return the text for the given comment.
func CommentBody[T CommentArg](arg T) string {
    if com := CommentNormalize(arg); com == nil {
        return ""
    } else {
        return com.Body
    }
}

// Indicate whether the two user objects refer to the same comment.
func SameComment[T1, T2 CommentArg](arg1 T1, arg2 T2) bool {
    if nil1, nil2 := util.IsNil(arg1), util.IsNil(arg2); nil1 || nil2 {
        return nil1 && nil2
    } else if num1, num2 := CommentNumber(arg1), CommentNumber(arg2); (num1 != "") || (num2 != "") {
        return num1 == num2
    } else {
        return CommentBody(arg1) == CommentBody(arg2)
    }
}
