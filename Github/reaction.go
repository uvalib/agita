// Github/reaction.go
//
// Functions supporting GitHub issue and comment reactions.

package Github

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported functions
// ============================================================================

const githubReactionFieldCount = 10

// Render a GitHub Reactions object as a string.
func ReactionString(r *github.Reactions) string {
    switch {
        case r == nil:            return "[nil]"
        case r.TotalCount == nil: return "[missing]"
    }
    array  := make([]string, 0, githubReactionFieldCount)
    addVal := func(key string, val any) {
        array = append(array, fmt.Sprintf("%s: %v", key, val))
    }
    addPos := func(key string, val int) {
        if val > 0 { addVal(key, val) }
    }
    if r.TotalCount != nil { addVal("Total",    *r.TotalCount) }
    if r.PlusOne    != nil { addPos("PlusOne",  *r.PlusOne) }
    if r.MinusOne   != nil { addPos("MinusOne", *r.MinusOne) }
    if r.Laugh      != nil { addPos("Laugh",    *r.Laugh) }
    if r.Confused   != nil { addPos("Confused", *r.Confused) }
    if r.Heart      != nil { addPos("Heart",    *r.Heart) }
    if r.Hooray     != nil { addPos("Hooray",   *r.Hooray) }
    if r.Rocket     != nil { addPos("Rocket",   *r.Rocket) }
    if r.Eyes       != nil { addPos("Eyes",     *r.Eyes) }
    if r.URL        != nil { addVal("URL",      *r.URL) }
    return strings.Join(array, " | ")
}
