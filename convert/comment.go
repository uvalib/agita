// convert/comment.go
//
// Conversion of Jira issue comments to GitHub issue comments.

package convert

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"lib.virginia.edu/agita/Github"
	"lib.virginia.edu/agita/Jira"
)

// ============================================================================
// Exported functions
// ============================================================================

// Translate a Jira issue comment into JSON.
func CommentToJson(src Jira.Comment) string {
    if bytes, err := src.MarshalJSON(); log.ErrorValue(err) == nil {
        return string(bytes)
    } else {
        return ""
    }
}

// Translate a Jira issue comment into a Github comment import object.
func Comment(comment Jira.Comment) *Github.CommentImport {
    fld  := map[string]any{}
    note := map[string]any{}
    skip := map[string]bool{}
    add  := func(key string, jiraValue any) {
        if githubValue, use := From(jiraValue); use {
            fld[key] = githubValue
        }
    }

    // If there was no comment body, explicitly show that.
    body := comment.Body()
    if body == "" {
        body = "_(no body)_"
    }

    // Avoid adding an update time if it is the same as the creation time.
    created := comment.Created()
    updated := comment.Updated()
    skip["Updated"] = (created == updated)

    add("CreatedAt",    Github.MakeTime(created))
    add("Body",         body)

    if lines := commentAnnotations(comment, note, skip); len(lines) > 0 {
        notes := strings.Join(lines, "\n")
        fld["Body"] = fmt.Sprintf("%s\n\n%s", notes, fld["Body"])
    }

    return Github.NewCommentImport(fld)
}

// ============================================================================
// Internal functions
// ============================================================================

// Generate lines to annotate the comment body with Jira comment properties
// that have no (updateable) GitHub comment equivalent.
func commentAnnotations(comment Jira.Comment, added map[string]any, skipped map[string]bool) []string {
    res  := []string{}
    tag  := Github.COMMENT_ANNOTATION_TAG
    max  := util.CharCount("UpdateAuthor")
    note := func(key string, jiraValue any) {
        if !skipped[key] {
            if githubValue, use := From(jiraValue); use {
                line := fmt.Sprintf("%s %-*s = %v", tag, max, key, githubValue)
                res = append(res, line)
            }
        }
    }

    // Begin with added notes, which should be exceptional enough that they
    // should appear before other lines.
    for key, val := range added {
        note(key, val)
    }

    // Only include UpdateAuthor if it's different than Author.  For either,
    // annotate with the full name of the Jira user.
    author  := comment.Author()
    updater := comment.UpdateAuthor()
    if author == updater {
        updater = ""
    }

    note("Author",          Jira.AppendFullName(author))
    note("UpdateAuthor",    Jira.AppendFullName(updater))
    note("Updated",         comment.Updated())
    note("Visibility",      comment.Visibility())

    return res
}
