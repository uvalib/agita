// export.go
//
// Generate JSON from Jira projects.

package main

import (
	"fmt"
	"slices"
	"strings"

	"lib.virginia.edu/agita/convert"
	"lib.virginia.edu/agita/util"

	"lib.virginia.edu/agita/Jira"
)

// ============================================================================
// Constants
// ============================================================================

const PROJECTS_KEY = "Projects"
const ISSUES_KEY   = "Issues"
const COMMENTS_KEY = "Comments"

// ============================================================================
// Functions
// ============================================================================

// Output JSON for all projects and their issues and comments.
//  NOTE: projectKeys must have ALL_PROJECTS or a list of Jira project keys.
//  NOTE: ignores issue range specifications
func ExportAll(projectKeys ...string) {
    projIssues := ValidateProjectKeys(projectKeys...)
    projectKeys = util.MapKeys(projIssues)
    fmt.Print(ProjectsJson(projectKeys...), "\n")
}

// Return JSON for all projects and their issues and comments.
//  NOTE: if projectKeys has ALL_PROJECTS then all projects are used.
func ProjectsJson(projectKeys ...string) string {
    result := fmt.Sprintf("{%q: [\n", PROJECTS_KEY)
    all    := slices.Contains(projectKeys, ALL_PROJECTS)
    items  := Jira.MainClient().GetProjects()
    if len(items) > 0 {
        for _, project := range items {
            if all || slices.Contains(projectKeys, project.Key()) {
                result += ProjectJson(project) + ",\n"
            }
        }
        result = strings.TrimSuffix(result, ",\n") + "\n]"
    }
    return result + "}"
}

// Return JSON for the project object and its issues and comments.
func ProjectJson(project *Jira.Project) string {
    result := convert.ProjectToJson(*project)
    items  := project.Issues()
    if len(items) > 0 {
        result = strings.TrimSuffix(result, "}")
        result += fmt.Sprintf(",\n%q: [\n", ISSUES_KEY)
        for _, item := range items {
            result += IssueJson(item) + ",\n"
        }
        result = strings.TrimSuffix(result, ",\n") + "\n]}"
    }
    return result
}

// Return JSON for the issue object and its comments.
func IssueJson(jiraIssue Jira.Issue) string {
    result := convert.IssueToJson(jiraIssue)
    items  := jiraIssue.Comments()
    if len(items) > 0 {
        result = strings.TrimSuffix(result, "}")
        result += fmt.Sprintf(",\n%q: [\n", COMMENTS_KEY)
        for _, item := range items {
            result += CommentJson(item) + ",\n"
        }
        result = strings.TrimSuffix(result, ",\n") + "\n]}"
    }
    return result
}

// Return JSON for the comment object.
func CommentJson(arg Jira.Comment) string {
    return convert.CommentToJson(arg)
}
