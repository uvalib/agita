// convert/issue.go
//
// Conversion of Jira issues to GitHub issues.

package convert

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"lib.virginia.edu/agita/Github"
	"lib.virginia.edu/agita/Jira"
)

// ============================================================================
// Exported functions
// ============================================================================

// Render a Jira issue comment into JSON.
func IssueToJson(src Jira.Issue) string {
    if bytes, err := src.MarshalJSON(); log.ErrorValue(err) == nil {
        return string(bytes)
    } else {
        return ""
    }
}

// Translate a Jira issue into a Github issue import object.
func Issue(issue Jira.Issue) *Github.IssueImport {
    fld  := map[string]any{"Body": ""} // Ensure that "Body" is not nil.
    note := map[string]any{}
    skip := map[string]bool{}
    add  := func(key string, jiraValue any) {
        if githubValue, use := From(jiraValue); use {
            fld[key] = githubValue
        }
    }

    // If the assignee does not have an equivalent GitHub account, then it is
    // added to the annotations.
    assignee := issue.Assignee()
    if githubUser := JiraToGithubUser[assignee]; githubUser != "" {
        assignee = githubUser
    } else {
        note["Assignee"] = Jira.AppendFullName(assignee)
        assignee = ""
    }

    add("Title",        issue.Summary())
    add("Body",         issue.Description())
    add("CreatedAt",    issue.Created())
    add("ClosedAt",     issue.Resolutiondate())
    add("UpdatedAt",    issue.Updated())
    add("Assignee",     assignee)
    add("Labels",       issue.Labels())

    if lines := issueAnnotations(issue, note, skip); len(lines) > 0 {
        notes := strings.Join(lines, "\n")
        fld["Body"] = fmt.Sprintf("%s\n\n%s", notes, fld["Body"])
    }

    return Github.NewIssueImport(fld)
}

// ============================================================================
// Internal functions
// ============================================================================

// Generate lines to annotate the issue body with Jira issue properties that
// have no (updateable) GitHub issue equivalent.
func issueAnnotations(issue Jira.Issue, added map[string]any, skipped map[string]bool) []string {
    res  := []string{}
    tag  := "ORIGINAL JIRA ISSUE"
    max  := util.CharCount("Resolution")
    note := func(key string, jiraValue any) {
        if !skipped[key] {
            if githubValue, use := From(jiraValue); use {
                line := fmt.Sprintf("%s %-*s = %v", tag, max, key, githubValue)
                res = append(res, line)
            }
        }
    }

    // Begin with the original issue key.
    note("Key", issue.Key())

    // Follow with added notes, which should be exceptional enough that they
    // should appear before other lines.
    for key, val := range added {
        note(key, val)
    }

    // Only include Creator if it's different than Reporter.  For either,
    // annotate with the full name of the Jira user.
    creator  := issue.Creator()
    reporter := issue.Reporter()
    if creator == reporter {
        creator = ""
    }

    note("Reporter",    Jira.AppendFullName(reporter))
    note("Creator",     Jira.AppendFullName(creator))
    note("Type",        issue.Type())
    note("Priority",    issue.Priority())
    note("Status",      issue.Status())
    note("Resolution",  issue.Resolution())

    return res
}
