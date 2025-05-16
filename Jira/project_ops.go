// Jira/project_ops.go
//
// Operations on jira.Project objects.

package Jira

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/util"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

type ProjId  = int
type ProjKey = string

// ============================================================================
// Exported variables
// ============================================================================

var ProjectByKey map[ProjKey]ProjId

// ============================================================================
// Exported functions - identity
// ============================================================================

// Perform a case-insensitive lookup into ProjectByKey.
func ProjectKeyToId(key ProjKey) ProjId {
    return ProjectByKey[strings.ToUpper(key)]
}

// Treat each name as either a ProjKey or an IssueKey.
//
// * If PROJ appears alone then that represents (implicitly) all of the issues
//   for Jira project PROJ.
// * If PROJ-num appears then that represents Jira project PROJ starting at
//   issue key "PROJ-num".
// * If PROJ-min and PROJ-max appear then that represents Jira project PROJ
//   issues from "PROJ-min" through "PROJ-max".
//
func ExpandProjectKeys(names ...string) map[string]([]string) {
    result := map[string]([]string){}
    for _, name := range names {
        project := name
        if re.Match(name, `^[A-Z]+-\d+$`) {
            project = re.ReplaceAll(name, `^([A-Z]+)-\d+$`, "$1")
            issue  := name
            if minMax := result[project]; len(minMax) > 0 {
                result[project] = append(minMax, issue)
            } else {
                result[project] = []string{issue}
            }
        } else {
            result[project] = []string{}
        }
    }
    return result
}

// ============================================================================
// Internal functions
// ============================================================================

// Get all projects for the Jira referenced by the client.
//  NOTE: all returned elements are non-nil
func getProjects(client *jira.Client) []*jira.Project {
    result   := []*jira.Project{}
    urlStr   := "rest/api/2/project"
    req, err := client.NewRequest("GET", urlStr, nil)
    if log.ErrorValue(err) == nil {
        _, err = client.Do(req, &result)
        log.ErrorValue(err)
    }
    return result
}

// Get the project with the given project key.
//  NOTE: returns nil on error
func getProjectByKey(client *jira.Client, key ProjKey) *jira.Project {
    projId := strconv.Itoa(ProjectKeyToId(key))
    project, _, err := client.Project.Get(projId)
    log.ErrorValue(err)
    return project
}

// // Get the project with the given identifier.
// //  NOTE: returns nil on error
// func getProjectById(client *jira.Client, id ProjId) *jira.Project {
//     projId := strconv.Itoa(id)
//     project, _, err := client.Project.Get(projId)
//     log.ErrorValue(err)
//     return project
// }

// ============================================================================
// Internal functions - initialization
// ============================================================================

// Generate a mapping of project ID to project key of all projects.
func projectByKey() map[ProjKey]ProjId {
    items  := NewClient().GetProjects()
    result := make(map[ProjKey]ProjId, len(items))
    for _, project := range items {
        if projectId, err := strconv.Atoi(project.ptr.ID); err == nil {
            id  := ProjId(projectId)
            key := ProjKey(project.ptr.Key)
            result[key] = id
        }
    }
    return result
}

// ============================================================================
// Internal methods - rendering
// ============================================================================

const jiraProjectFieldCount = 16

// Render a Jira project object as a multiline string.
//  @see PROJECT_MARSHAL
func projectDetails(p *jira.Project) string {
    if p == nil { return "" }
    ary := make([]string, 0, jiraProjectFieldCount)
    wid := util.CharCount("ProjectCategory")
    add := func(key string, val any) {
        ary = append(ary, fmt.Sprintf("%-*s %v", wid, key, val))
    }

    var Lead string
    if testing.Testing() {
        Lead = Account(&p.Lead)
    } else {
        Lead = UserLabel(&p.Lead)
    }
    Category := p.ProjectCategory.Name

    if p.Expand       != "" { add("Expand",             p.Expand) }
    if p.Self         != "" { add("Self",               p.Self) }
    if p.ID           != "" { add("ID",                 p.ID) }
    if p.Key          != "" { add("Key",                p.Key) }
    if p.Description  != "" { add("Description",        p.Description) }
    if Lead           != "" { add("Lead",               Lead) }
    if false                { add("Components",         p.Components) }
    if false                { add("IssueTypes",         p.IssueTypes) }
    if p.URL          != "" { add("URL",                p.URL) }
    if p.Email        != "" { add("Email",              p.Email) }
    if p.AssigneeType != "" { add("AssigneeType",       p.AssigneeType) }
    if p.Name         != "" { add("Name",               p.Name) }
    if false                { add("Roles",              p.Roles) }
    if false                { add("AvatarUrls",         p.AvatarUrls) }
    if Category       != "" { add("ProjectCategory",    Category) }

    return strings.Join(ary, "\n")
}

// ============================================================================
// Module initialization
// ============================================================================

// Initialize variables related to Jira projects.
func setupProject() {
    if ProjectByKey == nil {
        ProjectByKey = projectByKey()
    }
}
