// re/flags.go
//
// Regular expression flag manipulation.

package re

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Exported constants
// ============================================================================

const REGEX_FLAGS = "imsU"

const FLAGS_START = "(?"
const FLAGS_END   = ")"

const FLAG_PREFIX = FLAGS_START

// ============================================================================
// Exported types
// ============================================================================

type FlagMap = map[string]bool

type StringLike interface { ~byte|~rune|~string }

// ============================================================================
// Exported functions
// ============================================================================

// Indicate whether the argument appears to have regular expression flags.
func HasRegexFlags(str string) bool {
    return str[0:2] == FLAG_PREFIX
}

// Sets the flag prefix on pattern.
// If `flags` is empty, pattern will have all flags removed.
func FlagPrefix[T StringLike](pattern Pattern, flags []T) Pattern {
    isSet := flagMap(flags)
    return flagPattern(pattern, isSet)
}

// Add the flag prefix to pattern.
// If `flags` is empty, pattern will be returned without change.
func AddFlags[T StringLike](pattern Pattern, flags []T) Pattern {
    isSet := flagMap(flags)
    if HasRegexFlags(pattern) {
        flagsEnd := strings.Index(pattern, FLAGS_END)
        existing := flagMap([]byte(pattern[2:flagsEnd]))
        isSet = util.MapMerge(existing, isSet)
    }
    return flagPattern(pattern, isSet)
}

// ============================================================================
// Internal functions
// ============================================================================

// Set flags in pattern based on the given settings.
func flagPattern(pattern Pattern, isSet FlagMap) string {
    flagChars := ""
    for _, char := range REGEX_FLAGS {
        if flag := flagString(char); isSet[flag] {
            flagChars += flag
        }
    }
    if flagChars != "" {
        flagChars = FLAGS_START + flagChars + FLAGS_END
    }
    return flagChars + stripFlags(pattern)
}

// Remove "(?...)" flags from a pattern.
func stripFlags(pattern Pattern) Pattern {
    if HasRegexFlags(pattern) {
        afterFlags := strings.Index(pattern, FLAGS_END) + len(FLAGS_END)
        return pattern[afterFlags:]
    } else {
        return pattern
    }
}

// Process `flags` into a FlagMap.
//  NOTE: Does not validate against REGEX_FLAGS.
func flagMap[T StringLike](flags []T) FlagMap {
    isSet := FlagMap{}
    for _, char := range flags {
        flag := flagString(char)
        if (flag != "") && strings.Contains(REGEX_FLAGS, flag) {
            isSet[flag] = true
        }
    }
    return isSet
}

// Convert `flagChar` into the form of a FlagMap key.
//  NOTE: Does not validate against REGEX_FLAGS.
func flagString[T StringLike](flagChar T) string {
    switch v := any(flagChar).(type) {
        case byte:   return string(v)
        case rune:   return string(v)
        case string: return string(v)
        default:     panic(fmt.Errorf("unexpected type: %v", v))
    }
}
