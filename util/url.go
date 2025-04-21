// util/url.go

package util

import (
	"net/url"
)

// ============================================================================
// Exported functions
// ============================================================================

// Indicate whether the two URL objects refer to the same URL.
func SameUrl(u1, u2 *url.URL) bool {
    if u1_nil, u2_nil := (u1 == nil), (u2 == nil); u1_nil || u2_nil {
        return u1_nil && u2_nil
    } else {
        return u1.Path == u2.Path
    }
}
