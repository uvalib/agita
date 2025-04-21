// Jira/comment_json.go
//
// JSON marshaling of jira.Comment objects.

package Jira

import (
	"encoding/json"
	"fmt"

	"lib.virginia.edu/agita/log"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported methods
// ============================================================================

// Translate a Jira issue comment into JSON.
func (c *Comment) MarshalJSON() ([]byte, error) {
    if m := c.AsCommentMarshal(); m == nil {
        return []byte{}, nil
    } else {
        return json.Marshal(m)
    }
}

// Create a Comment object suitable for JSON marshaling.
func (c *Comment) AsCommentMarshal() *CommentMarshal {
    if noComment(c) {
        return nil
    } else {
        return asCommentMarshal(c.ptr)
    }
}

// ============================================================================
// Exported functions
// ============================================================================

// Reconstitute a Comment object from JSON.
func CommentFromJson(src string) *Comment {
    dst := &jira.Comment{}
    if err := json.Unmarshal([]byte(src), dst); log.ErrorValue(err) == nil {
        return &Comment{ptr: dst}
    } else {
        return nil
    }
}

// ============================================================================
// Internal functions
// ============================================================================

// Create a Comment object suitable for JSON marshaling.
func asCommentMarshal(c *jira.Comment) *CommentMarshal {
    if (c == nil) || (c.ID == "") { return nil }
    use := func(field string) bool { return COMMENT_MARSHAL[field] }
    res := CommentMarshal{}
    if use("ID")           { res.ID           = asStringMarshal(c.ID) }
    if use("Self")         { res.Self         = asStringMarshal(c.Self) }
    if use("Name")         { res.Name         = asStringMarshal(c.Name) }
    if use("Author")       { res.Author       = asUserReference(c.Author) }
    if use("Body")         { res.Body         = asStringMarshal(c.Body) }
    if use("UpdateAuthor") { res.UpdateAuthor = asUserReference(c.UpdateAuthor) }
    if use("Updated")      { res.Updated      = asStringMarshal(c.Updated) }
    if use("Created")      { res.Created      = asStringMarshal(c.Created) }
    if use("Visibility")   { res.Visibility   = asCommentVisibilityMarshal(c.Visibility) }
    if BogusTime(res.Updated) { res.Updated = nil }
    if BogusTime(res.Created) { res.Created = nil }
    return &res
}

// Create an object suitable for JSON marshaling.
func asCommentVisibilityMarshal(arg any) *jira.CommentVisibility {
    var result *jira.CommentVisibility
    switch v := arg.(type) {
        case jira.CommentVisibility:    result = &v
        case *jira.CommentVisibility:   result = v
        default:                        panic(fmt.Errorf("unexpected: %v", v))
    }
    if (result == nil) || (result.Type == "") {
        return nil
    } else {
        return result
    }
}

// ============================================================================
// Exported variables
// ============================================================================

// Controls which jira.Comment fields are converted by asCommentMarshal().
//  @see commentDetails()
var COMMENT_MARSHAL = map[string]bool{
    "ID":           true,
    "Self":         ____,
    "Name":         true,
    "Author":       true,
    "Body":         true,
    "UpdateAuthor": true,
    "Updated":      true,
    "Created":      true,
    "Visibility":   true,
}

// ============================================================================
// Exported types
// ============================================================================

// Used to facilitate jira.Comment JSON marshaling.
type CommentMarshal struct {
    ID           *string                    `json:"id,omitempty" structs:"id,omitempty"`
    Self         *string                    `json:"self,omitempty" structs:"self,omitempty"`
    Name         *string                    `json:"name,omitempty" structs:"name,omitempty"`
    Author       *UserMarshal               `json:"author,omitempty" structs:"author,omitempty"`
    Body         *string                    `json:"body,omitempty" structs:"body,omitempty"`
    UpdateAuthor *UserMarshal               `json:"updateAuthor,omitempty" structs:"updateAuthor,omitempty"`
    Updated      *string                    `json:"updated,omitempty" structs:"updated,omitempty"`
    Created      *string                    `json:"created,omitempty" structs:"created,omitempty"`
    Visibility   *jira.CommentVisibility    `json:"visibility,omitempty" structs:"visibility,omitempty"` // [2)
}
