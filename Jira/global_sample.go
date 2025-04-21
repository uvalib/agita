// Jira/global_sample.go
//
// Sample values used throughout the package.

package Jira

import (
	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Exported constants - samples
// ============================================================================

const SAMPLE_PROJ = "EMMA"
const SAMPLE_USER = "rwl"

// ============================================================================
// Exported variables - samples - issues
// ============================================================================

// Pre-existing issues for SAMPLE_PROJ and their properties.
var SAMPLE_ISSUE_MAP = map[string]struct{Title, Body string}{
//  "EMMA-1":   {"Initial Ruby-on-Rails project",                       re.Pattern(`^An empty Ruby-on-Rails application and GitHub project `)},     // 1 comment
//  "EMMA-2":   {"Docker deployment",                                   re.Pattern(`^Initial setup for deployment to Docker\.`)},                   // 1 comment
//  "EMMA-3":   {"TeamCity build/deploy",                               re.Pattern(`^Build and deploy scripts for EMMA `)},                         // 1 comment
//  "EMMA-4":   {"AWS deployment",                                      re.Pattern(`^Infrastructure setup to allow the EMMA test application `)},   // 1 comment
//  "EMMA-5":   {"Initial monitoring and metrics",                      re.Pattern(`^Initial logic to support extensible reporting `)},             // 1 comment
//  "EMMA-6":   {"API test controller",                                 re.Pattern(`^An initial test of the Bookshare API including:\n`)},          // 1 comment
//  "EMMA-7":   {"Integration with Bookshare authorization",            re.Pattern(`^Implement OAuth2 authorization flow to Bookshare `)},          // 1 comment
//  "EMMA-8":   {"Initial stylesheet support for responsive design",    re.Pattern(`^Establish CSS and SASS definitions `)},                        // 1 comment
//  "EMMA-10":  {"Basic title search",                                  re.Pattern(`^Provide an input for performing a keyword search `)},          // 1 comment
//  "EMMA-11":  {"Artifact format URL parameter",                       re.Pattern(`^Rails handles the URL "format" parameter `)},                  // 1 comment
//  "EMMA-12":  {"Browser compatibility adjustments",                   re.Pattern(`^Some CSS/SCSS adjustments are needed `)},                      // 1 comment
//  "EMMA-25":  {"Upgrade to Rails 6",                                  re.Pattern(`^Rails 6 has been released `)},                                 // 1 comment
//  "EMMA-32":  {"Automated Bookshare API compliance tests",            re.Pattern(`^The version 2 of the Bookshare API \(used by EMMA\) `)},       // 1 comment
    "EMMA-33":  {"Updates for Bookshare API v5.6.12",                   re.Pattern(`^Accommodate deltas from v.5.6.10:\n`)},                        // 2 comments
//  "EMMA-34":  {"Upgrade to Ruby 2.6.3",                               re.Pattern(`^Benchmarks indicate that Ruby 2.6 is a smidge better `)},      // 1 comment
//  "EMMA-35":  {"Exception handling",                                  re.Pattern(`^Currently exceptions that occur when `)},                      // 1 comment
//  "EMMA-40":  {"Federated search",                                    re.Pattern(`^Prepare a test controller to track development `)},            // 1 comment
    "EMMA-44":  {"File upload",                                         ""},                                                                        // 1 comment
//  "EMMA-46":  {"Multi-select search facets",                          re.Pattern(`^Support requesting of search results `)},                      // 1 comment
//  "EMMA-47":  {"Update to Ruby 2.7",                                  ""},                                                                        // 1 comment
//  "EMMA-48":  {"AWS S3 storage",                                      re.Pattern(`^Reconfigure the file upload process `)},                       // 1 comment
//  "EMMA-49":  {"File download",                                       re.Pattern(`^Support download of content items `)},                         // 1 comment
//  "EMMA-51":  {"Index ingest",                                        re.Pattern(`^Make use of the Benetech "EMMA Federated Ingest" service`)},   // 1 comment
//  "EMMA-63":  {"Help text updates",                                   ""},                                                                        // 1 comment
//  "EMMA-75":  {`Add OAuth2 token revocation to "Sign Out"`,           re.Pattern(`^Benetech has modified their OAuth2 service`)},                 // 2 comments
//  "EMMA-92":  {"Improved error handling and reporting",               re.Pattern(`^Improvements to the infrastructure are needed `)},             // 1 comment
//  "EMMA-122": {"Matomo Analytics",                                    re.Pattern(`^Add analytics:\n`)},                                           // 0 comment; table
    "EMMA-131": {"IA download redesign",                                re.Pattern(`^Internet Archive has a \[new API|https://docs.google.com/document/d/1mcZwhvTGUmkSOcS7UpgH8Bni7e9dQ03jTiOirEIWK_U/view\]  `)}, // 1 comment
//  "EMMA-133": {"Inactive status",                                     re.Pattern(`^Implement the mechanics of Inactive status `)},                // 1 comment
}

// Pre-existing issues for SAMPLE_PROJ.
var SAMPLE_ISSUES = util.MapKeys(SAMPLE_ISSUE_MAP)

// Pre-existing issue for SAMPLE_PROJ.
var SAMPLE_ISSUE = "EMMA-33"

// ============================================================================
// Exported variables - samples - comments
// ============================================================================

type sampleCommentMap map[int]struct{Body string}

// Pre-existing issues for SAMPLE_PROJ and their comments.
var SAMPLE_ISSUE_COMMENT_MAP = map[string]sampleCommentMap{
    "EMMA-33": {
        1148864: {re.Pattern(`^Adjustment for Bookshare API v5.6.13:\n`)},
        1148865: {re.Pattern(`^The new API requests and data structures have been added.\n`)},
    },
}

// Pre-existing test comments for SAMPLE_ISSUE and their properties.
var SAMPLE_COMMENT_MAP = SAMPLE_ISSUE_COMMENT_MAP[SAMPLE_ISSUE]

// Pre-existing test comments for SAMPLE_ISSUE.
var SAMPLE_COMMENTS = util.MapKeys(SAMPLE_COMMENT_MAP)

// Pre-existing test comment for SAMPLE_ISSUE.
var SAMPLE_COMMENT = SAMPLE_COMMENTS[0]
