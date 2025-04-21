// Github/repository_util.go

package Github

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// A generic repository reference.
type RepositoryArg interface {
    Repository | github.Repository | *Repository | *github.Repository
}

// ============================================================================
// Exported functions
// ============================================================================

// Transform a generic repository reference.
func NormalizeRepository[T RepositoryArg](arg T) *github.Repository {
    switch v := any(arg).(type) {
        case Repository:         return v.ptr
        case github.Repository:  return &v
        case *Repository:        if v != nil { return v.ptr }
        case *github.Repository: return v
        default:                 panic(fmt.Errorf("unexpected: %v", v))
    }
    return nil
}

// Return the repository ID.
func RepoNumber[T RepositoryArg](arg T) int64 {
    if repo := NormalizeRepository(arg); (repo == nil) || (repo.ID == nil) {
        return 0
    } else {
        return *repo.ID
    }
}

// Return the login name for the repository owner.
func RepoOwner[T RepositoryArg](arg T) string {
    switch v := any(arg).(type) {
        case Repository:  if v.Owner != "" { return v.Owner }
        case *Repository: if (v != nil) && (v.Owner != "") { return v.Owner }
    }
    if repo := NormalizeRepository(arg); (repo == nil) || (repo.Owner == nil) {
        return ""
    } else {
        return Account(repo.Owner)
    }
}

// Return the repository name.
func RepoName[T RepositoryArg](arg T) string {
    switch v := any(arg).(type) {
        case Repository:  if v.Name != "" { return v.Name }
        case *Repository: if (v != nil) && (v.Name != "") { return v.Name }
    }
    if repo := NormalizeRepository(arg); (repo == nil) || (repo.Name == nil) {
        return ""
    } else {
        return *repo.Name
    }
}

// Return `owner/repo` for the given repository.
//  NOTE: a nil pointer returns "-/-"
func RepoPath[T RepositoryArg](arg T) string {
    owner, name := RepoOwner(arg), RepoName(arg)
    if owner == "" { owner = MISSING }
    if name  == "" { name  = MISSING }
    return owner + "/" + name
}

// Indicate whether the two repository objects refer to the same repository.
func SameRepository[T1, T2 RepositoryArg](arg1 T1, arg2 T2) bool {
    if nil1, nil2 := util.IsNil(arg1), util.IsNil(arg2); nil1 || nil2 {
        return nil1 && nil2
    } else if num1, num2 := RepoNumber(arg1), RepoNumber(arg2); (num1 != 0) || (num2 != 0) {
        return num1 == num2
    } else {
        return RepoPath(arg1) == RepoPath(arg2)
    }
}
