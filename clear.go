// clear.go
//
// Remove all issues and comments from GitHub repositories.

package main

import (
	"fmt"
	"slices"

	"lib.virginia.edu/agita/Github"
)

// ============================================================================
// Functions
// ============================================================================

// Remove all GitHub repository issues and comments.
func ClearAll(names ...string) {
    names  = ValidateRepoNames(names...)
    count := 0
    all   := slices.Contains(names, ALL_REPOS)
    cli   := Github.MainClient()
    org   := Github.ORG
    for _, repo := range cli.GetRepos() {
        name := repo.Name
        if all || slices.Contains(names, name) {
            if num := Github.DeleteIssues(cli, org, name); num == 0 {
                fmt.Printf("Repository %q - no issues removed\n", name)
            } else {
                fmt.Printf("Repository %q - %d issues removed\n", name, num)
            }
            count++
        }
    }
    if count == 1 {
        fmt.Printf("1 repository cleared\n")
    } else {
        fmt.Printf("%d repositories cleared\n", count)
    }
}
