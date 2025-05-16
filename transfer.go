// transfer.go
//
// Generate GitHub issues and comments from Jira issues and comments.

package main

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"lib.virginia.edu/agita/convert"
	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/util"

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

// If true, create "issues-only" project repositories for Jira projects which
// do not map to known existing GitHub repository.
const PROJECT_REPOS = true

// If true, always create "issues-only" project repositories for Jira projects
// regardless of whether they map to a known existing GitHub repository.
const PROJECT_REPOS_ONLY = PROJECT_REPOS && true

// GitHub API accepts no more than 80 content-generating requests per minute.
const REQUESTS_PER_MINUTE = 80

// GitHub API accepts no more than 5000 requests per hour of any kind.
const REQUESTS_PER_HOUR = 5000

// ============================================================================
// Variables
// ============================================================================

// If set, only report output is produced; no GitHub updates are made.
var FakeTransfer bool

// Time the current run was started; used to determine whether GitHub updates
// should be throttled to avoid the primary rate limit.
var TotalTime time.Time

// Request counter used in conjuction with TotalTime.
var TotalCount int

// Time the current batch of requests was started; used to determine whether
// GitHub updates should be throttled to avoid the secondary rate limit.
var BatchTime time.Time

// Request counter used in conjuction with BatchTime.
var BatchCount int

// ============================================================================
// Functions
// ============================================================================

// Generate GitHub issues and comments from Jira issues and comments for all
// Jira projects.
//
// * If PROJECT_REPOS is *false*, action will be limited to those Jira projects
//   which have a known one-to-one mapping with a GitHub repository.
// * If PROJECT_REPOS is *true*, a "project-PROJ" repo will be created if there
//   is no known equivalent GitHub repository.
// * PROJECT_REPOS_ONLY is *true*, even Jira projects with a known equivalent
//   GitHub repository will have a "project-PROJ" repo created.
//
func TransferAll(projectKeys ...string) {
    projIssues := ValidateProjectKeys(projectKeys...)
    projectKeys = util.MapKeys(projIssues)
    all   := slices.Contains(projectKeys, ALL_PROJECTS)
    count := 0
    for _, project := range Jira.MainClient().GetProjects() {
        if proj := project.Key(); all || slices.Contains(projectKeys, proj) {
            repo, projRepo := convert.ProjectToRepo[proj], false
            if PROJECT_REPOS_ONLY {
                repo, projRepo = convert.ProjectRepositoryFor(proj), true
            } else if PROJECT_REPOS && (repo == "") {
                repo, projRepo = convert.RepositoryNameFor(proj), true
            }
            if projRepo && !FakeTransfer {
                if Github.GetProjRepo(Github.MainClient(), repo) == nil {
                    logError("FAILED TO GET PROJECT REPO %q", repo)
                    repo = ""
                }
            }
            if (repo != "") || FakeTransfer {
                minMax := []string{}
                if !all {
                    minMax = projIssues[proj]
                }
                if count == 0 {
                    resetSecondaryRateLimit()
                }
                if TransferProject(project, repo, minMax) {
                    count++
                }
            }
        }
    }
    logSummary("PROJECTS TRANSFERRED: %d", count)
}

// Generate GitHub issues and comments from Jira issues and comments for the
// given Jira project.
func TransferProject(project *Jira.Project, repo string, minMax []string) bool {
    var min, max string
    switch len(minMax) {
        case 0:  min, max = "", ""
        case 1:  min, max = minMax[0], ""
        default: min, max = minMax[0], minMax[1]
    }
    first, last, total := "", "", 0
    for _, issue := range project.GetIssues(min, max) {
        if TransferIssue(issue, repo) {
            key := issue.Key()
            if first == "" { first = key }
            last = key
            total++
        }
    }
    logSummary("%s (%s) ISSUES TRANSFERRED: %d [%q through %q]", project.Key(), project.Name(), total, first, last)
    return total > 0
}

// Generate GitHub issue/comments for a specific Jira issue and its comments.
func TransferIssue(jiraIssue Jira.Issue, repo string) bool {
    // Convert the issue.
    key   := jiraIssue.Key()
    issue := convert.Issue(jiraIssue)
    issue.Body = convertAttachments(issue.Body, key)
    logIssueFields(&jiraIssue, issue)

    // Convert its related comments.
    comments := []*Github.CommentImport{}
    for _, fromJira := range jiraIssue.Comments() {
        toGithub := convert.Comment(fromJira)
        toGithub.Body = convertAttachments(toGithub.Body, key)
        logCommentFields(&jiraIssue, &fromJira, toGithub)
        comments = append(comments, toGithub)
    }

    // Avoid updating GitHub for a fake transfer.
    if FakeTransfer {
        return true
    }

    // Ensure a repository destination was given.
    if repo == "" {
        logError("NO REPO DESTINATION FOR ISSUE %q", key)
        return false
    }
    client := Github.MainClient()

    // Save attachments.
    for _, attach := range jiraIssue.Attachments() {
        checkPrimaryRateLimit()
        file := key + "-" + attach.Filename
        src  := Jira.DownloadAttachment(nil, attach.ID)
        Github.CreateProjAttachment(client, repo, file, src)
    }

    // Create the matching GitHub issue and comments.
    if !checkPrimaryRateLimit() { checkSecondaryRateLimit() }
    Github.ImportIssue(client, Github.ORG, repo, issue, comments...)
    return true
}

// ============================================================================
// Internal functions
// ============================================================================

// Replace inline Jira attachment references with GitHub references to files
// which have been preserved in the project repository's ATTACH_DIR.
func convertAttachments(text, issue Jira.IssueKey) string {
    if issue != "" {
        issue += "-"
    }
    root := "../blob/main/" + Github.ATTACH_DIR
    repl := func(match string) string {
        file := issue + re.ReplaceAll(match, `^!(.*?)(\|.*)?!$`, "$1")
        alt  := ""
        return fmt.Sprintf("![%s](%s/%s?raw=true)", alt, root, file)
    }
    return re.ReplaceAllFunc(text, `!\S[^!\n]+?!`, repl)
}

// Called before every GitHub request to ensure that no more than
// REQUESTS_PER_HOUR are performed.
func checkPrimaryRateLimit() bool {
    if limit := Github.RateLimit(); limit.Remaining <= 1 {
        pause := time.Until(limit.Reset.Time) + (10 * time.Second)
        logWarning("GITHUB PRIMARY RATE LIMIT PAUSE: %v", pause)
        time.Sleep(pause)
        resetSecondaryRateLimit()
        return true
    } else {
        return false
    }
}

// Called before every content-generating GitHub request to ensure that no more
// than REQUESTS_PER_MINUTE are potentially pending.
func checkSecondaryRateLimit() bool {
    if BatchCount++; BatchCount >= REQUESTS_PER_MINUTE - 1 {
        pause := pauseTime(time.Minute, BatchTime)
        logWarning("GITHUB SECONDARY RATE LIMIT PAUSE: %v", pause)
        time.Sleep(pause)
        resetSecondaryRateLimit()
        return true
    } else {
        return false
    }
}

// Initialize secondary rate limit values.
func resetSecondaryRateLimit() {
    BatchTime  = time.Now()
    BatchCount = 0
}

// Determine a time span to wait to avoid a GitHub rate limit.  The time since
// `start` is deducted if it is less than `pause`.
func pauseTime(pause time.Duration, start time.Time) time.Duration {
    if delta := time.Since(start); delta < pause {
        return pause - delta
    } else {
        return pause
    }
}

// ============================================================================
// Internal functions - reporting
// ============================================================================

// Report a problem.
func logError(format string, args ...any) {
    if !strings.HasSuffix(format, "\n") { format += "\n" }
    format = "!!! " + format
    fmt.Printf(format, args...)
}

// Report a condition.
func logWarning(format string, args ...any) {
    if !strings.HasSuffix(format, "\n") { format += "\n" }
    format = "!!! " + format
    fmt.Printf(format, args...)
}

// Report a summary of object transfers.
func logSummary(format string, args ...any) {
    if !LOG_SUMMARIES { return }
    if !strings.HasSuffix(format, "\n") { format += "\n" }
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
        for line := range strings.SplitSeq(lines, "\n") {
            new_lines += leader + line + "\n"
        }
        lines = new_lines
    }
    fmt.Print(lines)
}
