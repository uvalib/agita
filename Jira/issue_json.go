// Jira/issue_json.go
//
// JSON marshaling of jira.Issue objects.

package Jira

import (
	"encoding/json"
	"fmt"

	"lib.virginia.edu/agita/log"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported methods
// ============================================================================

// Translate a Jira issue into JSON.
func (i *Issue) MarshalJSON() ([]byte, error) {
    if m := i.AsIssueMarshal(); m == nil {
        return []byte{}, nil
    } else {
        return json.Marshal(m)
    }
}

// Create an Issue object suitable for JSON marshaling.
func (i *Issue) AsIssueMarshal() *IssueMarshal {
    if noIssue(i) {
        return nil
    } else {
        return asIssueMarshal(i.ptr)
    }
}

// ============================================================================
// Exported functions
// ============================================================================

// Reconstitute an Issue object from JSON.
func IssueFromJson(src string) *Issue {
    dst := &jira.Issue{}
    if err := json.Unmarshal([]byte(src), dst); log.ErrorValue(err) == nil {
        return &Issue{ptr: dst}
    } else {
        return nil
    }
}

// ============================================================================
// Internal functions
// ============================================================================

// Create an Issue object suitable for JSON marshaling.
func asIssueMarshal(i *jira.Issue) *IssueMarshal {
    if (i == nil) || (i.Key == "") { return nil }
    use := func(field string) bool { return ISSUE_MARSHAL[field] }
    res := IssueMarshal{Fields: &IssueFieldsMarshal{}}
    if use("Key")            { res.Key            = asStringMarshal(i.Key) }
    if use("ID")             { res.ID             = asStringMarshal(i.ID) }
    if use("Expand")         { res.Expand         = asStringMarshal(i.Expand) }
    if use("Self")           { res.Self           = asStringMarshal(i.Self) }
    if use("Fields")         { res.Fields         = asIssueFieldsMarshal(i.Fields) }
    if use("RenderedFields") { res.RenderedFields = i.RenderedFields }
    if use("Changelog")      { res.Changelog      = i.Changelog }
    if use("Transitions")    { res.Transitions    = i.Transitions }
    if use("Names")          { res.Names          = i.Names }
    return &res
}

// Create an IssueFields object suitable for JSON marshaling.
func asIssueFieldsMarshal[T jira.IssueFields|*jira.IssueFields](arg T) *IssueFieldsMarshal {
    var src *jira.IssueFields
    switch v := any(arg).(type) {
        case jira.IssueFields:  src = &v
        case *jira.IssueFields: if src = v; src == nil { return nil }
        default:                panic(fmt.Errorf("unexpected: %v", v))
    }
    use := func(field string) bool { return ISSUE_FIELDS_MARSHAL[field] }
    res := IssueFieldsMarshal{}
    if use("Expand")                        { res.Expand                        = asStringMarshal(src.Expand) }
    if use("Type")                          { res.Type                          = asIssueTypeReference(src.Type) }
    if use("Project")                       { res.Project                       = asProjectReference(src.Project) }
    if use("Environment")                   { res.Environment                   = asStringMarshal(src.Environment) }
    if use("Resolution")                    { res.Resolution                    = asResolutionReference(src.Resolution) }
    if use("Priority")                      { res.Priority                      = asPriorityReference(src.Priority) }
    if use("Resolutiondate")                { res.Resolutiondate                = asTimeMarshal(src.Resolutiondate) }
    if use("Created")                       { res.Created                       = asTimeMarshal(src.Created) }
    if use("Duedate")                       { res.Duedate                       = asDateMarshal(src.Duedate) }
    if use("Watches")                       { res.Watches                       = src.Watches }
    if use("Assignee")                      { res.Assignee                      = asUserReference(src.Assignee) }
    if use("Updated")                       { res.Updated                       = asTimeMarshal(src.Updated) }
    if use("Description")                   { res.Description                   = asStringMarshal(src.Description) }
    if use("Summary")                       { res.Summary                       = asStringMarshal(src.Summary) }
    if use("Creator")                       { res.Creator                       = asUserReference(src.Creator) }
    if use("Reporter")                      { res.Reporter                      = asUserReference(src.Reporter) }
    if use("Components")                    { res.Components                    = src.Components }
    if use("Status")                        { res.Status                        = asStatusReference(src.Status) }
    if use("Progress")                      { res.Progress                      = asProgressMarshal(src.Progress) }
    if use("AggregateProgress")             { res.AggregateProgress             = asProgressMarshal(src.AggregateProgress) }
    if use("TimeTracking")                  { res.TimeTracking                  = src.TimeTracking }
    if use("TimeSpent")                     { res.TimeSpent                     = asIntMarshal(src.TimeSpent) }
    if use("TimeEstimate")                  { res.TimeEstimate                  = asIntMarshal(src.TimeEstimate) }
    if use("TimeOriginalEstimate")          { res.TimeOriginalEstimate          = asIntMarshal(src.TimeOriginalEstimate) }
    if use("Worklog")                       { res.Worklog                       = src.Worklog }
    if use("IssueLinks")                    { res.IssueLinks                    = src.IssueLinks }
    if use("Comments")                      { res.Comments                      = asCommentsMarshal(src.Comments) }
    if use("FixVersions")                   { res.FixVersions                   = src.FixVersions }
    if use("AffectsVersions")               { res.AffectsVersions               = src.AffectsVersions }
    if use("Labels")                        { res.Labels                        = src.Labels }
    if use("Subtasks")                      { res.Subtasks                      = asSubtasksMarshal(src.Subtasks) }
    if use("Attachments")                   { res.Attachments                   = src.Attachments }
    if use("Epic")                          { res.Epic                          = src.Epic }
    if use("Sprint")                        { res.Sprint                        = src.Sprint }
    if use("Parent")                        { res.Parent                        = src.Parent }
    if use("AggregateTimeOriginalEstimate") { res.AggregateTimeOriginalEstimate = asIntMarshal(src.AggregateTimeOriginalEstimate) }
    if use("AggregateTimeSpent")            { res.AggregateTimeSpent            = asIntMarshal(src.AggregateTimeSpent) }
    if use("AggregateTimeEstimate")         { res.AggregateTimeEstimate         = asIntMarshal(src.AggregateTimeEstimate) }
    return &res
}

// Create a Progress object suitable for JSON marshaling.
func asProgressMarshal(arg any) *jira.Progress {
    var progress jira.Progress
    switch v := arg.(type) {
        case int:               progress.Percent = v
        case jira.Progress:     progress = v
        case *int:              if v != nil { progress.Percent = *v }
        case *jira.Progress:    if v != nil { progress = *v }
        default:                panic(fmt.Errorf("unexpected: %v", v))
    }
    if progress.Percent == 0 {
        return nil
    } else {
        return &progress
    }
}

// Create an object suitable for JSON marshaling of an array of comments.
func asCommentsMarshal(src *jira.Comments) *CommentsMarshal {
    if src != nil {
        if size := len(src.Comments); size > 0 {
            comments := make([]*CommentMarshal, size)
            for _, comment := range src.Comments {
                comments = append(comments, asCommentMarshal(comment))
            }
            return &CommentsMarshal{Comments: comments}
        }
    }
    return nil
}

// Create an object suitable for JSON marshaling of an array of subtasks.
func asSubtasksMarshal(src []*jira.Subtasks) []*SubtasksMarshal {
    if size := len(src); size > 0 {
        subtasks := make([]*SubtasksMarshal, 0, size)
        for _, subtask := range src {
            subtasks = append(subtasks, &SubtasksMarshal{Key: &subtask.Key})
        }
        return subtasks
    }
    return []*SubtasksMarshal{}
}

// A limited object that contains only the project key.
func asProjectReference(arg any) *ProjectMarshal {
    var key string
    switch v := arg.(type) {
        case string:        key = v
        case jira.Project:  key = v.Key
        case *string:       if v != nil { key = *v }
        case *jira.Project: if v != nil { key = v.Key }
        default:            panic(fmt.Errorf("unexpected: %v", v))
    }
    if key == "" {
        return nil
    } else {
        return &ProjectMarshal{Key: &key}
    }
}

// A limited object that contains only the issue type name.
func asIssueTypeReference(arg any) *jira.IssueType {
    var name string
    switch v := arg.(type) {
        case string:            name = v
        case jira.IssueType:    name = v.Name
        case *string:           if v != nil { name = *v }
        case *jira.IssueType:   if v != nil { name = v.Name }
        default:                panic(fmt.Errorf("unexpected: %v", v))
    }
    if name == "" {
        return nil
    } else {
        return &jira.IssueType{Name: name}
    }
}

// A limited object that contains only the resolution name.
func asResolutionReference(arg any) *ResolutionMarshal {
    var name string
    switch v := arg.(type) {
        case string:            name = v
        case jira.Resolution:   name = v.Name
        case *string:           if v != nil { name = *v }
        case *jira.Resolution:  if v != nil { name = v.Name }
        default:                panic(fmt.Errorf("unexpected: %v", v))
    }
    if name == "" {
        return nil
    } else {
        return &ResolutionMarshal{Name: name}
    }
}

// A limited object that contains only the priority name.
func asPriorityReference(arg any) *jira.Priority {
    var name string
    switch v := arg.(type) {
        case string:            name = v
        case jira.Priority:     name = v.Name
        case *string:           if v != nil { name = *v }
        case *jira.Priority:    if v != nil { name = v.Name }
        default:                panic(fmt.Errorf("unexpected: %v", v))
    }
    if name == "" {
        return nil
    } else {
        return &jira.Priority{Name: name}
    }
}

// A limited object that contains only the status name.
func asStatusReference(arg any) *StatusMarshal {
    var name string
    switch v := arg.(type) {
        case string:        name = v
        case jira.Status:   name = v.Name
        case *string:       if v != nil { name = *v }
        case *jira.Status:  if v != nil { name = v.Name }
        default:            panic(fmt.Errorf("unexpected: %v", v))
    }
    if name == "" {
        return nil
    } else {
        return &StatusMarshal{Name: &name}
    }
}

// ============================================================================
// Exported variables
// ============================================================================

// Controls which jira.Issue fields are converted by asIssueMarshal().
//  @see issueDetails()
var ISSUE_MARSHAL = map[string]bool{
    "Expand":           ____,
    "ID":               true,
    "Self":             ____,
    "Key":              true,
    "Fields":           true,
    "RenderedFields":   ____,
    "Changelog":        ____,
    "Transitions":      ____,
    "Names":            ____,
}

// Controls which jira.IssueFields are converted by asIssueFieldsMarshal().
//  @see issueFieldsDetails()
var ISSUE_FIELDS_MARSHAL = map[string]bool{
    "Expand":                           ____,
    "Type":                             true,
    "Project":                          ____,
    "Environment":                      ____,
    "Resolution":                       true,
    "Priority":                         true,
    "Resolutiondate":                   true,
    "Created":                          true,
    "Duedate":                          true,
    "Watches":                          ____,
    "Assignee":                         true,
    "Updated":                          true,
    "Description":                      true,
    "Summary":                          true,
    "Creator":                          true,
    "Reporter":                         true,
    "Components":                       ____,
    "Status":                           true,
    "Progress":                         true,
    "AggregateProgress":                ____,
    "TimeTracking":                     ____,
    "TimeSpent":                        ____,
    "TimeEstimate":                     ____,
    "TimeOriginalEstimate":             ____,
    "Worklog":                          ____,
    "IssueLinks":                       ____,
    "Comments":                         ____,
    "FixVersions":                      true,
    "AffectsVersions":                  true,
    "Labels":                           true,
    "Subtasks":                         true,
    "Attachments":                      true,
    "Epic":                             true,
    "Sprint":                           ____,
    "Parent":                           ____,
    "AggregateTimeOriginalEstimate":    ____,
    "AggregateTimeSpent":               ____,
    "AggregateTimeEstimate":            ____,
    "Unknowns":                         ____,
}

// ============================================================================
// Exported types
// ============================================================================

// Used to facilitate jira.Issue JSON marshaling.
type IssueMarshal struct {
	Key            *string                      `json:"key,omitempty" structs:"key,omitempty"`
	ID             *string                      `json:"id,omitempty" structs:"id,omitempty"`
	Expand         *string                      `json:"expand,omitempty" structs:"expand,omitempty"`
	Self           *string                      `json:"self,omitempty" structs:"self,omitempty"`
	Fields         *IssueFieldsMarshal          `json:"fields,omitempty" structs:"fields,omitempty"`
	RenderedFields *jira.IssueRenderedFields    `json:"renderedFields,omitempty" structs:"renderedFields,omitempty"`
	Changelog      *jira.Changelog              `json:"changelog,omitempty" structs:"changelog,omitempty"`
	Transitions    []jira.Transition            `json:"transitions,omitempty" structs:"transitions,omitempty"`
	Names          map[string]string            `json:"names,omitempty" structs:"names,omitempty"`
}

// Used to facilitate jira.IssueFields JSON marshaling.
type IssueFieldsMarshal struct {
	Expand                        *string                   `json:"expand,omitempty" structs:"expand,omitempty"`
	Type                          *jira.IssueType           `json:"issuetype,omitempty" structs:"issuetype,omitempty"`
	Project                       *ProjectMarshal           `json:"project,omitempty" structs:"project,omitempty"`
	Environment                   *string                   `json:"environment,omitempty" structs:"environment,omitempty"`
	Resolution                    *ResolutionMarshal        `json:"resolution,omitempty" structs:"resolution,omitempty"`
	Priority                      *jira.Priority            `json:"priority,omitempty" structs:"priority,omitempty"`
	Resolutiondate                *Time                     `json:"resolutiondate,omitempty" structs:"resolutiondate,omitempty"`
	Created                       *Time                     `json:"created,omitempty" structs:"created,omitempty"`
	Duedate                       *Date                     `json:"duedate,omitempty" structs:"duedate,omitempty"`
	Watches                       *jira.Watches             `json:"watches,omitempty" structs:"watches,omitempty"`
	Assignee                      *UserMarshal              `json:"assignee,omitempty" structs:"assignee,omitempty"`
	Updated                       *Time                     `json:"updated,omitempty" structs:"updated,omitempty"`
	Description                   *string                   `json:"description,omitempty" structs:"description,omitempty"`
	Summary                       *string                   `json:"summary,omitempty" structs:"summary,omitempty"`
	Creator                       *UserMarshal              `json:"Creator,omitempty" structs:"Creator,omitempty"`
	Reporter                      *UserMarshal              `json:"reporter,omitempty" structs:"reporter,omitempty"`
	Components                    []*jira.Component         `json:"components,omitempty" structs:"components,omitempty"`
	Status                        *StatusMarshal            `json:"status,omitempty" structs:"status,omitempty"`
	Progress                      *jira.Progress            `json:"progress,omitempty" structs:"progress,omitempty"`
	AggregateProgress             *jira.Progress            `json:"aggregateprogress,omitempty" structs:"aggregateprogress,omitempty"`
	TimeTracking                  *jira.TimeTracking        `json:"timetracking,omitempty" structs:"timetracking,omitempty"`
	TimeSpent                     *int                      `json:"timespent,omitempty" structs:"timespent,omitempty"`
	TimeEstimate                  *int                      `json:"timeestimate,omitempty" structs:"timeestimate,omitempty"`
	TimeOriginalEstimate          *int                      `json:"timeoriginalestimate,omitempty" structs:"timeoriginalestimate,omitempty"`
	Worklog                       *jira.Worklog             `json:"worklog,omitempty" structs:"worklog,omitempty"`
	IssueLinks                    []*jira.IssueLink         `json:"issuelinks,omitempty" structs:"issuelinks,omitempty"`
	Comments                      *CommentsMarshal          `json:"comment,omitempty" structs:"comment,omitempty"`
	FixVersions                   []*jira.FixVersion        `json:"fixVersions,omitempty" structs:"fixVersions,omitempty"`
	AffectsVersions               []*jira.AffectsVersion    `json:"versions,omitempty" structs:"versions,omitempty"`
	Labels                        []string                  `json:"labels,omitempty" structs:"labels,omitempty"`
	Subtasks                      []*SubtasksMarshal        `json:"subtasks,omitempty" structs:"subtasks,omitempty"`
	Attachments                   []*jira.Attachment        `json:"attachment,omitempty" structs:"attachment,omitempty"`
	Epic                          *jira.Epic                `json:"epic,omitempty" structs:"epic,omitempty"`
	Sprint                        *jira.Sprint              `json:"sprint,omitempty" structs:"sprint,omitempty"`
	Parent                        *jira.Parent              `json:"parent,omitempty" structs:"parent,omitempty"`
	AggregateTimeOriginalEstimate *int                      `json:"aggregatetimeoriginalestimate,omitempty" structs:"aggregatetimeoriginalestimate,omitempty"`
	AggregateTimeSpent            *int                      `json:"aggregatetimespent,omitempty" structs:"aggregatetimespent,omitempty"`
	AggregateTimeEstimate         *int                      `json:"aggregatetimeestimate,omitempty" structs:"aggregatetimeestimate,omitempty"`
}

// Used to facilitate jira.Resolution JSON marshaling.
//  NOTE: this just fixes the field tags which don't include "omitempty"
type ResolutionMarshal struct {
    Self        string `json:"self,omitempty" structs:"self,omitempty"`
    ID          string `json:"id,omitempty" structs:"id,omitempty"`
    Description string `json:"description,omitempty" structs:"description,omitempty"`
    Name        string `json:"name,omitempty" structs:"name,omitempty"`
}

// Used to facilitate jira.Status JSON marshaling.
type StatusMarshal struct {
	Self           *string                  `json:"self,omitempty" structs:"self,omitempty"`
	Description    *string                  `json:"description,omitempty" structs:"description,omitempty"`
	IconURL        *string                  `json:"iconUrl,omitempty" structs:"iconUrl,omitempty"`
	Name           *string                  `json:"name,omitempty" structs:"name,omitempty"`
	ID             *string                  `json:"id,omitempty" structs:"id,omitempty"`
	StatusCategory *StatusCategoryMarshal   `json:"statusCategory,omitempty" structs:"statusCategory,omitempty"`
}

// Used to facilitate jira.Subtasks JSON marshaling.
//  NOTE: this just fixes the field tags which don't include "omitempty"
type SubtasksMarshal struct {
	ID     *string              `json:"id,omitempty" structs:"id,omitempty"`
	Key    *string              `json:"key,omitempty" structs:"key,omitempty"`
	Self   *string              `json:"self,omitempty" structs:"self,omitempty"`
	Fields *IssueFieldsMarshal  `json:"fields,omitempty" structs:"fields,omitempty"`
}

// Used to facilitate jira.StatusCategory JSON marshaling.
//  NOTE: this just fixes the field tags which don't include "omitempty"
type StatusCategoryMarshal struct {
	Self      string `json:"self,omitempty" structs:"self,omitempty"`
	ID        int    `json:"id,omitempty" structs:"id,omitempty"`
	Name      string `json:"name,omitempty" structs:"name,omitempty"`
	Key       string `json:"key,omitempty" structs:"key,omitempty"`
	ColorName string `json:"colorName,omitempty" structs:"colorName,omitempty"`
}

// Comments represents a list of Comment.
type CommentsMarshal struct {
	Comments []*CommentMarshal `json:"comments,omitempty" structs:"comments,omitempty"`
}
