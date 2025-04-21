// Github/label.go
//
// Functions supporting GitHub issue labels.

package Github

import (
	"fmt"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported functions
// ============================================================================

// Render a GitHub Label object as a string.
func LabelString(label *github.Label) string {
    switch {
        case label == nil:      return "[nil]"
        case label.Name == nil: return "[missing]"
        default:                return *label.Name
    }
}

// Render a slice of GitHub Label objects as a slice of strings.
func LabelStrings(labels []*github.Label) []string {
    result := make([]string, 0, len(labels))
    objNil, missing := 0, 0
    for _, label := range labels {
        var name string
        if label == nil {
            objNil++
            name = fmt.Sprintf("[nil-%d]", objNil)
        } else if label.Name == nil {
            missing++
            name = fmt.Sprintf("[missing-%d]", missing)
        } else {
            name = LabelString(label)
        }
        result = append(result, name)
    }
    return result
}
