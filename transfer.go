// transfer.go
//
// Generate GitHub issues and comments from Jira issues and comments.

package main

import (
	"fmt"
	"slices"
	"strings"

	"lib.virginia.edu/agita/convert"

	"lib.virginia.edu/agita/Github"
	"lib.virginia.edu/agita/Jira"
)

// ============================================================================
// Constants
// ============================================================================

// If true, output summary information on transfers.
const LOG_SUMMARIES = true

// If true, output the Jira source object with its GitHub import object.
const LOG_CONVERSIONS = true

// ============================================================================
// Variables
// ============================================================================

// If set, only report output is produced; no GitHub updates are made.
var FakeTransfer bool

// ============================================================================
// Functions
// ============================================================================

// Generate GitHub issues and comments from Jira issues and comments for all
// Jira projects.
//  NOTE: currently is limited to those Jira projects which have a known
//  one-to-one mapping with a GitHub repository.
func TransferAll(projectKeys ...string) {
    projectKeys = ValidateProjectKeys(projectKeys...)
    all   := slices.Contains(projectKeys, ALL_PROJECTS)
    count := 0
    for _, project := range Jira.MainClient().GetProjects() {
        if proj := project.Key(); all || slices.Contains(projectKeys, proj) {
            repo := convert.ProjectToRepo[proj]
            if (repo != "") || FakeTransfer {
                if TransferProject(project, repo) {
                    count++
                }
            }
        }
    }
    logSummary("PROJECTS TRANSFERRED: %d\n", count)
}

// Generate GitHub issues and comments from Jira issues and comments for the
// given Jira project.
func TransferProject(project *Jira.Project, repo string) bool {
    count := 0
    for _, issue := range project.Issues() {
        if TransferIssue(issue, repo) {
            count++
        }
    }
    logSummary("%s (%s) ISSUES TRANSFERRED: %d\n", project.Key(), project.Name(), count)
    return count > 0
}

// Generate GitHub issue/comments for a specific Jira issue and its comments.
func TransferIssue(jiraIssue Jira.Issue, repo string) bool {
    // Convert the issue.
    issue := convert.Issue(jiraIssue)
    logIssueFields(&jiraIssue, issue)

    // Convert its related comments.
    comments := []*Github.CommentImport{}
    for _, fromJira := range jiraIssue.Comments() {
        toGithub := convert.Comment(fromJira)
        comments = append(comments, toGithub)
        logCommentFields(&jiraIssue, &fromJira, toGithub)
    }

    // Create the matching GitHub issue and comments.
    if !FakeTransfer {
        if repo == "" {
            fmt.Printf("!!! NO REPO DESTINATION FOR ISSUE %q", jiraIssue.Key())
            return false
        }
        client := Github.MainClient()
        owner  := Github.ORG
        Github.ImportIssue(client, owner, repo, issue, comments...)
    }
    return true
}

// ============================================================================
// Internal functions - reporting
// ============================================================================

// Report a summary of object transfers.
func logSummary(format string, args ...any) {
    if !LOG_SUMMARIES { return }
    format = "\n\n*** " + format
    fmt.Printf(format, args...)
}

// Report on issue field conversions.
func logIssueFields(jira *Jira.Issue, github *Github.IssueImport) {
    if !LOG_CONVERSIONS { return }
    heading := fmt.Sprintf("Issue %q conversion:", jira.Key())
    logFields(heading, jira, github, 0)
}

// Report on comment field conversions.
func logCommentFields(issue *Jira.Issue, jira *Jira.Comment, github *Github.CommentImport) {
    if !LOG_CONVERSIONS { return }
    heading := fmt.Sprintf("Issue %q Comment conversion:", issue.Key())
    logFields(heading, jira, github, 4)
}

// Any type which has a Details() method.
type Details interface {
    Details() string
}

// Report on field conversions.
func logFields(heading string, jira Details, github Details, indent int) {
    if !LOG_CONVERSIONS { return }
    parts := []string{
        "\n*** " + heading,
        "\n--- JIRA:",
        jira.Details(), 
        "\n--- GITHUB:",
        github.Details(),
    }
    lines := ""
    for _, part := range parts {
        lines += part
        if !strings.HasSuffix(part, "\n") {
            lines += "\n"
        }
    }
    if indent > 0 {
        leader := strings.Repeat(" ", indent)
        new_lines := ""
        for _, line := range strings.Split(lines, "\n") {
            new_lines += leader + line + "\n"
        }
        lines = new_lines
    }
    fmt.Print(lines)
}
