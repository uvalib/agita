// Jira/project_util.go
//
// Operations on jira.Project objects.

package Jira

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported types
// ============================================================================

// A generic project reference.
type ProjectArg interface {
    Project | jira.Project | *Project | *jira.Project
}

// ============================================================================
// Exported functions
// ============================================================================

// Transform a generic project reference.
func NormalizeProject[T ProjectArg](arg T) *jira.Project {
    switch v := any(arg).(type) {
        case Project:       return v.ptr
        case jira.Project:  return &v
        case *Project:      if v != nil { return v.ptr }
        case *jira.Project: return v
        default:            panic(fmt.Errorf("unexpected: %v", v))
    }
    return nil
}

// Return the project key.
func ProjectKey[T ProjectArg](arg T) string {
    if proj := NormalizeProject(arg); proj == nil {
        return ""
    } else {
        return proj.Key
    }
}

// Indicate whether the two project objects refer to the same project.
func SameProject[T1, T2 ProjectArg](arg1 T1, arg2 T2) bool {
    if nil1, nil2 := util.IsNil(arg1), util.IsNil(arg2); nil1 || nil2 {
        return nil1 && nil2
    } else {
        return ProjectKey(arg1) == ProjectKey(arg2)
    }
}
