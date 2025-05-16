// markdown/jira.go
//
// Convert Jira markdown.

package markdown

import (
	"fmt"
	"strings"

	"lib.virginia.edu/agita/re"
	"lib.virginia.edu/agita/util"
)

// ============================================================================
// Exported constants
// ============================================================================

// Markdown resulting in an HTML "<hr>".
const HORIZONTAL_RULE = "----"

const CODE_PREFIX  = ""
const NOFMT_PREFIX = ""
const QUOTE_PREFIX = "> "
const PANEL_PREFIX = "> "

// Temporary construct which indicates an embedded newline in a Jira table cell
// which needs to be converted to "<br>" for the GitHub table.
const CELL_NEWLINE = "[[[NEWLINE]]]"

// ============================================================================
// Exported functions
// ============================================================================

// Convert Jira markdown to Github markdown.
//
// Caveats:
//
// {code} - Style attributes ignored because GitHub sanitizes them.
//
// {panel} - There is no GitHub equivalent but it is "faked" here using
// horizontal rules.  Style attributes ignored because GitHub sanitizes them.
//
func JiraToGithub(src string) string {
    var inCodeBlock, inNoformat, inQuote, inPanel, inTable, inTableCell bool
    depthCounter := map[int]int{}
    lastDepth    := 0
    indent       := ""
    resetIndent  := func() { lastDepth, indent = 0, "" }
    color        := ""
    lines        := strings.Split(src, "\n")

    for i, line := range lines {
        line = util.StripRight(line)
        blockEnd := ""

        // States where there is no further formatting of the line unless a
        // portion of the line comes after the block termination marker.

        if inCodeBlock {
            if marker := "{code}"; strings.Contains(line, marker) {
                inCodeBlock = false
                lines[i] = "```"
                pat := re.New(`{code[^}]*}`)
                line = jiraEndBlock(&lines[i], line, pat, CODE_PREFIX)
            } else {
                line = ""
            }
            if line == "" {
                continue
            } else {
                blockEnd = lines[i]
            }
        }
        if inNoformat {
            if marker := "{noformat}"; strings.Contains(line, marker) {
                inNoformat = false
                lines[i] = "```"
                line = jiraEndBlock(&lines[i], line, marker, NOFMT_PREFIX)
            } else {
                line = ""
            }
            if line == "" {
                continue
            } else {
                blockEnd = lines[i]
            }
        }

        // States where there is further formatting of the line unless it is
        // the terminating token for the state.

        if inQuote {
            if marker := "{quote}"; strings.Contains(line, marker) {
                inQuote  = false
                lines[i] = ""
                line = jiraEndBlock(&lines[i], line, marker, QUOTE_PREFIX)
                if line == "" {
                    continue
                } else {
                    blockEnd = lines[i]
                }
            }
        }
        if inPanel {
            if marker := "{panel}"; strings.Contains(line, marker) {
                inPanel  = false
                lines[i] = HORIZONTAL_RULE
                line = jiraEndBlock(&lines[i], line, marker, PANEL_PREFIX)
                if line == "" {
                    continue
                } else {
                    blockEnd = lines[i]
                }
            }
        }

        // Table lines are handled specially here.

        if inTableCell {
            if strings.HasSuffix(line, "|") {
                if strings.HasSuffix(line, "||") {
                    lines[i] = strings.TrimSuffix(line, "|")
                } else {
                    lines[i] = line
                }
                inTableCell = false
            } else {
                lines[i] = line + CELL_NEWLINE
            }
            continue

        } else if inTable {
            if strings.HasPrefix(line, "|") {
                lines[i], inTableCell = jiraTableLine(line, false)
                continue
            } else {
                inTable = false
            }
        }

        // Process the line for state changes where a block begins with a
        // line containing a token like "{noformat}" which directs the
        // interpretation of ensuing lines until the closing token is found.

        if (len(line) == 0) && (lastDepth == 1) {
            // Reset numbering depth on blank line
            resetIndent()
            lines[i] = line
            continue

        } else if (len(line) < 3) {
            // No line subsititutions applicable.
            lines[i] = line
            continue

        } else if strings.Contains(line, "{noformat}") {
            // Start of a {noformat} block.
            inNoformat = true
            lines[i] = "```"
            continue

        } else if strings.Contains(line, "{quote}") {
            // Start of a {quote} block.
            inQuote = true
            lines[i] = ""
            continue

        } else if strings.Contains(line, "{panel}") {
            // Start of a {panel} block with no attributes.
            inPanel = true
            lines[i] = "\n" + HORIZONTAL_RULE
            continue

        } else if strings.Contains(line, "{panel:") {
            // Start of a {panel} block with attributes.
            inPanel = true
            param := re.ReplaceAll(line, `.*{panel:(.+)}.*`, "$1")
            title := ""
            for part := range strings.SplitSeq(param, "|") {
                if strings.HasPrefix(part, "title=") {
                    title = strings.TrimPrefix(part, "title=")
                    title = re.ReplaceAll(title, `^"(.*)"$`, "$1")
                }
            }
            lines[i] = "\n" + HORIZONTAL_RULE
            if title != "" {
                lines[i] += "\n> " + title
                lines[i] += "\n> " + HORIZONTAL_RULE
            }
            continue

        } else if re.Match(line, `^\s*{code}$`) {
            // Start of a {code} block with no language specified.
            inCodeBlock = true
            lines[i] = "```"
            continue

        } else if re.Match(line, `^\s*{code:[^}]+}$`) {
            // Start of a {code} block with attributes (possibly language).
            inCodeBlock = true
            title := ""
            lang  := ""
            param := re.ReplaceAll(line, `.*{code:(.*)}.*`, "$1")
            for part := range strings.SplitSeq(param, "|") {
                if !strings.Contains(part, "=") {
                    lang  = part
                } else if strings.HasPrefix(part, "title=") {
                    title = strings.TrimPrefix(part, "title=")
                    title = re.ReplaceAll(title, `^"(.*)"$`, "$1")
                }
            }
            lines[i] = re.ReplaceAll(line, `{code:(.*)}`, ("```" + lang))
            if title != "" {
                lines[i] = "\n**" + title + "**\n" + lines[i]
            }
            continue

        } else if re.Match(line, `^\s*{color:[^}]+}$`) {
            // Begin a set of colored text lines.
            color = re.ReplaceAll(line, `^\s*{color:(.+)}`, "$1")
            lines[i] = ""
            continue

        } else if strings.Contains(line, "{color:") {
            // Embedded colors are handled below.

        } else if re.Match(line, `{color}`) {
            // End colored text lines, possibly with uncolored text following.
            part := strings.SplitN(line, "{color}", 2)
            blockColor := color
            blockEnd, line = part[0], part[1]
            color = ""
            if blockEnd != "" {
                blockEnd = jiraInline(blockEnd, blockColor)
            } else if line == "" {
                continue
            }
        }

        // Process the line.

        if strings.HasPrefix(line, "bq.") {
            // A single-line block quote.
            lines[i] = QUOTE_PREFIX + jiraInline(line[3:], color)

        } else if strings.HasPrefix(line, "||") {
            // First row of a table.
            inTable = true
            lines[i], inTableCell = jiraTableLine(line, true)

        } else if re.Match(line, `^h\d\.`) {
            // Transform "h1." through "h6." headers.
            resetIndent()
            switch line[1] {
                case '1': line = "# "      + line[3:]
                case '2': line = "## "     + line[3:]
                case '3': line = "### "    + line[3:]
                case '4': line = "#### "   + line[3:]
                case '5': line = "##### "  + line[3:]
                case '6': line = "###### " + line[3:]
            }
            lines[i] = jiraInline(line, color)

        } else if re.Match(line, `^\s*[*-]+\s`) {
            // Transform unordered list lines.
            space := re.ReplaceAll(line, `^(\s*)[*-]+.*`, "$1")
            marks := re.ReplaceAll(line, `^\s*([*-]+).*`, "$1")
            depth := len(marks)
            chars := len(space) + depth
            text  := jiraInline(util.StripLeft(line[chars:]), color)
            indent = strings.Repeat("   ", depth-1)
            lines[i] = fmt.Sprintf("%s%c %s", indent, marks[0], text)

        } else if re.Match(line, `^\s*#+\s`) {
            // Transform autonumbered list lines.
            space := re.ReplaceAll(line, `^(\s*)#+.*`, "$1")
            marks := re.ReplaceAll(line, `^\s*(#+).*`, "$1")
            depth := len(marks)
            chars := len(space) + depth
            text  := jiraInline(util.StripLeft(line[chars:]), color)
            indent = strings.Repeat("   ", depth-1)
            if len(space) > len(indent) {
                indent = space
            }
            count := depthCounter[depth] + 1
            lines[i] = fmt.Sprintf("%s%d. %s", indent, count, text)
            if lastDepth > depth {
                depthCounter[lastDepth] = 0
            }
            depthCounter[depth] = count
            lastDepth = depth

        } else {
            lines[i] = jiraInline(line, color)
            if lastDepth > 0 {
                num := fmt.Sprintf("%d. ", depthCounter[lastDepth])
                add := strings.Repeat(" ", len(num))
                lines[i] = indent + add + lines[i]
            }
        }

        // State-based post-processing.

        switch {
            case inPanel: lines[i] = PANEL_PREFIX + lines[i]
            case inQuote: lines[i] = QUOTE_PREFIX + lines[i]
        }

        // Restore saved content which terminates a block.

        if blockEnd != "" {
            lines[i] = blockEnd + "\n" + lines[i]
        }
    }

    // Convert embedded Jira table cell newlines into the GitHub equivalent.

    result := strings.Join(lines, "\n")
    return strings.ReplaceAll(result, (CELL_NEWLINE + "\n"), "<br>")
}

// ============================================================================
// Internal functions
// ============================================================================

// Accounts for the possibility that the block-terminating marker is not on its
// own line by spliting `text` by `marker`.
//
// * The `linePtr` is assumed to already contain the proper GitHub sequence to
// terminate the block.
//
// * Any non-blank text before `marker` is injected prior to the termination so
// that it is part of the block.
//
// * Any non-blank text after `marker` is returned.
//
func jiraEndBlock[T string|re.Regex|*re.Regex](linePtr *string, text string, marker T, prefix string) string {
    var part []string
    switch v := any(marker).(type) {
        case string:    part = strings.SplitN(text, v, 2)
        case re.Regex:  part = v.Split(text, 2)
        case *re.Regex: part = v.Split(text, 2)
    }
    if inBlock := util.Strip(part[0]); inBlock != "" {
        *linePtr = prefix + inBlock + "\n" + *linePtr
    }
    return util.Strip(part[1])
}

// Parse a line already determined to be part of a Jira markdown table
// (i.e., a line that starts with '|').
//
// If this is known to be the first row of a table, a divider is inserted to
// ensure that GitHub sees it as a table.
//
// Any cell of a Jira table may start with "||" which makes it a "header cell".
//
// Since table cells can contain {panel} and other blocks, each cell content is
// processed through JiraToGithub().
//
func jiraTableLine(line string, firstRow bool) (string, bool) {
    headerCount := 0
    cellCount   := 0
    ended       := false
    row         := "|"

    // For each table cell...
    for rdr := strings.NewReader(line); ; {
        // If the cell does not begin with '|' or "||" then we must be at the
        // end of the line.
        if char1, _, err := rdr.ReadRune(); (err != nil) || (char1 != '|') {
            break
        }

        // If the cell is a header (i.e., begins with "||") then the next
        // character should be '|'; otherwise it is the first content character
        // of the cell.
        char2, _, err := rdr.ReadRune()
        if err != nil {
            break
        }

        // At this point the string reader is in a cell.
        cellCount++
        headerCell := (char2 == '|')
        if headerCell {
            headerCount++
        } else {
            rdr.UnreadRune()
        }

        // Read cell characters until the next cell or end of line is reached.
        isCell := true
        chars  := []rune{}
        for {
            if char, _, err := rdr.ReadRune(); err != nil {
                // Final cell ends with no '|'.
                ended  = (len(chars) == 0)
                isCell = !ended
                break
            } else if char == '|' {
                // Start of next cell.
                ended = true
                rdr.UnreadRune()
                break
            } else {
                // Accumulate cell content.
                chars = append(chars, char)
            }
        }

        // Append the cell to the output row.
        if isCell {
            cell := string(chars)
            if ended {
                cell = JiraToGithub(cell)
                if headerCell && !strings.HasPrefix(cell, "**") {
                    cell = "**" + cell + "**"
                }
                row += " " + cell + " |"
            } else {
                row += " " + cell + CELL_NEWLINE
            }
        } else {
            cellCount--
            if headerCell {
                headerCount--
            }
        }
    }

    // If this first table line is all header cells then append the divider.
    // Otherwise, a bogus header row needs to be prepended in order to cause
    // GitHub to see it as a table.
    if firstRow {
        divider := "|" + strings.Repeat(" --- |", cellCount)
        if headerCount == cellCount {
            row = row + "\n" + divider
        } else {
            heading := "|" + strings.Repeat(" |", cellCount)
            row = heading + "\n" + divider + "\n" + row
        }
    }

    return row, !ended
}

// Perform in-line substitutions of Jira Markdown forms with Github Markdown.
func jiraInline(line, color string) string {
    line = jiraInlineStrong(line)
    line = jiraInlineEmphasis(line)
    line = jiraInlineMonospaced(line)
    line = jiraInlineDeleted(line)
    line = jiraInlineInserted(line)
    line = jiraInlineSuperscript(line)
    line = jiraInlineSubscript(line)
    line = jiraInlineCitation(line)
    line = jiraInlineHyperlink(line)
    line = jiraInlineCode(line)
    if color == "" {
        return jiraInlineColor(line)
    } else {
        return githubColorize(line, color)
    }
}

// Transform Jira bold markdown.
//  EXAMPLE: *strong* => **strong**
func jiraInlineStrong(line string) string {
    return re.ReplaceAll(line, `(^|\s+)\*(.+?)\*(\s+|$)`, "$1**$2**$3")
}

// Transform Jira italic markdown.
//  EXAMPLE: _emphasis_ => *emphasis*
func jiraInlineEmphasis(line string) string {
    return re.ReplaceAll(line, `(^|\s+)_(.+?)_(\s+|$)`, "$1*$2*$3")
}

// Transform Jira italic markdown.
//  EXAMPLE: {{monospaced}} => `monospaced`
func jiraInlineMonospaced(line string) string {
    return re.ReplaceAll(line, `{{(.+?)}}`, "`$1`")
}

// Transform Jira deleted markdown.
//  EXAMPLE: -deleted- => ~~deleted~~
func jiraInlineDeleted(line string) string {
    // return re.ReplaceAll(line, `(^|\s+)\-(.+?)\-(\s+|$)`, "$1~~$2~~$3")
    return line
}

// Transform Jira inserted markdown.
//  EXAMPLE: +inserted+ => <ins>inserted</ins>
func jiraInlineInserted(line string) string {
    // return re.ReplaceAll(line, `(^|\s+)\+(.+?)\+(\s+|$)`, "$1<ins>$2</ins>$3")
    return line
}

// Transform Jira superscript markdown.
//  EXAMPLE: ^superscript^ => <sup>superscript</sup>
func jiraInlineSuperscript(line string) string {
    return re.ReplaceAll(line, `\^(.+?)\^`, "<sup>$1</sup>")
}

// Transform Jira subscript markdown.
//  EXAMPLE: ~subscript~ => <sub>subscript</sub>
func jiraInlineSubscript(line string) string {
    return re.ReplaceAll(line, `(^|[^<])\~(\S+?)\~`, "$1<sub>$2</sub>")
}

// Transform Jira citation markdown.
//  EXAMPLE: ??citation?? => <cite>citation</cite>
func jiraInlineCitation(line string) string {
    return re.ReplaceAll(line, `\?\?(.+?)\?\?`, "<cite>$1</cite>")
}

// Jira hyperlink markdown.
//  EXAMPLE: [https://x.com]     => <https://x.com>
//  EXAMPLE: [XYZ|https://x.com] => [XYZ](https://x.com)
func jiraInlineHyperlink(line string) string {
    line = re.ReplaceAll(line, `\[([^|]+?)\]`,     "<$1>")
    return re.ReplaceAll(line, `\[(.+?)\|(.+?)\]`, "[$1]($2)")
}

// Transform all Jira inline code directives in the line.
//  EXAMPLE: {code:javascript}stuff{code} => `stuff`
func jiraInlineCode(line string) string {
    return re.ReplaceAll(line, `{code[^}]*}(.*?){code}`, "`$1`")
}

// Transform all Jira inline color directives in the line.
//  EXAMPLE: {color:red}red text{color} => the Github LaTeX equivalent
func jiraInlineColor(line string) string {
    if re.Match(line, `{color:.*{color}`) {
        repl := func(match string) string {
            text  := re.ReplaceAll(match, `{color:.+}(.*){color}`, "$1")
            color := re.ReplaceAll(match, `{color:(.+)}.*{color}`, "$1")
            return githubColorize(text, color)
        }
        line = re.ReplaceAllFunc(line, `{color:([^}]+)}.*?{color}`, repl)
    }
    return line
}

// Apply the LaTeX sequence honored by GitHub which colorizes text.
// In the case of blank text or text colored "black" no colorization is done.
func githubColorize(text, color string) string {
    if (color != "black") && (color != "#000000") && !re.Match(text, `^\s*$`) {
        // Escape sequences GitHub LaTeX doesn't like.
        text = re.ReplaceAll(text, `([_#*])`, `\\\$1`)
        text = fmt.Sprintf(`$\color{%s}{\textsf{%s}}$`, color, text)
    }
    return text
}
