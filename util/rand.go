// util/rnd.go
//
// Random value definitions.

package util

import (
	"fmt"
	"math"
	"math/rand"
)

// ============================================================================
// Internal constants
// ============================================================================

const hexDigits = 6
const hexBase   = 16

const bigDigits = 13
const bigBase   = 36

// ============================================================================
// Internal variables
// ============================================================================

var hexRange = int(math.Pow(hexBase, hexDigits))

// ============================================================================
// Exported functions
// ============================================================================

// Render a 6-digit string of hexadecimal characters.
func HexRand() string {
    number := Base(rand.Intn(hexRange), hexBase)
    return fmt.Sprintf("%0*s", hexDigits, number)
}

// Render a 13-digit string of base36 characters.
func BigRand() string {
    number := Base(rand.Int(), bigBase)
    return fmt.Sprintf("%0*s", bigDigits, number)
}

// Generate a random number string for use in constructing random identifiers.
func RandomId() string {
    return HexRand()
}

// Modify str by appending a random number string.
func Randomize(str string) string {
    return AppendNote(str,  RandomId())
}
