// Jira/issue_ops.go
//
// Operations on jira.Issue objects.

package Jira

import (
	"fmt"
	"reflect"
	"strings"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

// type IssueId  = MapId
type IssueKey = string

// ============================================================================
// Exported variables
// ============================================================================

// Set with searchFields() on module initialization.
var SEARCH_FIELDS []string

// ============================================================================
// Internal functions
// ============================================================================

// Use ISSUE_FIELDS_MARSHAL to determine which fields Search() should return.
//  NOTE: This is for the sake of getting "attachment" returned.
func searchFields() []string {
    Type   := reflect.TypeOf(jira.IssueFields{})
    count  := Type.NumField()
    result := make([]string, 0, count)
	for i := range count {
		field := Type.Field(i)
        name  := field.Name
        if ISSUE_FIELDS_MARSHAL[name] {
            tag := field.Tag.Get("json")
            if tag == "" {
                tag = name
            } else {
                tag = strings.SplitN(tag, ",", 2)[0]
            }
            result = append(result, tag)
        }
	}
    return result
}

// Get all issues for the indicated project.
//  NOTE: may return partial results on error
func getIssues(client *jira.Client, project ProjKey) []jira.Issue {
    return getIssueRange(client, project, "", "")
}

// Get issues for the indicated project, between keys inclusive.
//  NOTE: may return partial results on error
//  NOTE: JQL will fail if a stated issue does not exist
//  NOTE: PROJ-0 and PROJ-1 will be ignored for `minKey`
func getIssueRange(client *jira.Client, project ProjKey, minKey, maxKey IssueKey) (result []jira.Issue) {

    // Build the query based on the indicated issue range.
    jql := fmt.Sprintf("project = %s", project)
    prj := project + "-"
    if min := minKey; min != "" {
        if !strings.HasPrefix(min, prj) {
            min = prj + min
        }
        if (min != prj + "0") && (min != prj + "1") {
            jql += fmt.Sprintf(" AND Key >= %q", min)
        }
    }
    if max := maxKey; max != "" {
        if !strings.HasPrefix(max, prj) {
            max = prj + max
        }
        jql += fmt.Sprintf(" AND Key <= %q", max)
    }
    jql += " ORDER BY Key Asc"

    // Specify issue fields to be returned.
    opt := &jira.SearchOptions{Fields: SEARCH_FIELDS, MaxResults: MAX_PER_PAGE}

    // Get items, possibly across multiple search response pages.
    for last, total := 0, 1; last < total; {
        opt.StartAt = last
		chunk, rsp, err := client.Issue.Search(jql, opt)
		if log.ErrorValue(err) != nil {
            break
		}
		last  = rsp.StartAt + len(chunk)
        total = rsp.Total
		if result == nil {
			result = make([]jira.Issue, 0, total)
		}
		result = append(result, chunk...)
	}
    if result == nil {
        result = []jira.Issue{}
    }
    return
}

// Get the issue with the given issue key.
//  NOTE: returns nil on error
func getIssueByKey(client *jira.Client, key IssueKey) (result *jira.Issue) {
    urlStr := fmt.Sprintf("rest/api/2/issue/%s", key)
    req, err := client.NewRequest("GET", urlStr, nil)
    if log.ErrorValue(err) == nil {
        buffer := jira.Issue{}
        _, err = client.Do(req, &buffer)
        if log.ErrorValue(err) == nil {
            result = &buffer
        }
    }
    return
}

// // Get the issue with the given issue ID.
// //  NOTE: returns nil on error
// func getIssueById(client *jira.Client, id IssueId) (result *jira.Issue) {
//     urlStr := fmt.Sprintf("rest/api/2/issue/%d", id)
//     req, err := client.NewRequest("GET", urlStr, nil)
//     if log.ErrorValue(err) == nil {
//         buffer := jira.Issue{}
//         _, err = client.Do(req, &buffer)
//         if log.ErrorValue(err) == nil {
//             result = &buffer
//         }
//     }
//     return
// }

// ============================================================================
// Internal functions - rendering
// ============================================================================

const jiraIssueFieldCount       = 9
const jiraIssueFieldsFieldCount = 39

// Render a Jira issue object as a multiline string.
//  @see ISSUE_MARSHAL
func issueDetails(i *jira.Issue) string {
    if i == nil { return "" }
    max := jiraIssueFieldCount + jiraIssueFieldsFieldCount - 1
    res := make([]string, 0, max)
    wid := util.CharCount("Fields.Resolutiondate")
    add := func(key string, val any) {
        res = append(res, fmt.Sprintf("%-*s %v", wid, key, val))
    }

    if i.Expand         != ""  { add("Expand",           i.Expand) }
    if i.Self           != ""  { add("Self",             i.Self) }
    if i.ID             != ""  { add("ID",               i.ID) }
    if i.Key            != ""  { add("Key",              i.Key) }
    if i.RenderedFields != nil { add("RenderedFields",  *i.RenderedFields) }
    if i.Changelog      != nil { add("Changelog",       *i.Changelog) }
    if len(i.Transitions)  > 0 { add("Transitions",      i.Transitions) }
    if len(i.Names)        > 0 { add("Names",            i.Names) }

    res = append(res, issueFieldsDetails(i))

    return strings.Join(res, "\n")
}

// Render Jira issue fields as a multiline string.
//  @see ISSUE_FIELDS_MARSHAL
func issueFieldsDetails(i *jira.Issue) string {
    if (i == nil) || (i.Fields == nil) { return "" }
    f   := i.Fields
    res := make([]string, 0, jiraIssueFieldsFieldCount)
    wid := util.CharCount("Resolutiondate")
    add := func(key string, val any) {
        res = append(res, fmt.Sprintf("Fields.%-*s %v", wid, key, val))
    }

    Resolutiondate := TimeString(f.Resolutiondate)
    Created        := TimeString(f.Created)
    DueDate        := DateString(f.Duedate)
    Updated        := TimeString(f.Updated)

    if f.Summary      != ""      { add("Summary",                       f.Summary) }
    if f.Expand       != ""      { add("Expand",                        f.Expand) }
    if f.Type.Name    != ""      { add("Type",                          f.Type.Name) }
    if false                     { add("Project",                       f.Project) }
    if f.Environment  != ""      { add("Environment",                   f.Environment) }
    if f.Resolution   != nil     { add("Resolution",                    f.Resolution.Name) }
    if f.Priority     != nil     { add("Priority",                      f.Priority.Name) }
    if Resolutiondate != NO_TIME { add("Resolutiondate",                Resolutiondate) }
    if Created        != NO_TIME { add("Created",                       Created) }
    if DueDate        != NO_DATE { add("Duedate",                       DueDate) }
    if false                     { add("Watches",                      *f.Watches) }
    if f.Assignee     != nil     { add("Assignee",                      UserLabel(f.Assignee)) }
    if Updated        != NO_TIME { add("Updated",                       Updated) }
    if f.Creator      != nil     { add("Creator",                       UserLabel(f.Creator)) }
    if f.Reporter     != nil     { add("Reporter",                      UserLabel(f.Reporter)) }
    if len(f.Components) > 0     { add("Components",                    f.Components) }
    if f.Status       != nil     { add("Status",                        f.Status.Name) }
    if f.Progress     != nil     { add("Progress",                     *f.Progress) }
    if false                     { add("AggregateProgress",            *f.AggregateProgress) }
    if false                     { add("TimeTracking",                 *f.TimeTracking) }
    if false                     { add("TimeSpent",                     f.TimeSpent) }
    if false                     { add("TimeEstimate",                  f.TimeEstimate) }
    if false                     { add("TimeOriginalEstimate",          f.TimeOriginalEstimate) }
    if false                     { add("Worklog",                      *f.Worklog) }
    if false                     { add("IssueLinks",                    f.IssueLinks) }
    if false                     { add("Comments",                     *f.Comments) }
    if false                     { add("FixVersions",                   f.FixVersions) }
    if false                     { add("AffectsVersions",               f.AffectsVersions) }
    if len(f.Labels) > 0         { add("Labels",                        f.Labels) }
    if len(f.Subtasks) > 0       { add("Subtasks",                      f.Subtasks) }
    if len(f.Attachments) > 0    { add("Attachments",                   f.Attachments) }
    if false                     { add("Epic",                         *f.Epic) }
    if false                     { add("Sprint",                       *f.Sprint) }
    if false                     { add("Parent",                       *f.Parent) }
    if false                     { add("AggregateTimeOriginalEstimate", f.AggregateTimeOriginalEstimate) }
    if false                     { add("AggregateTimeSpent",            f.AggregateTimeSpent) }
    if false                     { add("AggregateTimeEstimate",         f.AggregateTimeEstimate) }
    if false                     { add("Unknowns",                      f.Unknowns) }

    if f.Description  != ""      { add("Description",                   f.Description) }
    if f.Comments     != nil     { add("Comments",                     *f.Comments) }

    return strings.Join(res, "\n")
}

// ============================================================================
// Module initialization
// ============================================================================

// Initialize variables related to Jira issues.
func setupIssue() {
    SEARCH_FIELDS = searchFields()
}
