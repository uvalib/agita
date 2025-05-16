// convert/user.go
//
// Mapping of Jira projects to GitHub repositories.

package convert

import (
	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"lib.virginia.edu/agita/Github"
	"lib.virginia.edu/agita/Jira"
)

// ============================================================================
// Exported variables
// ============================================================================

// A mapping of Jira project key to its matching GitHub repository name.
// Each is annotated with its number of issues as of 2025-05-15.
var ProjectToRepo = map[string]string {
//  Jira project    GitHub repo    Jira project name                    Issue Count
//  ------------    -----------    ---------------------------------    -----------
    "AST":          "avalon",   // Avalon Service Team                  38
    "CSH":          "",         // Computer Systems Help                9828
    "DSP":          "",         // Desktop Support Projects             57
    "DHUV":         "",         // DH Portal                            179
    "DCP":          "",         // Digital Collections Projects         202
    "DEPP":         "",         // Digital Exhibits Platform Project    25
    "DPGPROJ":      "",         // dpg3k-project                        112
    "EMMA":         "emma",     // EMMA                                 149
    "EMHELP":       "",         // EMMA Help                            427
    "LIBRAONE":     "libra-oa", // Libra 1.x                            30
    "LIBRA":        "Libra2",   // Libra 2.x                            1141
    "LIBRADATA":    "",         // Libra Data                           141
    "LIBRAOC":      "libra-oc", // Libra OC                             608
    "MA2":          "",         // Mandala 2.0                          6
    "MLPS":         "",         // Misc LIT Project Support             39
    "PPK":          "",         // Preservation Projects                47
    "RWL":          "",         // RWL                                  53
    "SAT":          "",         // Search & Access Technologies         71
    "SP":           "",         // Sysadmin Projects                    25
    "TDG":          "",         // Technology Development Group         1931
    "ULBE":         "",         // UVA Library Browser Extension        97
    "VIRGO":        "virgo",    // Virgo 3.x                            403
    "VIRGONEW":     "",         // Virgo 4.x                            2810
    "VIRGOEXP":     "virgo4",   // Virgo 4.x (previous)                 91
    "WAWG":         "",         // Web Archiving Working Group          1309
}

// A mapping of GitHub repository name to its matching Jira project key.
var RepoToProject = util.MapInvert(util.MapCompact(ProjectToRepo))

// ============================================================================
// Exported functions
// ============================================================================

// Find or infer the GitHub repository name for the given project name.
func RepositoryNameFor(name string) string {
    switch {
        case name == "":                return name
        case RepoToProject[name] != "": return name
        case ProjectToRepo[name] != "": return ProjectToRepo[name]
        default:                        return ProjectRepositoryFor(name)
    }
}

// Give the GitHub project repository name for the given Jira project name.
func ProjectRepositoryFor(jiraKey string) string {
    return Github.PROJ_NAME_PREFIX + jiraKey
}

// Render a Jira project into JSON.
func ProjectToJson(src Jira.Project) string {
    if bytes, err := src.MarshalJSON(); log.ErrorValue(err) == nil {
        return string(bytes)
    } else {
        return ""
    }
}
