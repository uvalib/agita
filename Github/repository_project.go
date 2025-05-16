// Github/repository_project.go
//
// Creation of project-oriented "issues-only" repositories.

package Github

import (
	"fmt"
	"strings"
	"time"

	"lib.virginia.edu/agita/log"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported constants
// ============================================================================

// Project repository property.
const (
    PROJ_NAME_PREFIX = "project-"
    PROJ_DESC_FORMAT = "Preserved issues and comments for Jira project %s."
)

// Template project repository property.
const (
    TEMPLATE_PROJ_NAME = "agita-proj-template"
    TEMPLATE_PROJ_FULL = "AGITA project template repository"
    TEMPLATE_PROJ_DESC = "Used as a source for generating issues-only repositories."
)

// Folder under the project repository where copied attachments are stored.
const ATTACH_DIR = "attachments"

// ============================================================================
// Exported functions
// ============================================================================

// Get an issues-only repository, creating it if necessary.
func GetProjRepo(client *Client, name string) (result *Repository) {
    if result = GetRepository(client, ORG, name, true); result == nil {
        result = CreateProjRepo(client, name)
    }
    return
}

// Create a new issues-only repository from the proj template repository.
func CreateProjRepo(client *Client, name string) *Repository {
    return createProjRepository(client, name)
}

// Add a copy of an attachment to an existing issues-only repository.
func CreateProjAttachment(client *Client, name, file, content string) {
    createProjAttachment(client, name, file, content)
}

// ============================================================================
// Exported functions - project repository
// ============================================================================

// Generate data for a new issues-only repository.
func ProjRepoTemplateData(name string) *github.TemplateRepoRequest {
    pref := PROJ_NAME_PREFIX
    if !strings.HasPrefix(name, pref) {
        panic(fmt.Sprintf("name %q missing prefix %q", name, pref))
    }
    proj := strings.ToUpper(strings.TrimPrefix(name, pref))
    desc := fmt.Sprintf(PROJ_DESC_FORMAT, proj)
    return RepoTemplateData(name, desc)
}

// Generate data for a new issues-only repository.
func ProjRepoTemplateDataAsRepository(name string) *github.Repository {
    data := ProjRepoTemplateData(name)
    return RepoTemplateDataAsRepository(data)
}

// ============================================================================
// Internal functions - project repository
// ============================================================================

// Number of tries to get the contents of the project's README.md after it has
// been created.
const maxGetReadmeRetries = 10

// GitHub project label for a project repository.
const projRepositoryLabel = "jira-project"

// Create a new issues-only repository from the proj template repository.
func createProjRepository(client *Client, name string) *Repository {
    getProjTemplateRepository(client)
    data := ProjRepoTemplateData(name)
    repo := createRepoFromTemplate(client, TEMPLATE_PROJ_NAME, data)
    if repo != nil {
        owner := ORG

        // Customize the new repository's README.md.  First, the hash of the
        // hash of the copied README.md must be obtained, which may take
        // repeated retries until because copying from the template repository
        // is asynchronous.
        file := "README.md"
        srv  := client.ptr.Repositories
        var current *github.RepositoryContent
        var rsp *github.Response
        var err error
        for range maxGetReadmeRetries {
            current, _, rsp, err = srv.GetContents(ctx, owner, name, file, nil)
            extractRateLimit(rsp)
            if err != nil {
                if strings.Contains(err.Error(), "This repository is empty") {
                    time.Sleep(time.Second)
                    continue
                } else {
                    log.ErrorValue(err)
                }
            } else if current == nil {
                log.Error("Repositories.GetContents failed")
            }
            break
        }
        if current != nil {
            opts := &github.RepositoryContentFileOptions{
                Message: github.Ptr("Updated " + file),
                Content: []byte(generateReadmeContent(data)),
                SHA:     current.SHA,
            }
            _, rsp, err := srv.UpdateFile(ctx, owner, name, file, opts)
            extractRateLimit(rsp)
            log.ErrorValue(err)
        }

        // Update the repository with an identifying label.
        setRepositoryTopics(client.ptr, owner, name, projRepositoryLabel)

    }
    return repo
}

// The body of the project repository README.md which follows the description.
const readmeContent = `
## [Issues](../../issues)

Each issue title is prefixed with the original Jira issue key.
At the top of the issue description, lines with "` + ISSUE_ANNOTATION_TAG +
`" convey information that cannot be represented in the GitHub issue.

At the top of each comment body, lines with "` + COMMENT_ANNOTATION_TAG +
`" convey information that cannot be represented in the GitHub comment.

## [Attachments](../main/` + ATTACH_DIR + `)

Each attachment file name is the original Jira attachment file name prefixed
with the Jira issue key with which it is associated.
Note that while some of these are referenced directly by any issue or comment
body (e.g. as embedded images), not all of them are.
`

// Produce README.md content tailored to the supplied repository data.
func generateReadmeContent(data *github.TemplateRepoRequest) string {
    desc := *data.Description
    name := *data.Name
    base := strings.TrimPrefix(name, PROJ_NAME_PREFIX)
    text := fmt.Sprintf("# Jira Project %s\n%s", base, desc)
    return text + "\n" + readmeContent
}

// Add a copy of an attachment to an existing issues-only repository.
//  NOTE: There is a to-be-determined maximum attachment size; may fail with:
//  "Sorry, the file is too large to be processed. Consider creating/updating
//  the file in a local clone and pushing it to GitHub."
func createProjAttachment(client *Client, name, file, content string) {
    opts := &github.RepositoryContentFileOptions{
        Message: github.Ptr("Stored attachment " + file),
        Content: []byte(content),
    }
    file = ATTACH_DIR + "/" + file
    _, rsp, err := client.ptr.Repositories.CreateFile(ctx, ORG, name, file, opts)
    extractRateLimit(rsp)
    log.ErrorValue(err)
}

// ============================================================================
// Internal variables - template repository for projects
// ============================================================================

// GitHub repository used as a template for generating project repositories.
var projTemplateRepo *Repository

// ============================================================================
// Internal functions - template repository for projects
// ============================================================================

// Fetch the project template repository, creating it if necessary.
func getProjTemplateRepository(client *Client) *Repository {
    if projTemplateRepo == nil {
        projTemplateRepo = getTemplateRepository(client, TEMPLATE_PROJ_NAME)
        if projTemplateRepo == nil {
            projTemplateRepo = createProjTemplateRepository(client)
        }
    }
    return projTemplateRepo
}

// Create the project template repository.
func createProjTemplateRepository(client *Client) *Repository {
    name := TEMPLATE_PROJ_NAME
    full := TEMPLATE_PROJ_FULL
    desc := TEMPLATE_PROJ_DESC
    return createTemplateRepository(client, name, full, desc)
}
