// re/ops.go
//
// Regular expression shortcuts.

package re

import (
	"regexp"
)

// ============================================================================
// Exported functions
// ============================================================================

// Indicate whether target matches the [Regexp] pattern.
func Match(pattern Pattern, target string) bool {
    return MatchString(pattern, target)
}

// Indicate whether target matches the [Regexp] pattern.
func MatchString(pattern Pattern, target string) bool {
    return regex(pattern).MatchString(target)
}

// Indicate whether target matches the [Regexp] pattern.
func MatchBytes(pattern Pattern, target []byte) bool {
    return regex(pattern).Match(target)
}

// Return a copy of target, replacing matches of the [Regexp] pattern with the
// replacement string repl. Inside repl, $ signs are interpreted as in
// [Regexp.Expand].
func ReplaceAll(pattern Pattern, target, repl string) string {
    return ReplaceAllString(pattern, target, repl)
}

// Return a copy of target, replacing matches of the [Regexp] pattern with the
// replacement string repl. Inside repl, $ signs are interpreted as in
// [Regexp.Expand].
func ReplaceAllString(pattern Pattern, target, repl string) string {
    return regex(pattern).ReplaceAllString(target, repl)
}

// Return a copy of target, replacing matches of the [Regexp] pattern with the
// replacement string repl. Inside repl, $ signs are interpreted as in
// [Regexp.Expand].
func ReplaceAllBytes(pattern Pattern, target, repl []byte) []byte {
    return regex(pattern).ReplaceAll(target, repl)
}

func ReplaceAllFunc(pattern Pattern, target []byte, repl func([]byte) []byte) []byte {
    return regex(pattern).ReplaceAllFunc(target, repl)
}

func Find(pattern Pattern, target string) string {
    return FindString(pattern, target)
}

func FindString(pattern Pattern, target string) string {
    return regex(pattern).FindString(target)
}

func FindBytes(pattern Pattern, target []byte) []byte {
    return regex(pattern).Find(target)
}

func FindAll(pattern Pattern, target string, n int) []string {
    return FindAllString(pattern, target, n)
}

func FindAllString(pattern Pattern, target string, n int) []string {
    return regex(pattern).FindAllString(target, n)
}

func FindAllBytes(pattern Pattern, target []byte, n int) [][]byte {
    return regex(pattern).FindAll(target, n)
}

func FindSubmatches(pattern Pattern, target string) [][]string {
    return FindAllSubmatch(pattern, target, -1)
}

func FindStringSubmatches(pattern Pattern, target string) [][]string {
    return FindAllStringSubmatch(pattern, target, -1)
}

func FindBytesSubmatches(pattern Pattern, target []byte,) [][][]byte {
    return FindAllBytesSubmatch(pattern, target, -1)
}

func FindAllSubmatch(pattern Pattern, target string, n int) [][]string {
    return FindAllStringSubmatch(pattern, target, n)
}

func FindAllStringSubmatch(pattern Pattern, target string, n int) [][]string {
    return regex(pattern).FindAllStringSubmatch(target, n)
}

func FindAllBytesSubmatch(pattern Pattern, target []byte, n int) [][][]byte {
    return regex(pattern).FindAllSubmatch(target, n)
}

// ============================================================================
// Internal functions
// ============================================================================

// Generate a regular expression matching a pattern.
func regex(pattern Pattern) *regexp.Regexp {
    return regexp.MustCompile(withoutSlashes(pattern))
}

// Remove leading and trailing slashes from a pattern which explicitly
// indicates that it is for the purpose of generating a regular expression.
func withoutSlashes(pattern Pattern) Pattern {
    last := len(pattern) - 1
    head, tail := pattern[0:1], pattern[last:]
    if (head == SLASH) && (tail == SLASH) {
        pattern = pattern[1:last]
    }
    return pattern
}
