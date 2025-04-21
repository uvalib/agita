// Github/global_sample.go
//
// Sample values used throughout the package.

package Github

import (
	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Exported constants - samples
// ============================================================================

const SAMPLE_ORG  = ORG
const SAMPLE_REPO = "emma"
const SAMPLE_USER = "RayLubinsky"

// ============================================================================
// Exported variables - samples - issues
// ============================================================================

// Pre-existing issues for SAMPLE_REPO and their properties.
var SAMPLE_ISSUE_MAP = map[int]struct{Title, Body string}{
    27: {"Test issue",     re.Pattern(`^This is a fake issue created only to facilitate testing\.\n`)},
    28: {"request title",  "request body"},
    29: {"request title",  "request body"},
//  31: {"imported title", "imported body"},
}

// Pre-existing issues for SAMPLE_REPO.
var SAMPLE_ISSUES = util.MapKeys(SAMPLE_ISSUE_MAP)

// Pre-existing issue for SAMPLE_REPO.
var SAMPLE_ISSUE = 27

// ============================================================================
// Exported variables - samples - comments
// ============================================================================

type sampleCommentMap map[int64]struct{Body string}

// Pre-existing issues for SAMPLE_REPO and their comments.
var SAMPLE_ISSUE_COMMENT_MAP = map[int]sampleCommentMap{
    27: {
        2649639293: {re.Pattern(`^First comment to test issue`)},
        2649639545: {re.Pattern(`^Second comment to test issue`)},
        2652317116: {re.Pattern(`^added comment`)},
        2660375930: {re.Pattern(`^added comment`)},
    //  2663699042: {re.Pattern(`^ORIGINAL COMMENTER "UVADave"`)},
    },
}

// Pre-existing test comments for SAMPLE_ISSUE and their properties.
var SAMPLE_COMMENT_MAP = SAMPLE_ISSUE_COMMENT_MAP[SAMPLE_ISSUE]

// Pre-existing test comments for SAMPLE_ISSUE.
var SAMPLE_COMMENTS = util.MapKeys(SAMPLE_COMMENT_MAP)

// Pre-existing test comment for SAMPLE_ISSUE.
var SAMPLE_COMMENT = SAMPLE_COMMENTS[0]
