// jira/user_json.go
//
// JSON marshal Jira user objects.

package Jira

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Internal functions
// ============================================================================

// A limited object that contains only the user account name.
func asUserReference(arg any) *UserMarshal {
    var name string
    switch v := arg.(type) {
        case string:        name = v
        case jira.User:     name = v.Name
        case *string:       if v != nil { name = *v }
        case *jira.User:    if v != nil { name = v.Name }
        default:            panic(fmt.Errorf("unexpected: %v", v))
    }
    if name == "" {
        return nil
    } else {
        return &UserMarshal{Name: &name}
    }
}

// ============================================================================
// Exported types
// ============================================================================

// Used to facilitate jira.User JSON marshaling.
type UserMarshal struct {
    Self            *string             `json:"self,omitempty" structs:"self,omitempty"`
    AccountID       *string             `json:"accountId,omitempty" structs:"accountId,omitempty"`
    AccountType     *string             `json:"accountType,omitempty" structs:"accountType,omitempty"`
    Name            *string             `json:"name,omitempty" structs:"name,omitempty"`
    Key             *string             `json:"key,omitempty" structs:"key,omitempty"`
  //Password        *string             `json:"-"` // [1]
    EmailAddress    *string             `json:"emailAddress,omitempty" structs:"emailAddress,omitempty"`
    AvatarUrls      *jira.AvatarUrls    `json:"avatarUrls,omitempty" structs:"avatarUrls,omitempty"` // [2]
    DisplayName     *string             `json:"displayName,omitempty" structs:"displayName,omitempty"`
    Active          *bool               `json:"active,omitempty" structs:"active,omitempty"`
    TimeZone        *string             `json:"timeZone,omitempty" structs:"timeZone,omitempty"`
    Locale          *string             `json:"locale,omitempty" structs:"locale,omitempty"`
    ApplicationKeys *[]string           `json:"applicationKeys,omitempty" structs:"applicationKeys,omitempty"`
}
