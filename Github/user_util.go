// Github/user_util.go

package Github

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported types
// ============================================================================

// A generic user account reference.
type UserArg interface {
    User | github.User | *User | *github.User
}

// ============================================================================
// Exported functions
// ============================================================================

// Transform a generic user account reference.
func NormalizeUser[T UserArg](arg T) *github.User {
    switch v := any(arg).(type) {
        case User:         return v.ptr
        case github.User:  return &v
        case *User:        if v != nil { return v.ptr }
        case *github.User: return v
        default:           panic(fmt.Errorf("unexpected: %v", v))
    }
    return nil
}

// Transform users into login names.
func Accounts[T UserArg](users []T) []string {
    result := make([]string, 0, len(users))
    for _, user := range users {
        if login := Account(user); login != "" {
            result = append(result, login)
        }
    }
    return result
}

// Return the login name for the given user.
func Account[T UserArg](arg T) string {
    switch v := any(arg).(type) {
        case User:  if v.Login != "" { return v.Login }
        case *User: if (v != nil) && (v.Login != "") { return v.Login }
    }
    if user := NormalizeUser(arg); (user == nil) || (user.Login == nil) {
        return ""
    } else {
        return *user.Login
    }
}

// Indicate whether the two user objects refer to the same user.
func SameAccount[T1, T2 UserArg](arg1 T1, arg2 T2) bool {
    if nil1, nil2 := util.IsNil(arg1), util.IsNil(arg2); nil1 || nil2 {
        return nil1 && nil2
    } else {
        return Account(arg1) == Account(arg2)
    }
}
