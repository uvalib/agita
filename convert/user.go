// convert/user.go
//
// Mapping of Jira user accounts to GitHub user accounts.

package convert

// ============================================================================
// Exported variables
// ============================================================================

// A mapping of Jira account to its matching GitHub account.
var JiraToGithubUser = map[string]string {

    // === Jira users with GitHub accounts

    "akl3b":    "amber-reichert",   // Amber Reichert
    "bcb4y":    "bc-butler",        // Brandon Butler
    "cmw5mc":   "ChristopherWelte", // Christopher Welte
    "dhc4z":    "dougchestnut",     // Doug Chestnut
    "dpg3k":    "UVADave",          // Dave Goldstein
    "ecr2c":    "ecr2c",            // Ellen Ramsey
    "jkb2b":    "jeremyboggs",      // Jeremy Boggs
    "jlj5aj":   "jlj5aj",           // Jason Jordan
    "jlk4p":    "jlk4p",            // Jack Kelly
    "jph9e":    "jph9e",            // Jeff Hill
    "jtb4t":    "jtbmadva",         // Jeremy Bartczak
    "lar4k":    "larockey",         // Leigh Rockey
    "ldc8n":    "ldc8n",            // Linda Vaughan
    "lf6f":     "louffoster",       // Lou Foster
    "lw2cd":    "lwforpres",        // Lauren Work
    "md5wz":    "mikedurbin",       // Mike Durbin
    "mhw8m":    "mhw8m",            // Mark Witteman
    "naw4t":    "nestorw",          // Nestor Walker
    "rar6u":    "rallisonr",        // Renee Reighart
    "rh9ec":    "haschart",         // Bob Haschart
    "rmg6f":    "rmgartland",       // Bob Gartland
    "rwl":      "RayLubinsky",      // Ray Lubinsky
    "sah":      "shlake",           // Sherry Lake
    "sss3s":    "SusanSNeal",       // Susan Neal
    "tss6n":    "UncleJefferson",   // Tim Stevens
    "wxf2nb":   "ruobingsu",        // Ruobing Su
    "xw5d":     "Xiaoming",         // Xiaoming Wang
    "ys2n":     "ys2n",             // Yuji Shinozaki

    // === Jira users without GitHub accounts

    "adj5j":    "",                 // Tony Jones
    "ag5vy":    "",                 // Anne Gaynor
    "arb5w":    "",                 // Ann Burns
    "au3vk":    "",                 // Ashish Upadhyaya
    "bds2mv":   "",                 // Ben Spector
    "bg9ba":    "",                 // Brenda Gunn
    "bga3d":    "",                 // Bethany Anderson
    "bmx7wf":   "",                 // Marco Battistella
    "bwb9f":    "",                 // Beth Blanton
    "bwo3db":   "",                 // Ben Ormond
    "cma4u":    "",                 // Carla Arton
    "cmm2t":    "",                 // Christina Deane
    "csk5vf":   "",                 // Connor Kenaston
    "ege5vd":   "",                 // Molly Fair
    "eks4x":    "",                 // Eric Seidel
    "emg3b":    "",                 // Eliza Gilligan
    "er8fn":    "",                 // Liz Rapp
    "gmk3d":    "",                 // Ginny Kois
    "hdn9c":    "",                 // Hanni Nabahe
    "hmh8eq":   "",                 // Akari Hernandez
    "ja3nu":    "",                 // Joseph Azizi
    "jgc4q":    "",                 // Zeke Crater
    "jh9xp":    "",                 // Jill Heinze
    "jmu2m":    "",                 // John Unsworth
    "jor2a":    "",                 // Jennifer Roper
    "ka7uz":    "",                 // Krystal Appiah
    "khj5c":    "",                 // Kristen Jensen
    "kmm6ef":   "",                 // Kara McClurken
    "libx-sat": "",                 // SAT Team
    "lsc6v":    "",                 // Lorrie Chisholm
    "lws4n":    "",                 // Lucie Stylianopoulos
    "lzl9b":    "",                 // Lauren Longwell
    "mdw7p":    "",                 // Melanie Williams
    "mhm8m":    "",                 // Heather Riser
    "mrc2x":    "",                 // Marc Campbell
    "mwm7b":    "",                 // Molly Minturn
    "nir4x":    "",                 // Nicole Royal
    "np6hd":    "",                 // Nitesh Parajuli
    "nqt2bq":   "",                 // Ryan Russo
    "pw7e":     "",                 // Peter Welch
    "rac8f":    "",                 // Rebecca Coleman
    "rah6w":    "",                 // Ric Hodges
    "rcm7e":    "",                 // Rennie Mapp
    "rds4w":    "",                 // Rob Smith
    "rfg":      "",                 // Rebecca Garver
    "rpk2kn":   "",                 // Ryan Kelly
    "rsl6m":    "",                 // Robin Ruggaber
    "sd3gz":    "",                 // Sue Donovan
    "sh2rd":    "",                 // Stephanie Hunter
    "smm3m":    "",                 // Sean McCord
    "sp7fg":    "",                 // Sony Prosper
    "spr7b":    "",                 // Sue Richeson
    "srv3f":    "",                 // Steven Villereal
    "stg2s":    "",                 // Stan Gunn
    "tnn7yc":   "",                 // Tho Nguyen
    "trb5f":    "",                 // Tammy Barbour
    "tsh2k":    "",                 // Trillian Hosticka
    "tvf3c":    "",                 // Tracy Fewell
    "vdg8v":    "",                 // Dave Griles
    "vlk4n":    "",                 // Vince Kois
    "wdw5ch":   "",                 // Will Wyatt
    "wkb5j":    "",                 // Winston Barham
    "wmr5a":    "",                 // Will Rourk
    "xjk5yt":   "",                 // Nora Wilkerson
    "ys2ck":    "",                 // Yukesh Sitoula

    // === Ignored GitHub accounts

    // "":      "lib-ole-ci",       // [non-user]
    // "":      "mfarish",          // Mitchell Farish
    // "":      "uvabamboo",        // [non-user]
    // "":      "UVABuilder",       // [non-user]
}
