// Jira/attachment.go
//
// Issue/comment attachment files.

package Jira

import (
	"io"

	"lib.virginia.edu/agita/log"
)

// ============================================================================
// Exported functions
// ============================================================================

// Get the content of the indicated attachment.
func DownloadAttachment(client *Client, attachmentId string) string {
    if client == nil {
        client = MainClient()
    }
    rsp, err := client.ptr.Issue.DownloadAttachment(attachmentId)
    if log.ErrorValue(err) == nil {
        defer rsp.Body.Close()
        bytes, err := io.ReadAll(rsp.Body)
        if log.ErrorValue(err) == nil {
            return string(bytes)
        }
    }
    return ""
}
