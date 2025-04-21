// Jira/global_test.go
//
// Test values used throughout the package.

package Jira

import (
	"reflect"
	"testing"

	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Tests
// ============================================================================

func TestCommentFromJson(t *testing.T) {
    const fn = "CommentFromJson"

    src := GetCommentById(TestClient, SAMPLE_ISSUE, SAMPLE_COMMENTS[0])
    buf, err := src.MarshalJSON()
    log.ErrorValue(err)
    dst := CommentFromJson(string(buf))

    mapSrc, mapDst := util.StructMap(src.ptr), util.StructMap(dst.ptr)
    for field, inDst := range COMMENT_MARSHAL {
        if got, want := mapDst[field], mapSrc[field]; inDst && testMarshaled(want) {
            testVerifyMarshal(fn, field, got, want, t)
        }
    }
}

func TestIssueFromJson(t *testing.T) {
    const fn = "IssueFromJson"

    src := GetIssueByKey(TestClient, SAMPLE_ISSUE)
    buf, err := src.MarshalJSON()
    log.ErrorValue(err)
    dst := IssueFromJson(string(buf))

    mapSrc, mapDst := util.StructMap(src.ptr), util.StructMap(dst.ptr)
    for field, inDst := range ISSUE_MARSHAL {
        if got, want := mapDst[field], mapSrc[field]; inDst && testMarshaled(want) {
            if field == "Fields" {
                mSrc, mDst := util.StructMap(want), util.StructMap(got)
                for fld, used := range ISSUE_FIELDS_MARSHAL {
                    if g, w := mDst[fld], mSrc[fld]; used && testMarshaled(w) {
                        testVerifyMarshal(fn, fld, g, w, t)
                    }
                }
            } else {
                testVerifyMarshal(fn, field, got, want, t)
            }
        }
    }
}

func TestProjectFromJson(t *testing.T) {
    const fn = "ProjectFromJson"

    src := GetProjectByKey(TestClient, SAMPLE_PROJ)
    buf, err := src.MarshalJSON()
    log.ErrorValue(err)
    dst := ProjectFromJson(string(buf))

    mapSrc, mapDst := util.StructMap(src.ptr), util.StructMap(dst.ptr)
    for field, inDst := range PROJECT_MARSHAL {
        if got, want := mapDst[field], mapSrc[field]; inDst && testMarshaled(want) {
            testVerifyMarshal(fn, field, got, want, t)
        }
    }
}

// ============================================================================
// Internal functions
// ============================================================================

// Indicate whether the value would be one that would actually get marshaled
// when converting a Jira object to JSON.  If not then the `got` object would
// not have a field value that the source `want` object happens to have.
func testMarshaled(arg any) bool {
    switch any(arg).(type) {
        case int, *int:                                         return nil != asIntMarshal(arg)
        case Date, *Date:                                       return nil != asDateMarshal(arg)
        case Time, *Time:                                       return nil != asTimeMarshal(arg)
        case string, *string:                                   return nil != asStringMarshal(arg)
        case jira.User, *jira.User:                             return nil != asUserReference(arg)
        case jira.Status, *jira.Status:                         return nil != asStatusReference(arg)
        case jira.Project, *jira.Project:                       return nil != asProjectReference(arg)
        case jira.Priority, *jira.Priority:                     return nil != asPriorityReference(arg)
        case jira.Progress, *jira.Progress:                     return nil != asProgressMarshal(arg)
        case jira.IssueType, *jira.IssueType:                   return nil != asIssueTypeReference(arg)
        case jira.Resolution, *jira.Resolution:                 return nil != asResolutionReference(arg)
        case jira.ProjectCategory, *jira.ProjectCategory:       return nil != asProjectCategoryReference(arg)
        case jira.ProjectComponent, *jira.ProjectComponent:     return nil != asProjectComponentReference(arg)
        case jira.CommentVisibility, *jira.CommentVisibility:   return nil != asCommentVisibilityMarshal(arg)
    }
    return !util.IsEmpty(arg)
}

// Display an error unless `got` is equivalent to `want`.
func testVerifyMarshal(fn, field string, got any, want any, t *testing.T) {
    if !testCompare(got, want) {
        t.Errorf("%s(): %s got != want", fn, field)
        t.Errorf("%s(): %s got  = %v", fn, field, got)
        t.Errorf("%s(): %s want = %v", fn, field, want)
    }
}

// Indicate whether two items are equivalent by determining whether their
// normalized values are equal.
//  NOTE: slices are assumed to have elements in the same order
func testCompare(got any, want any) bool {
    // If one is empty then they're different unless both are empty.
    if emt1, emt2 := util.IsEmpty(got), util.IsEmpty(want); emt1 || emt2 {
        return emt1 && emt2
    }

    // If one is an array or slice then they're different if any of their
    // elements are different.
    if ary1, ary2 := util.IsArray(got), util.IsArray(want); ary1 || ary2 {
        if !ary1 || !ary2 {
            return false
        } else if s1, s2 := testSlice(got), testSlice(want); len(s1) != len(s2) {
            return false
        } else {
            for i := range s1 {
                if !testCompare(s1[i], s2[i]) {
                    return false
                }
            }
        }
        return true
    }

    // Otherwise compare their values as they would be marshaled.
    switch want.(type) {
        case int, *int:                                         got, want = asIntMarshal(got),                  asIntMarshal(want)
        case Date, *Date:                                       got, want = asDateMarshal(got),                 asDateMarshal(want)
        case Time, *Time:                                       got, want = asTimeMarshal(got),                 asTimeMarshal(want)
        case string, *string:                                   got, want = asStringMarshal(got),               asStringMarshal(want)
        case jira.User, *jira.User:                             got, want = Account(got),                       Account(want)
        case jira.Status, *jira.Status:                         got, want = asStatusReference(got),             asStatusReference(want)
        case jira.Project, *jira.Project:                       got, want = asProjectReference(got),            asProjectReference(want)
        case jira.Progress, *jira.Progress:                     got, want = asProgressMarshal(got),             asProgressMarshal(want)
        case jira.Priority, *jira.Priority:                     got, want = asPriorityReference(got),           asPriorityReference(want)
        case jira.IssueType, *jira.IssueType:                   got, want = asIssueTypeReference(got),          asIssueTypeReference(want)
        case jira.Resolution, *jira.Resolution:                 got, want = asResolutionReference(got),         asResolutionReference(want)
        case jira.ProjectCategory, *jira.ProjectCategory:       got, want = asProjectCategoryReference(got),    asProjectCategoryReference(want)
        case jira.ProjectComponent, *jira.ProjectComponent:     got, want = asProjectComponentReference(got),   asProjectComponentReference(want)
        case jira.CommentVisibility, *jira.CommentVisibility:   got, want = asCommentVisibilityMarshal(got),    asCommentVisibilityMarshal(want)
    }
    return reflect.DeepEqual(got, want)
}

// Create a slice from an arbitrary source known to be an array or slice.
func testSlice(src any) []any {
    value  := reflect.ValueOf(src)
    size   := value.Len()
    result := make([]any, 0, size)
    for i := range size {
        if v := value.Index(i); v.CanInterface() {
            result = append(result, v.Interface())
        }
    }
    return result
}