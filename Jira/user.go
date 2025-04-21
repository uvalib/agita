// jira/user.go

package Jira

import (
	"fmt"

	"lib.virginia.edu/agita/util"

	"github.com/andygrunwald/go-jira"
)

// ============================================================================
// Exported variables
// ============================================================================

// A mapping of Jira account to user full name.
var JiraUser = map[string]string {
    "adj5j":    "Tony Jones",
    "ag5vy":    "Anne Gaynor",
    "akl3b":    "Amber Reichert",
    "arb5w":    "Ann Burns",
    "au3vk":    "Ashish Upadhyaya",
    "bcb4y":    "Brandon Butler",
    "bds2mv":   "Ben Spector",
    "bg9ba":    "Brenda Gunn",
    "bga3d":    "Bethany Anderson",
    "bmx7wf":   "Marco Battistella",
    "bwb9f":    "Beth Blanton",
    "bwo3db":   "Ben Ormond",
    "cma4u":    "Carla Arton",
    "cmm2t":    "Christina Deane",
    "cmw5mc":   "Christopher Welte",
    "csk5vf":   "Connor Kenaston",
    "dhc4z":    "Doug Chestnut",
    "dpg3k":    "Dave Goldstein",
    "ecr2c":    "Ellen Ramsey",
    "ege5vd":   "Molly Fair",
    "eks4x":    "Eric Seidel",
    "emg3b":    "Eliza Gilligan",
    "er8fn":    "Liz Rapp",
    "gmk3d":    "Ginny Kois",
    "hdn9c":    "Hanni Nabahe",
    "hmh8eq":   "Akari Hernandez",
    "ja3nu":    "Joseph Azizi",
    "jgc4q":    "Zeke Crater",
    "jh9xp":    "Jill Heinze",
    "jkb2b":    "Jeremy Boggs",
    "jlj5aj":   "Jason Jordan",
    "jlk4p":    "Jack Kelly",
    "jmu2m":    "John Unsworth",
    "jor2a":    "Jennifer Roper",
    "jph9e":    "Jeff Hill",
    "jtb4t":    "Jeremy Bartczak",
    "ka7uz":    "Krystal Appiah",
    "khj5c":    "Kristen Jensen",
    "kmm6ef":   "Kara McClurken",
    "lar4k":    "Leigh Rockey",
    "ldc8n":    "Linda Vaughan",
    "lf6f":     "Lou Foster",
    "libx-sat": "SAT Team",
    "lsc6v":    "Lorrie Chisholm",
    "lw2cd":    "Lauren Work",
    "lws4n":    "Lucie Stylianopoulos",
    "lzl9b":    "Lauren Longwell",
    "md5wz":    "Mike Durbin",
    "mdw7p":    "Melanie Williams",
    "mhm8m":    "Heather Riser",
    "mhw8m":    "Mark Witteman",
    "mrc2x":    "Marc Campbell",
    "mwm7b":    "Molly Minturn",
    "naw4t":    "Nestor Walker",
    "nir4x":    "Nicole Royal",
    "np6hd":    "Nitesh Parajuli",
    "nqt2bq":   "Ryan Russo",
    "pw7e":     "Peter Welch",
    "rac8f":    "Rebecca Coleman",
    "rah6w":    "Ric Hodges",
    "rar6u":    "Renee Reighart",
    "rcm7e":    "Rennie Mapp",
    "rds4w":    "Rob Smith",
    "rfg":      "Rebecca Garver",
    "rh9ec":    "Bob Haschart",
    "rmg6f":    "Bob Gartland",
    "rpk2kn":   "Ryan Kelly",
    "rsl6m":    "Robin Ruggaber",
    "rwl":      "Ray Lubinsky",
    "sah":      "Sherry Lake",
    "sd3gz":    "Sue Donovan",
    "sh2rd":    "Stephanie Hunter",
    "smm3m":    "Sean McCord",
    "sp7fg":    "Sony Prosper",
    "spr7b":    "Sue Richeson",
    "srv3f":    "Steven Villereal",
    "sss3s":    "Susan Neal",
    "stg2s":    "Stan Gunn",
    "tnn7yc":   "Tho Nguyen",
    "trb5f":    "Tammy Barbour",
    "tss6n":    "Tim Stevens",
    "tsh2k":    "Trillian Hosticka",
    "tvf3c":    "Tracy Fewell",
    "vdg8v":    "Dave Griles",
    "vlk4n":    "Vince Kois",
    "wdw5ch":   "Will Wyatt",
    "wkb5j":    "Winston Barham",
    "wmr5a":    "Will Rourk",
    "wxf2nb":   "Ruobing Su",
    "xjk5yt":   "Nora Wilkerson",
    "xw5d":     "Xiaoming Wang",
    "ys2ck":    "Yukesh Sitoula",
    "ys2n":     "Yuji Shinozaki",
}

// ============================================================================
// Exported types
// ============================================================================

// A generic user account reference.
type UserArg interface {
    jira.User | *jira.User
}

// ============================================================================
// Exported functions
// ============================================================================

// Return the Jira account name in a form that includes the user's full name
// if included in the JiraUser map.
func AppendFullName(jiraAccount string) string {
    if fullName := JiraUser[jiraAccount]; fullName == "" {
        return jiraAccount
    } else {
        return fmt.Sprintf("%s (%s)", jiraAccount, fullName)
    }
}

// Return a non-blank string representing the Jira account for the argument.
func UserLabel[T UserArg](arg T) string {
    name := Account(arg)
    switch {
        case name != "":      return name
        case util.IsNil(arg): return "[nil]"
        default:              return "[blank]"
    }
}

// ============================================================================
// Exported functions
// ============================================================================

// Return the login name for the given user.
func Account(arg any) string {
    switch v := arg.(type) {
        case string:     return v
        case *string:    if v != nil { return *v }
        case jira.User:  return v.Name
        case *jira.User: if v != nil { return v.Name }
        default:         panic(fmt.Errorf("unexpected: %v", v))
    }
    return ""
}

// Indicate whether the two User objects refer to the same user.
func SameAccount[T1, T2 any](arg1 T1, arg2 T2) bool {
    if nil1, nil2 := util.IsNil(arg1), util.IsNil(arg2); nil1 || nil2 {
        return nil1 && nil2
    } else {
        return Account(arg1) == Account(arg2)
    }
}
