// Github/rate_limit.go
//
// Functions supporting GitHub rate limit processing.
//
// Note that GitHub has a rate limit on API requests per hour and also a
// secondary rate limits which apply to content-creating requests.
//
// @see https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api

package Github

import (
	"lib.virginia.edu/agita/log"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported variables
// ============================================================================

// The last acquired core rate limit status.
var LastRate github.Rate

// ============================================================================
// Exported functions
// ============================================================================

// Get the last rate limit status.
func RateLimit() github.Rate {
    if LastRate.Limit == 0 {
        GetRateLimit(nil)
    }
    return LastRate
}

// Get the current rate limit status.
//  NOTE: This does not count against secondary rate limits but it *does* count
//  against the primary rate limit.
func GetRateLimit(client *Client) *github.RateLimits {
    if client == nil { client = MainClient() }
    result, _, err := client.ptr.RateLimit.Get(ctx)
    log.ErrorValue(err)
    if (result != nil) && (result.Core != nil) {
        LastRate = *result.Core
        return result
    } else {
        return &github.RateLimits{Core: &LastRate}
    }
}

// ============================================================================
// Internal functions
// ============================================================================

// Get the current rate limit status from a GitHub API response.
//  NOTE: every GitHub API request should be followed by this function.
func extractRateLimit(response *github.Response) {
    if response != nil {
        LastRate = response.Rate
    }
}

