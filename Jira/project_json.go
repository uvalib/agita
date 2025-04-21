// Jira/project_json.go
//
// JSON marshaling of jira.Project objects.

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

// Translate a Jira project into JSON.
func (p *Project) MarshalJSON() ([]byte, error) {
    if m := p.AsProjectMarshal(); m == nil {
        return []byte{}, nil
    } else {
        return json.Marshal(m)
    }
}

// Create a Project object suitable for JSON marshaling.
func (p *Project) AsProjectMarshal() *ProjectMarshal {
    if noProject(p) {
        return nil
    } else {
        return asProjectMarshal(p.ptr)
    }
}

// ============================================================================
// Exported functions
// ============================================================================

// Reconstitute a Project object from JSON.
func ProjectFromJson(src string) *Project {
    dst := &jira.Project{}
    if err := json.Unmarshal([]byte(src), dst); log.ErrorValue(err) == nil {
        return &Project{ptr: dst}
    } else {
        return nil
    }
}

// ============================================================================
// Internal functions
// ============================================================================

// Create a Project object suitable for JSON marshaling.
func asProjectMarshal(p *jira.Project) *ProjectMarshal {
    if (p == nil) || (p.Key == "") { return nil }
    use := func(field string) bool { return PROJECT_MARSHAL[field] }
    res := ProjectMarshal{}
    if use("Key")               { res.Key               = asStringMarshal(p.Key) }
    if use("ID")                { res.ID                = asStringMarshal(p.ID) }
    if use("Expand")            { res.Expand            = asStringMarshal(p.Expand) }
    if use("Self")              { res.Self              = asStringMarshal(p.Self) }
    if use("Description")       { res.Description       = asStringMarshal(p.Description) }
    if use("Lead")              { res.Lead              = asUserReference(p.Lead) }
    if use("Components")        { res.Components        = asProjectComponentReferences(p.Components) }
    if use("IssueTypes")        { res.IssueTypes        = asIssueTypeReferences(p.IssueTypes) }
    if use("URL")               { res.URL               = asStringMarshal(p.URL) }
    if use("Email")             { res.Email             = asStringMarshal(p.Email) }
    if use("AssigneeType")      { res.AssigneeType      = asStringMarshal(p.AssigneeType) }
    if use("Versions")          { res.Versions          = p.Versions }
    if use("Name")              { res.Name              = asStringMarshal(p.Name) }
    if use("Roles")             { res.Roles             = p.Roles }
    if use("AvatarUrls")        { res.AvatarUrls        = &p.AvatarUrls }
    if use("ProjectCategory")   { res.ProjectCategory   = asProjectCategoryReference(p.ProjectCategory) }
    return &res
}

// A limited object that contains only the project category name.
func asProjectCategoryReference(arg any) *ProjectCategoryMarshal {
    var name string
    switch v := arg.(type) {
        case string:                name = v
        case jira.ProjectCategory:  name = v.Name
        case *string:               if v != nil { name = *v }
        case *jira.ProjectCategory: if v != nil { name = v.Name }
        default:                    panic(fmt.Errorf("unexpected: %v", v))
    }
    if name == "" {
        return nil
    } else {
        return &ProjectCategoryMarshal{Name: name}
    }
}

// A limited object that contains only the project component name.
func asProjectComponentReference(arg any) *ProjectComponentMarshal {
    var name string
    switch v := arg.(type) {
        case string:                    name = v
        case jira.ProjectComponent:     name = v.Name
        case *string:                   if v != nil { name = *v }
        case *jira.ProjectComponent:    if v != nil { name = v.Name }
        default:                        panic(fmt.Errorf("unexpected: %v", v))
    }
    if name == "" {
        return nil
    } else {
        return &ProjectComponentMarshal{Name: &name}
    }
}

// A slice of limited objects containing only the project component names.
func asProjectComponentReferences(src []jira.ProjectComponent) []ProjectComponentMarshal {
    dst := make([]ProjectComponentMarshal, 0, len(src))
    for _, elem := range src {
        dst = append(dst, *asProjectComponentReference(elem))
    }
    return dst
}

// A slice of limited objects containing only the issue type names.
func asIssueTypeReferences(src []jira.IssueType) []jira.IssueType {
    dst := make([]jira.IssueType, 0, len(src))
    for _, elem := range src {
        dst = append(dst, *asIssueTypeReference(elem))
    }
    return dst
}

// ============================================================================
// Exported variables
// ============================================================================

// Controls which jira.Project fields are converted by asProjectMarshal().
//  @see projectDetails()
var PROJECT_MARSHAL = map[string]bool{
    "Expand":           ____,
    "Self":             ____,
    "ID":               true,
    "Key":              true,
    "Description":      true,
    "Lead":             true,
    "Components":       true,
    "IssueTypes":       true,
    "URL":              ____,
    "Email":            ____,
    "AssigneeType":     ____,
    "Versions":         ____,
    "Name":             true,
    "Roles":            ____,
    "AvatarUrls":       ____,
    "ProjectCategory":  true,
}

// ============================================================================
// Exported types
// ============================================================================

// Used to facilitate jira.Project JSON marshaling.
type ProjectMarshal struct {
	Key             *string                     `json:"key,omitempty" structs:"key,omitempty"`
	ID              *string                     `json:"id,omitempty" structs:"id,omitempty"`
	Expand          *string                     `json:"expand,omitempty" structs:"expand,omitempty"`
	Self            *string                     `json:"self,omitempty" structs:"self,omitempty"`
	Description     *string                     `json:"description,omitempty" structs:"description,omitempty"`
	Lead            *UserMarshal                `json:"lead,omitempty" structs:"lead,omitempty"`
	Components      []ProjectComponentMarshal   `json:"components,omitempty" structs:"components,omitempty"`
	IssueTypes      []jira.IssueType            `json:"issueTypes,omitempty" structs:"issueTypes,omitempty"`
	URL             *string                     `json:"url,omitempty" structs:"url,omitempty"`
	Email           *string                     `json:"email,omitempty" structs:"email,omitempty"`
	AssigneeType    *string                     `json:"assigneeType,omitempty" structs:"assigneeType,omitempty"`
	Versions        []jira.Version              `json:"versions,omitempty" structs:"versions,omitempty"`
	Name            *string                     `json:"name,omitempty" structs:"name,omitempty"`
	Roles           map[string]string           `json:"roles,omitempty" structs:"roles,omitempty"`
	AvatarUrls      *jira.AvatarUrls            `json:"avatarUrls,omitempty" structs:"avatarUrls,omitempty"`
	ProjectCategory *ProjectCategoryMarshal     `json:"projectCategory,omitempty" structs:"projectCategory,omitempty"`
}

// Used to facilitate jira.ProjectCategory JSON marshaling.
//  NOTE: this just fixes the field tags which don't include "omitempty"
type ProjectCategoryMarshal struct {
	Self        string `json:"self,omitempty" structs:"self,omitempty"`
	ID          string `json:"id,omitempty" structs:"id,omitempty"`
	Name        string `json:"name,omitempty" structs:"name,omitempty"`
	Description string `json:"description,omitempty" structs:"description,omitempty"`
}

// Used to facilitate jira.ProjectComponent JSON marshaling.
type ProjectComponentMarshal struct {
	Self                *string         `json:"self,omitempty" structs:"self,omitempty"`
	ID                  *string         `json:"id,omitempty" structs:"id,omitempty"`
	Name                *string         `json:"name,omitempty" structs:"name,omitempty"`
	Description         *string         `json:"description,omitempty" structs:"description,omitempty"`
	Lead                *UserMarshal    `json:"lead,omitempty" structs:"lead,omitempty"`
	AssigneeType        *string         `json:"assigneeType,omitempty" structs:"assigneeType,omitempty"`
	Assignee            *UserMarshal    `json:"assignee,omitempty" structs:"assignee,omitempty"`
	RealAssigneeType    *string         `json:"realAssigneeType,omitempty" structs:"realAssigneeType,omitempty"`
	RealAssignee        *UserMarshal    `json:"realAssignee,omitempty" structs:"realAssignee,omitempty"`
	IsAssigneeTypeValid *bool           `json:"isAssigneeTypeValid,omitempty" structs:"isAssigneeTypeValid,omitempty"`
	Project             *string         `json:"project,omitempty" structs:"project,omitempty"`
	ProjectID           *int            `json:"projectId,omitempty" structs:"projectId,omitempty"`
}
