// Github/repository_template.go
//
// Management and use of GitHub template repositories.

package Github

import (
	"lib.virginia.edu/agita/log"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported functions
// ============================================================================

// Generate data for a new repository.
func RepoTemplateData(name, desc string) *github.TemplateRepoRequest {
    return &github.TemplateRepoRequest{
        Name:           github.Ptr(name),
        Owner:          github.Ptr(ORG),
        Description:    github.Ptr(desc),
        Private:        github.Ptr(true),
    }
}

// Generate data for a new repository.
func RepoTemplateDataAsRepository(req *github.TemplateRepoRequest) *github.Repository {
    return &github.Repository{
        Name:           req.Name,
        Owner:          &github.User{Login: req.Owner},
        Description:    req.Description,
        Private:        req.Private,
    }
}

// ============================================================================
// Internal functions
// ============================================================================

// Create a new repository from a template repository.
func createRepoFromTemplate(client *Client, templateName string, data *github.TemplateRepoRequest) (result *Repository) {
    repo, rsp, err := client.ptr.Repositories.CreateFromTemplate(ctx, ORG, templateName, data)
    extractRateLimit(rsp)
    if log.ErrorValue(err) == nil {
        result = AsRepositoryType(client, repo)
    }
    return
}

// ============================================================================
// Internal functions - template repository
// ============================================================================

// Fetch a template repository.
func getTemplateRepository(client *Client, name string) *Repository {
    log.SuppressPanic()
    defer log.RestorePanic()
    return GetRepository(client, ORG, name, true)
}

// Create a template repository.
func createTemplateRepository(client *Client, name, full, desc string) *Repository {
    data := github.Repository{
        Owner:          getUser(client.ptr, ORG),
        Name:           github.Ptr(name),
        FullName:       github.Ptr(full),
        Description:    github.Ptr(desc),
        Private:        github.Ptr(true),
        AutoInit:       github.Ptr(true),
        IsTemplate:     github.Ptr(true),
    }
    return CreateRepository(client, &RepositoryRequest{Repository: data})
}
