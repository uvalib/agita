// args.go
//
// Command line arguments.

package main

import (
	"flag"
	"os"
	"strings"

	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Constants
// ============================================================================

// Program operational mode.
const (
    ModeNone     = 0
    ModeTransfer = 1 << iota
    ModeExport   = 1 << iota
    ModeClear    = 1 << iota
    ModeTrial    = 1 << iota
    ModeHelp     = 1 << iota
)

// ============================================================================
// Variables
// ============================================================================

// Program operational mode.
var Mode = ModeHelp

// List of Jira project source(s) or GitHub repos for "-clear".
var Args = []string{}

// ============================================================================
// Functions
// ============================================================================

// Set Mode and Projects based on command line arguments.
func GetArgs() {
    // In launch.json, using "${input:...}" with "args" results in a single
    // string which must be tokenized here.
    if util.InDebugger() && (len(os.Args) > 1) {
        args := []string{}
        for _, arg := range os.Args[1:] {
            args = append(args, strings.Fields(arg)...)
        }
        os.Args = append(os.Args[0:1], args...)
    }

    xfer   := flag.Bool("transfer", false, "Create GitHub issues and comments from Jira issues and comments.")
    export := flag.Bool("export",   false, "Generate JSON from Jira projects, issues, and comments.")
    clear  := flag.Bool("clear",    false, "Remove GitHub issues and comments.")
    trial  := flag.Bool("trial",    false, "Exercise Jira and GitHub APIs; see below.")
    help   := flag.Bool("help",     false, "Show program usage help.")

    flag.Usage = showUsage
    flag.Parse()

    if Args = flag.Args(); len(Args) > 0 {
        for _, arg := range Args {
            if strings.HasPrefix(arg, "-") {
                abort("flags not acceptable after project list")
            }
        }
    }

    mode := ModeNone
    if *xfer   { mode = mode | ModeTransfer }
    if *export { mode = mode | ModeExport }
    if *clear  { mode = mode | ModeClear }
    if *trial  { mode = mode | ModeTrial }
    if *help   { mode = mode | ModeHelp }
    if mode != ModeNone {
        Mode = mode
    }

    switch Mode {
        case ModeTransfer:  // ok
        case ModeExport:    // ok
        case ModeClear:     // ok
        case ModeTrial:     // ok
        case ModeHelp:      usage(NORMAL_EXIT)
        case ModeNone:      abort("no default mode defined")
        default:            abort("only one mode flag is acceptable")
    }
}

// ============================================================================
// Internal functions
// ============================================================================

// Print a message on stderr and exit.
func abort(msg string, arg ...any) {
    ShowError(msg, arg...)
    usage(ABORT_EXIT)
}

// Print program usage help on stderr and exit.
func usage(code int) {
    flag.Usage()
    os.Exit(code)
}

// Print program usage help on stderr.
func showUsage() {
    prog := util.Progname()

    Show("Usage: %s mode      names...", prog)
    Show("Usage: %s -transfer %s | Jira_projects...", prog, ALL_PROJECTS)
    Show("Usage: %s -export   %s | Jira_projects...", prog, ALL_PROJECTS)
    Show("Usage: %s -clear    %s | GitHub_repos...",  prog, ALL_REPOS)
    Show("Usage: %s -trial    [args...]", prog)
    Show("Usage: %s -help", prog)
    Show("")
    Show("Mode Flags:")

    flag.PrintDefaults()

    Show("")
    Show("Exactly one mode flag must be given.")
    Show("If no projects are specified then all projects are assumed.")

    Show("")
    Show("Trial arguments:")
    Show("\tall                                  - All trials")
    Show("\tjira[:projects,issues,comments]      - Demo Jira API")
    Show("\tgithub[:users,repos,issues,comments] - Demo GitHub API")
    Show("\tgraphql                              - Demo GitHub GraphQL API")
    Show("\ttransfer [Jira_projects...]          - Simulate a transfer")

    Show("")
    Show("Each element of `Jira_projects` may be one of:")
    Show("\tPROJ               - All issues from Jira project PROJ")
    Show("\tPROJ-min           - PROJ issues starting with PROJ-min")
    Show("\tPROJ-min PROJ-max  - PROJ issues in the range [PROJ-min,PROJ-max]")

    Show("")
}
