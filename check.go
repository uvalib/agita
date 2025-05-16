// check.go
//
// Validate execution parameters.

package main

import (
	"slices"

	"lib.virginia.edu/agita/Jira"
)

// ============================================================================
// Constants
// ============================================================================

// Command line parameter used to positively indicate that all relevant names
// are to be used.
const ALL = "ALL"

const ALL_REPOS    = ALL
const ALL_PROJECTS = ALL

// ============================================================================
// Functions
// ============================================================================

// Abort unless `names` contains just ALL_REPOS or a non-empty list of GitHub
// repository names.
func ValidateRepoNames(names ...string) []string {
    if count := len(names); count == 0 {
        Abort("must specify %q or a list of repo names", ALL_REPOS)
    } else if (count > 1) && slices.Contains(names, ALL_REPOS) {
        Abort("when %q is given no other names are accepted", ALL_REPOS)
    }
    return RepositoryNames(names...)
}

// Abort unless `names` contains just ALL_PROJECTS or a non-empty list of Jira
// project keys or project issue range bounds.
func ValidateProjectKeys(names ...string) map[string]([]string) {
    if count := len(names); count == 0 {
        Abort("must specify %q or a list of Jira projects", ALL_PROJECTS)
    } else if (count > 1) && slices.Contains(names, ALL_PROJECTS) {
        Abort("when %q is given no other names are accepted", ALL_PROJECTS)
    }
    return Jira.ExpandProjectKeys(names...)
}
