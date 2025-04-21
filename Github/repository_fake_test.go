// Github/repository_fake_test.go
//
// Generation of temporary repositories to support testing.
//
// Exported functions deal with the temporary fake repositories created during
// execution.
//
// Internal functions deal with the mechanics of creating temporary
// repositories.

package Github

import (
	"lib.virginia.edu/agita/log"
	"lib.virginia.edu/agita/util"

	"github.com/google/go-github/v69/github"
)

// ============================================================================
// Exported constants
// ============================================================================

// Temporary fake repository property.
const (
    FAKE_REPO_NAME = "agita-fake"
    FAKE_REPO_DESC = "temporary fake repository"
)

// Template repository property.
const (
    TEMPLATE_NAME = "agita-test-template"
    TEMPLATE_FULL = "AGITA test template repository"
    TEMPLATE_DESC = "Used as a source for generating repositories in test"
)

// ============================================================================
// Internal variables - fakes
// ============================================================================

var fakeInitialized = !preCleanAll

var fakeRepos []*Repository

// ============================================================================
// Exported functions - fakes
// ============================================================================

// Ensure all GitHub temporary objects are removed.
func InitFakeRepos(client *Client) {
    if !fakeInitialized {
        DeleteAllTemporaryRepos(client)
        fakeInitialized = true
    }
}

// Get fake repositories generated during execution.
func GetFakeRepos(client *Client) []*Repository {
    InitFakeRepos(client)
    if fakeRepos == nil {
        fakeRepos = []*Repository{}
    }
    return fakeRepos
}

// Get a fake repository, creating one if necessary.
func GetFakeRepo(client *Client) *Repository {
    if repos := GetFakeRepos(client); len(repos) == 0 {
        return CreateFakeRepo(client)
    } else {
        return repos[0]
    }
}

// Create a new fake repository.
func CreateFakeRepo(client *Client) *Repository {
    result := createTemporaryRepo(client)
    AddFakeRepo(client, result)
    return result
}

// Insert an entry to `fakeRepos`.
func AddFakeRepo(client *Client, repo *Repository) {
    if repo != nil {
        fakeRepos = append(fakeRepos, repo)
    }
}

// Remove an entry from `fakeRepos`.
func ClearFakeRepo(client *Client, repo *Repository) {
    if repo != nil {
        ClearFakeRepos(client, repo)
    }
}

// Remove entries from `fakeRepos`.
//  NOTE: for use in contexts where the actual GitHub repository has already
//  been deleted.
func ClearFakeRepos(client *Client, repos ...*Repository) {
    if len(repos) == 0 {
        fakeRepos = nil
    } else if fakeRepos != nil {
        for _, repo := range repos {
            for i, fake := range fakeRepos {
                if (fake != nil) && SameRepository(repo, fake) {
                    fakeRepos[i] = nil
                }
            }
        }
        if fakeRepos = util.CompactSlice(fakeRepos); len(fakeRepos) == 0 {
            fakeRepos = nil
        }
    }
}

// Delete a fake repository.
func DeleteFakeRepo(client *Client, repo *Repository) {
    DeleteFakeRepos(client, repo)
}

// Remove (all) fake repositories generated in tests.
func DeleteFakeRepos(client *Client, repos ...*Repository) {
    if len(repos) > 0 {
        deleteTemporaryRepos(client, repos...)
        ClearFakeRepos(client, repos...)
    } else if len(fakeRepos) > 0 {
        deleteTemporaryRepos(client, fakeRepos...)
        ClearFakeRepos(client)
    } else {
        fakeRepos = nil
    }
}

// Remove all temporary test repositories found on GitHub.
func DeleteAllTemporaryRepos(client *Client) {
    deleteTemporaryRepos(client)
    fakeRepos = nil
}

// ============================================================================
// Exported functions - temporary repository
// ============================================================================

// Generate data for a unique temporary repository.
func TemplateRepoData() *github.TemplateRepoRequest {
    return &github.TemplateRepoRequest{
        Name:           github.Ptr(util.Randomize(FAKE_REPO_NAME)),
        Owner:          github.Ptr(ORG),
        Description:    github.Ptr(util.Randomize(FAKE_REPO_DESC)),
        Private:        github.Ptr(true),
    }
}

// Generate data for a unique temporary repository.
func TemplateRepoDataAsRepository() *github.Repository {
    req := TemplateRepoData()
    return &github.Repository{
        Name:           req.Name,
        Owner:          &github.User{Login: req.Owner},
        Description:    req.Description,
        Private:        req.Private,
    }
}

// ============================================================================
// Internal functions - temporary repository
// ============================================================================

// Create a new repository from the test template repository.
func createTemporaryRepo(client *Client) (result *Repository) {
    if templateRepo == nil {
        getTemplateRepository(client)
    }
    srv   := client.ptr.Repositories
    owner := ORG
    name  := TEMPLATE_NAME
    data  := TemplateRepoData()
    repo, _, err := srv.CreateFromTemplate(ctx, owner, name, data)
    if log.ErrorValue(err) == nil {
        result = AsRepositoryType(client, repo)
    }
    return
}

// Get all temporary test repositories from GitHub.
func getTemporaryRepos(client *Client) (result []*Repository) {
    result = []*Repository{}
    res, _, err := client.ptr.Search.Repositories(ctx, FAKE_REPO_NAME, nil)
    if log.ErrorValue(err) == nil {
        for _, repo := range res.Repositories {
            result = append(result, AsRepositoryType(client, repo))
        }
    }
    return result
}

// Remove (all) temporary test repositories found on GitHub.
func deleteTemporaryRepos(client *Client, repos ...*Repository) {
    if len(repos) == 0 {
        repos = getTemporaryRepos(client)
    }
    for _, repo := range repos {
        DeleteRepository(client, repo.Owner, repo.Name)
    }
}

// ============================================================================
// Internal variables - template repository
// ============================================================================

var templateRepo *Repository

// ============================================================================
// Internal functions - template repository
// ============================================================================

// Fetch the test template repository, creating it if necessary.
func getTemplateRepository(client *Client) *Repository {
    if templateRepo == nil {
        log.SuppressPanic()
        defer log.RestorePanic()
        templateRepo = GetRepository(client, ORG, TEMPLATE_NAME)
        if templateRepo == nil {
            log.RestorePanic()
            templateRepo = createTemplateRepository(client)
        }
    }
    return templateRepo
}

// Create the test template repository.
func createTemplateRepository(client *Client) *Repository {
    data := github.Repository{
        Owner:          getUser(client.ptr, ORG),
        Name:           github.Ptr(TEMPLATE_NAME),
        FullName:       github.Ptr(TEMPLATE_FULL),
        Description:    github.Ptr(TEMPLATE_DESC),
        Private:        github.Ptr(true),
        AutoInit:       github.Ptr(true),
        IsTemplate:     github.Ptr(true),
    }
    return CreateRepository(client, &RepositoryRequest{Repository: data})
}
