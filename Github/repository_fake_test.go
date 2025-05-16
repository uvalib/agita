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

// Template fake repository property.
const (
    TEMPLATE_FAKE_NAME = "agita-test-template"
    TEMPLATE_FAKE_FULL = "AGITA test template repository"
    TEMPLATE_FAKE_DESC = "Used as a source for generating repositories in test"
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
    result := createFakeRepository(client)
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
func FakeRepoTemplateData() *github.TemplateRepoRequest {
    name := util.Randomize(FAKE_REPO_NAME)
    desc := util.Randomize(FAKE_REPO_DESC)
    return RepoTemplateData(name, desc)
}

// Generate data for a unique temporary repository.
func FakeRepoTemplateDataAsRepository() *github.Repository {
    req := FakeRepoTemplateData()
    return RepoTemplateDataAsRepository(req)
}

// ============================================================================
// Internal functions - temporary repository
// ============================================================================

// Create a new repository from the test template repository.
func createFakeRepository(client *Client) *Repository {
    getFakeTemplateRepository(client)
    data := FakeRepoTemplateData()
    return createRepoFromTemplate(client, TEMPLATE_FAKE_NAME, data)
}

// Get all temporary test repositories from GitHub.
func getTemporaryRepos(client *Client) []*Repository {
    result := []*Repository{}
    res, rsp, err := client.ptr.Search.Repositories(ctx, FAKE_REPO_NAME, nil)
    extractRateLimit(rsp)
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
// Internal variables - template repository for fakes
// ============================================================================

// GitHub repository used as a template for generating fake repositories.
var fakeTemplateRepo *Repository

// ============================================================================
// Internal functions - template repository for fakes
// ============================================================================

// Fetch the test template repository, creating it if necessary.
func getFakeTemplateRepository(client *Client) *Repository {
    if fakeTemplateRepo == nil {
        fakeTemplateRepo = getTemplateRepository(client, TEMPLATE_FAKE_NAME)
        if fakeTemplateRepo == nil {
            fakeTemplateRepo = createFakeTemplateRepository(client)
        }
    }
    return fakeTemplateRepo
}

// Create the test template repository.
func createFakeTemplateRepository(client *Client) *Repository {
    name := TEMPLATE_FAKE_NAME
    full := TEMPLATE_FAKE_FULL
    desc := TEMPLATE_FAKE_DESC
    return createTemplateRepository(client, name, full, desc)
}
