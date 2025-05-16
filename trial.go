// trial.go
//
// Trials only work on information extracted from Jira and/or GitHub.
// They never affect the actual Jira projects or GitHub repositories.

package main

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/util"

	"lib.virginia.edu/agita/Github"
	"lib.virginia.edu/agita/Jira"
)

// ============================================================================
// Constants used to set flags if no command line arguments are given.
// ============================================================================

const JiraTrials        = false
const JiraProjects      = true
const JiraIssues        = true
const JiraComments      = true

const GithubTrials      = false
const GithubRate        = true
const GithubUsers       = true
const GithubRepos       = true
const GithubIssues      = true
const GithubComments    = true

const GraphQlTrials     = true

const TransferTrials    = false

// If this -trial argument is given, all flags are set regardless of the above
// constant values.
const ALL_TRIALS = "all"

// ============================================================================
// Types
// ============================================================================

type TrialFlags struct {
    jira            bool
    jiraProjects    bool
    jiraIssues      bool
    jiraComments    bool
    github          bool
    githubRate      bool
    githubUsers     bool
    githubRepos     bool
    githubIssues    bool
    githubComments  bool
    graphQl         bool
    transfer        bool
}

// ============================================================================
// Variables
// ============================================================================

var trial TrialFlags

// ============================================================================
// Variables
// ============================================================================

var TRIAL_PROJECTS = []string{"RWL"}

var TRIAL_REPOS = []string{"emma-ia"}

// ============================================================================
// Functions
// ============================================================================

// Run the trials specified on the command line.
func TrialAll(args ...string) {
    names := trialArgs(args...)
    if trial.jira     { TrialJira();             fmt.Println("") }
    if trial.github   { TrialGithub();           fmt.Println("") }
    if trial.graphQl  { TrialGraphQL(names...);  fmt.Println("") }
    if trial.transfer { TrialTransfer(names...); fmt.Println("") }
    fmt.Println("TRIALS COMPLETE")
}

// Sets trial flags based on the provided arguments.
// Any arguments which do not match a flag name are returned.
// If `args` is empty then trial flags are initialized based on the constants.
// Otherwise, only the trial flag(s) specified by the arguments will be set.
func trialArgs(args ...string) []string {
    names := []string{}
    if len(args) == 0 {
        trial.jira           = JiraTrials
        trial.jiraProjects   = JiraProjects
        trial.jiraIssues     = JiraIssues
        trial.jiraComments   = JiraComments
        trial.github         = GithubTrials
        trial.githubRate     = GithubRate
        trial.githubUsers    = GithubUsers
        trial.githubRepos    = GithubRepos
        trial.githubIssues   = GithubIssues
        trial.githubComments = GithubComments
        trial.graphQl        = GraphQlTrials
        trial.transfer       = TransferTrials

    } else if args[0] == ALL_TRIALS {

        trial.jira           = true
        trial.jiraProjects   = JiraProjects
        trial.jiraIssues     = JiraIssues
        trial.jiraComments   = JiraComments
        trial.github         = true
        trial.githubRate     = GithubRate
        trial.githubUsers    = GithubUsers
        trial.githubRepos    = GithubRepos
        trial.githubIssues   = GithubIssues
        trial.githubComments = GithubComments
        trial.graphQl        = true
        trial.transfer       = true

    } else {
        for _, arg := range trialExpandedArgs(args...) {
            switch {
                case trialProcessArg(arg):          // try directly
                case trialProcessArg(arg + "s"):    // try pluralized
                default:                            names = append(names, arg)
            }
        }

        if trial.jiraProjects || trial.jiraIssues  || trial.jiraComments {
            trial.jira = true
        } else if trial.jira {
            trial.jiraProjects = JiraProjects
            trial.jiraIssues   = JiraIssues
            trial.jiraComments = JiraComments
        }

        if trial.githubRate || trial.githubUsers || trial.githubRepos || trial.githubIssues || trial.githubComments {
            trial.github = true
        } else if trial.github {
            trial.githubRate     = GithubRate
            trial.githubUsers    = GithubUsers
            trial.githubRepos    = GithubRepos
            trial.githubIssues   = GithubIssues
            trial.githubComments = GithubComments
        }
    }
    return names
}

// Set the trial flag indicated by `arg`.
//  NOTE: `arg` is assumed to be lowercase.
func trialProcessArg(arg string) bool {
    switch strings.ToLower(arg) {
        case "jira":            trial.jira              = true
        case "jiraprojects":    trial.jiraProjects      = true
        case "jiraissues":      trial.jiraIssues        = true
        case "jiracomments":    trial.jiraComments      = true
        case "github":          trial.github            = true
        case "githubrate":      trial.githubRate        = true
        case "githubusers":     trial.githubUsers       = true
        case "githubrepos":     trial.githubRepos       = true
        case "githubissues":    trial.githubIssues      = true
        case "githubcomments":  trial.githubComments    = true
        case "graphql":         trial.graphQl           = true
        case "transfer":        trial.transfer          = true
        default:                return false
    }
    return true
}

// Analyze arguments into trial flags.
func trialExpandedArgs(args ...string) []string {
    expanded := []string{}
    for _, arg := range args {
        expanded = append(expanded, trialFlagArgs(arg)...)
    }
    return expanded
}

// Transform an argument like
//  "Jira:Issues,Comments"
// to return an array of trial flags like
//  []string{"jira", "jiraissues", "jiracomments"}
func trialFlagArgs(arg string) []string {
    part   := strings.Split(arg, ":")
    parts  := len(part)
    result := []string{}
    if parts > 0 {
        if parts > 2 {
            panic("malformed")
        } else if parts > 1 {
            names := strings.Split(strings.ToLower(part[1]), ",")
            switch mode := strings.ToLower(part[0]); mode {
                case "jira", "github":
                    for i, name := range names {
                        names[i] = mode + name
                    }
            }
            result = append(result, names...)
        } else {
            result = append(result, arg)
        }
    }
    return result
}

// ============================================================================
// Jira testing
// ============================================================================

// Exercise the Jira API.
func TrialJira() bool {
    trialTitle("Jira")
    fmt.Println()
    if !Jira.Initialize() {
        fmt.Println("*** FAILED TO INITIALIZE JIRA INTERFACE ***")
        return false
    }

    JIRA_PROJ    := Jira.SAMPLE_PROJ
    JIRA_ISSUE   := Jira.IssueKey("EMMA-7")
    JIRA_COMMENT := Jira.CommentId(1148857)

    fmt.Printf("*** JIRA_PROJ    = %v\n", JIRA_PROJ)
    fmt.Printf("*** JIRA_ISSUE   = %v\n", JIRA_ISSUE)
    fmt.Printf("*** JIRA_COMMENT = %v\n", JIRA_COMMENT)
    fmt.Printf("*** ProjectByKey = %v\n", Jira.ProjectByKey)

    cli := Jira.NewClient()
    prj := cli.GetProjectByKey(JIRA_PROJ)

    if trial.jiraProjects {
        prj.Print()
        cli.PrintProjects()
    }

    if trial.jiraIssues {
        all := prj.Issues()
        fmt.Printf("\n*** FETCH EMMA JIRA ISSUES (%d):\n", len(all))
        for _, issue := range all {
            issue.Print()
        }

        min, max := "CSH-1300", "CSH-1399"
        fmt.Printf("\n*** FETCH CSH JIRA ISSUES %q to %q:\n", min, max)
        for i, issue := range cli.GetProjectByKey("CSH").GetIssues(min, max) {
            fmt.Printf("*** %d: %s\n", i, issue.Key())
        }
    }

    if trial.jiraComments {
        key := JIRA_ISSUE
        iss := prj.GetIssue(key)
        all := iss.Comments()
        fmt.Printf("\n*** FETCH COMMENTS FOR JIRA ISSUE %q (%d):\n", key, len(all))
        for _, comment := range all {
            comment.Print()
        }
        com := iss.Comment(JIRA_COMMENT)
        com.Print()
    }

    return true
}

// ============================================================================
// GitHub testing
// ============================================================================

// Exercise the GitHub API.
func TrialGithub() bool {
    trialTitle("GitHub")
    fmt.Println()
    if !Github.Initialize() {
        fmt.Println("*** FAILED TO INITIALIZE GITHUB INTERFACE ***")
        return false
    }

    var GITHUB_ORG      = Github.SAMPLE_ORG
    var GITHUB_REPO     = Github.SAMPLE_REPO
    var GITHUB_ISSUE    = int(Github.SAMPLE_ISSUE)
    var GITHUB_COMMENT  = int64(2649639293)
    var GITHUB_USERS    = []string{"RayLubinsky", "nestorw", "UVADave"}
    var GITHUB_USER     = GITHUB_USERS[0]

    fmt.Printf("*** GITHUB_ORG     = %v\n", GITHUB_ORG)
    fmt.Printf("*** GITHUB_REPO    = %v\n", GITHUB_REPO)
    fmt.Printf("*** GITHUB_ISSUE   = %v\n", GITHUB_ISSUE)
    fmt.Printf("*** GITHUB_COMMENT = %v\n", GITHUB_COMMENT)
    fmt.Printf("*** GITHUB_USERS   = %v\n", GITHUB_USERS)
    fmt.Printf("*** GITHUB_USER    = %v\n", GITHUB_USER)

    cli := Github.NewClient()
    usr := cli.GetUser(GITHUB_USER)
    rep := cli.GetRepository(GITHUB_ORG, GITHUB_REPO)
    iss := rep.GetIssue(GITHUB_ISSUE)
    com := iss.GetComment(GITHUB_COMMENT)

    if trial.githubRepos {
        usr.PrintOrgs()
        usr.PrintRepos()
        cli.PrintOrgRepos(GITHUB_ORG)
    }

    if trial.githubIssues {
        rep.PrintIssues()
        txt := iss.Details()
        fmt.Printf("\n*** FETCH GITHUB %s/%s/%d:\n%s\n", GITHUB_ORG, GITHUB_REPO, GITHUB_ISSUE, txt)
    }

    if trial.githubComments {
        rep.PrintComments()
        txt := com.Details()
        fmt.Printf("\n*** FETCH GITHUB %s/%s/%d/comment:\n%s\n", GITHUB_ORG, GITHUB_REPO, GITHUB_ISSUE, txt)
    }

    if trial.githubUsers {
        for _, u := range GITHUB_USERS {
            usr := cli.GetUser(u)
            txt := usr.Details()
            fmt.Printf("\n*** FETCH GITHUB user %q:\n%s\n", u, txt)
        }
    }

    if trial.githubRate {
        limit := Github.GetRateLimit(nil)
        fmt.Printf("\n*** GITHUB RATE LIMITS:\n")
        fmt.Printf("Core:    %v\n", limit.Core)
        fmt.Printf("Search:  %v\n", limit.Search)
        fmt.Printf("GraphQL: %v\n", limit.GraphQL)
    }

    return true
}

// ============================================================================
// Transfer testing
// ============================================================================

// Exercise the GitHub GraphQL API.
//  NOTE: limited to TRIAL_REPOS if no repos are given.
func TrialGraphQL(repos ...string) bool {
    trialTitle("GitHub GraphQL")
    repos = RepositoryNames(repos...)
    if len(repos) == 0 { repos = TRIAL_REPOS }
    Github.TrialGraphQL(repos...)
    return true
}

// ============================================================================
// Transfer testing
// ============================================================================

// Exercise transfer logic without actually creating GitHub issues/comments.
//  NOTE: limited to TRIAL_PROJECTS if no projectKeys are given.
func TrialTransfer(projectKeys ...string) bool {
    trialTitle("Transfer")
    if len(projectKeys) == 0 { projectKeys = TRIAL_PROJECTS }
    FakeTransfer = true
    TransferAll(projectKeys...)
    return true
}

// ============================================================================
// Internal functions
// ============================================================================

// Output a section title to visually separate trial outputs.
func trialTitle(name string) {
    char   := "="
    span   := strings.Repeat(char, 80)
    cols   := util.CharCount(span)
    side   := strings.Repeat(char, 3)
    indent := util.CharCount(side)
    name   += " trials"
    width  := util.CharCount(name)
    spaces := (cols - width - (2 * indent)) / 2
    gap    := strings.Repeat(" ", spaces)
    title  := side + gap + name + gap
    if util.CharCount(title) + indent < cols { title += " " }
    title += side
    fmt.Println(span)
    fmt.Println(span)
    fmt.Println(title)
    fmt.Println(span)
    fmt.Println(span)
}
