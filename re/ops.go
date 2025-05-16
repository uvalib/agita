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
func Match(target string, pattern Pattern) bool {
    return MatchString(target, pattern)
}

// Indicate whether target matches the [Regexp] pattern.
func MatchString(target string, pattern Pattern) bool {
    return regex(pattern).MatchString(target)
}

// Indicate whether target matches the [Regexp] pattern.
func MatchBytes(target []byte, pattern Pattern) bool {
    return regex(pattern).Match(target)
}

// Return a copy of target, replacing matches of the [Regexp] pattern with the
// replacement string repl. Inside repl, $ signs are interpreted as in
// [Regexp.Expand].
func ReplaceAll(target string, pattern Pattern, repl string) string {
    return ReplaceAllString(target, pattern, repl)
}

// Return a copy of target, replacing matches of the [Regexp] pattern with the
// replacement string repl. Inside repl, $ signs are interpreted as in
// [Regexp.Expand].
func ReplaceAllString(target string, pattern Pattern, repl string) string {
    return regex(pattern).ReplaceAllString(target, repl)
}

// Return a copy of target, replacing matches of the [Regexp] pattern with the
// replacement string repl. Inside repl, $ signs are interpreted as in
// [Regexp.Expand].
func ReplaceAllBytes(target []byte, pattern Pattern, repl []byte) []byte {
    return regex(pattern).ReplaceAll(target, repl)
}

// Return a copy of target in which all matches of the [Regexp] have been
// replaced by the return value of function repl applied to each matched
// substring.
// The replacement returned by repl is substituted directly, without using
// [Regexp.Expand].
//
// The function receives the matching substring and returns the replacement for
// that substring (does not use [Regexp.Expand]).
//
func ReplaceAllFunc(target string, pattern Pattern, repl func(string) string) string {
    return ReplaceAllStringFunc(target, pattern, repl)
}

// Return a copy of target in which all matches of the [Regexp] have been
// replaced by the return value of function repl applied to each matched
// substring.
//
// The function receives the matching substring and returns the replacement for
// that substring (does not use [Regexp.Expand]).
//
func ReplaceAllStringFunc(target string, pattern Pattern, repl func(string) string) string {
    return regex(pattern).ReplaceAllStringFunc(target, repl)
}

func ReplaceAllBytesFunc(target []byte, pattern Pattern, repl func([]byte) []byte) []byte {
    return regex(pattern).ReplaceAllFunc(target, repl)
}

func Find(target string, pattern Pattern) string {
    return FindString(pattern, target)
}

func FindString(target string, pattern Pattern) string {
    return regex(pattern).FindString(target)
}

func FindBytes(target []byte, pattern Pattern) []byte {
    return regex(pattern).Find(target)
}

func FindAll(target string, pattern Pattern, n int) []string {
    return FindAllString(target, pattern, n)
}

func FindAllString(target string, pattern Pattern, n int) []string {
    return regex(pattern).FindAllString(target, n)
}

func FindAllBytes(target []byte, pattern Pattern, n int) [][]byte {
    return regex(pattern).FindAll(target, n)
}

func FindSubmatches(target string, pattern Pattern) [][]string {
    return FindAllSubmatch(target, pattern, -1)
}

func FindStringSubmatches(target string, pattern Pattern) [][]string {
    return FindAllStringSubmatch(target, pattern, -1)
}

func FindBytesSubmatches(target []byte, pattern Pattern) [][][]byte {
    return FindAllBytesSubmatch(target, pattern, -1)
}

func FindAllSubmatch(target string, pattern Pattern, n int) [][]string {
    return FindAllStringSubmatch(target, pattern, n)
}

func FindAllStringSubmatch(target string, pattern Pattern, n int) [][]string {
    return regex(pattern).FindAllStringSubmatch(target, n)
}

func FindAllBytesSubmatch(target []byte, pattern Pattern, n int) [][][]byte {
    return regex(pattern).FindAllSubmatch(target, n)
}

func Split(target string, pattern Pattern, n int) []string {
    return regex(pattern).Split(target, n)
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
