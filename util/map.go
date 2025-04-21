// util/map.go
//
// Functions supporting map manipulation.

package util

import (
	"fmt"
	"maps"
	"slices"
)

// ============================================================================
// Exported functions
// ============================================================================

// Return a slice with the keys of the given map in random order.
func MapKeys[Map ~map[K]V, K comparable, V any](arg Map) []K {
    return slices.Collect(maps.Keys(arg))
}

// Return a slice with the values from the given map in random order.
func MapValues[Map ~map[K]V, K comparable, V any](arg Map) []V {
    return slices.Collect(maps.Values(arg))
}

// Create a copy of `arg1` with entries from `arg2` added.
func MapMerge[Map ~map[K]V, K comparable, V any](arg1, arg2 Map) Map {
    result := maps.Clone(arg1)
    maps.Copy(result, arg2)
    return result
}

// Generate a map with keys and values swapped.
//  NOTE: assumes that all original values are unique
func MapInvert[K comparable, V comparable](arg map[K]V) map[V]K {
    result := make(map[V]K, len(arg))
    for k, v := range arg {
        result[v] = k
    }
    return result
}

// Return a copy of the map with nil values removed.
func MapCompact[Map ~map[K]V, K comparable, V any](arg Map) Map {
    result := make(Map, len(arg))
    for k, v := range arg {
        if !IsNil(v) {
            result[k] = v
        }
    }
    return result
}

// ============================================================================
// Exported functions - reporting
// ============================================================================

// Output a map sorted by key and aligned on value.
func PrintSortedMap[K string, V any](arg map[K]V) {
    keys := make([]K, 0, len(arg))
    max  := 0
    for k := range arg {
        keys = append(keys, k)
        if max < len(k) {
            max = len(k)
        }
    }
    slices.Sort(keys)
    fmt.Println("")
    for _, k := range keys {
        fmt.Printf("\t%-*s - %v\n", max, k, arg[k])
    }
}
