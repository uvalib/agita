// convert/user.go
//
// Mapping of Jira projects to GitHub repositories.

package convert

import (
	"lib.virginia.edu/agita/log"

	"lib.virginia.edu/agita/Jira"
)

// ============================================================================
// Exported variables
// ============================================================================

// A mapping of Jira project key to its matching GitHub repository name.
var ProjectToRepo = map[string]string {
    "AST":          "avalon",   // Avalon Service Team
    "CSH":          "",         // Computer Systems Help
    "DSP":          "",         // Desktop Support Projects
    "DHUV":         "",         // DH Portal
    "DCP":          "",         // Digital Collections Projects
    "DEPP":         "",         // Digital Exhibits Platform Project
    "DPGPROJ":      "",         // dpg3k-project
    "EMMA":         "emma",     // EMMA
    "EMHELP":       "",         // EMMA Help
    "LIBRAONE":     "libra-oa", // Libra 1.x
    "LIBRA":        "Libra2",   // Libra 2.x
    "LIBRADATA":    "",         // Libra Data
    "LIBRAOC":      "libra-oc", // Libra OC
    "MA2":          "",         // Mandala 2.0
    "MLPS":         "",         // Misc LIT Project Support
    "PPK":          "",         // Preservation Projects
    "RWL":          "",         // RWL
    "SAT":          "",         // Search & Access Technologies
    "SP":           "",         // Sysadmin Projects
    "TDG":          "",         // Technology Development Group
    "ULBE":         "",         // UVA Library Browser Extension
    "VIRGO":        "virgo",    // Virgo 3.x
    "VIRGONEW":     "",         // Virgo 4.x
    "VIRGOEXP":     "virgo4",   // Virgo 4.x (previous)
    "WAWG":         "",         // Web Archiving Working Group
}

// ============================================================================
// Exported functions
// ============================================================================

// Render a Jira project into JSON.
func ProjectToJson(src Jira.Project) string {
    if bytes, err := src.MarshalJSON(); log.ErrorValue(err) == nil {
        return string(bytes)
    } else {
        return ""
    }
}
