// util/slice.go
//
// Functions supporting slice manipulation.

package util

// ============================================================================
// Exported functions
// ============================================================================

// Return a copy of the slice with nil values removed.
func CompactSlice[T []V, V any](arg T) T {
    result := make(T, 0, len(arg))
    for _, v := range arg {
        if !IsNil(v) {
            result = append(result, v)
        }
    }
    return result
}
