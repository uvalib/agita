// re/type.go
//
// Type wrapping regexp with tailored logic.

package re

import (
	"regexp"
	"slices"
	"strings"
)

// ============================================================================
// Exported types
// ============================================================================

type Regex struct {
    pattern Pattern
    simple  bool
    exact   bool
    regex   *regexp.Regexp
}

// ============================================================================
// Exported functions
// ============================================================================

// Create a new instace.  Assignment of the `regex` field is deferred until it
// is actually needed.
func New(pattern Pattern) *Regex {
    exact  := (pattern == "")
    simple := exact || IsSimple(pattern)
    return &Regex{pattern, simple, exact, nil}
}

// ============================================================================
// Exported methods
// ============================================================================

func (r *Regex) Simple() bool {
    return r.simple
}

func (r *Regex) Exact() bool {
    return r.exact
}

func (r *Regex) String() string {
    return r.pattern
}

func (r *Regex) Bytes() []byte {
    return []byte(r.pattern)
}

// Indicate whether target matches the instance pattern.
func (r *Regex) Match(target string) bool {
    return r.MatchString(target)
}

// Indicate whether target matches the instance pattern.
func (r *Regex) MatchString(target string) bool {
    switch {
        case r.exact:   return target == r.pattern
        case r.simple:  return strings.Contains(target, r.pattern)
        default:        return r.regexp().MatchString(target)
    }
}

// Indicate whether target matches the instance pattern.
func (r *Regex) MatchBytes(target []byte) bool {
    switch {
        case r.exact:   return slices.Equal(target, []byte(r.pattern))
        default:        return r.regexp().Match(target)
    }
}

// Return a copy of target, replacing matches of the instance pattern with the
// replacement string repl. Inside repl, $ signs are interpreted as in
// [Regexp.Expand].
func (r *Regex) ReplaceAll(target, repl string) string {
    return r.ReplaceAllString(target, repl)
}

// Return a copy of target, replacing matches of the instance pattern with the
// replacement string repl. Inside repl, $ signs are interpreted as in
// [Regexp.Expand].
func (r *Regex) ReplaceAllString(target, repl string) string {
    return r.regexp().ReplaceAllString(target, repl)
}

// Return a copy of target, replacing matches of the instance pattern with the
// replacement string repl. Inside repl, $ signs are interpreted as in
// [Regexp.Expand].
func (r *Regex) ReplaceAllBytes(target, repl []byte) []byte {
    return r.regexp().ReplaceAll(target, repl)
}

// Return a copy of target, replacing matches of the instance pattern with the
// return value of function repl applied to each matched substring.
//
// The function receives the matching substring and returns the replacement for
// that substring (does not use [Regexp.Expand]).
//
func (r *Regex) ReplaceAllFunc(target string, repl func(string) string) string {
    return r.ReplaceAllStringFunc(target, repl)
}

// Return a copy of target, replacing matches of the instance pattern with the
// return value of function repl applied to each matched substring.
//
// The function receives the matching substring and returns the replacement for
// that substring (does not use [Regexp.Expand]).
//
func (r *Regex) ReplaceAllStringFunc(target string, repl func(string) string) string {
    return r.regexp().ReplaceAllStringFunc(target, repl)
}

func (r *Regex) ReplaceAllBytesFunc(target []byte, repl func([]byte) []byte) []byte {
    return r.regexp().ReplaceAllFunc(target, repl)
}

func (r *Regex) Find(target string) string {
    return r.FindString(target)
}

func (r *Regex) FindString(target string) string {
    return r.regexp().FindString(target)
}

func (r *Regex) FindBytes(target []byte) []byte {
    return r.regexp().Find(target)
}

func (r *Regex) FindAll(target string, n int) []string {
    return r.FindAllString(target, n)
}

func (r *Regex) FindAllString(target string, n int) []string {
    return r.regexp().FindAllString(target, n)
}

func (r *Regex) FindAllBytes(target []byte, n int) [][]byte {
    return r.regexp().FindAll(target, n)
}

func (r *Regex) FindSubmatches(target string) [][]string {
    return r.FindAllSubmatch(target, -1)
}

func (r *Regex) FindStringSubmatches(target string) [][]string {
    return r.FindAllStringSubmatch(target, -1)
}

func (r *Regex) FindBytesSubmatches(target []byte,) [][][]byte {
    return r.FindAllBytesSubmatch(target, -1)
}

func (r *Regex) FindAllSubmatch(target string, n int) [][]string {
    return r.FindAllStringSubmatch(target, n)
}

func (r *Regex) FindAllStringSubmatch(target string, n int) [][]string {
    return r.regexp().FindAllStringSubmatch(target, n)
}

func (r *Regex) FindAllBytesSubmatch(target []byte, n int) [][][]byte {
    return r.regexp().FindAllSubmatch(target, n)
}

func (r *Regex) Split(target string, n int) []string {
    return r.regexp().Split(target, n)
}

// ============================================================================
// Internal methods
// ============================================================================

func (r *Regex) regexp() *regexp.Regexp {
    if r.regex == nil {
        r.regex = regex(r.pattern)
    }
    return r.regex
}
