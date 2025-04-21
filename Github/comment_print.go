// Github/comment.go
//
// Reporting on comment objects.

package Github

import (
	"fmt"
)

// ============================================================================
// Exported methods
// ============================================================================

// Show details about the instance.
func (c *Comment) Print() {
    id := int64(0)
    if (c.ptr != nil) && (c.ptr.ID != nil) { id = *c.ptr.ID }
    fmt.Printf("\n*** GITHUB Comment %d:\n", id)
    fmt.Println(c.Details())
}
