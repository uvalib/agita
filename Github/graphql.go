// Github/graphql.go
//
// Access GitHub via the GraphQL API.

package Github

import (
	"context"
	"fmt"
	"log"
	"time"

	"lib.virginia.edu/agita/util"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// ============================================================================
// Internal constants
// ============================================================================

// If *true* end program execution on error.
const gqlABORT = false

// ============================================================================
// Internal variables
// ============================================================================

var gql_client *githubv4.Client

// ============================================================================
// Exported functions
// ============================================================================

// Remove the indicated issue by its unique node ID.
func GqlDeleteIssue(issueNodeId string) bool {
    var Mutation struct {
        DeleteIssue struct {
            ClientMutationID string
        } `graphql:"deleteIssue(input: $input)"`
    }
    input := githubv4.DeleteIssueInput{
        IssueID: githubv4.ID(issueNodeId),
    }
    return gqlMutate(&Mutation, input)
}

// ============================================================================
// Internal functions
// ============================================================================

// Send a GraphQL mutation to GitHub.
func gqlMutate(mut any, input githubv4.Input, variables ...map[string]any) bool {
    ctx := context.Background()
    err := gqlClient().Mutate(ctx, mut, input, gqlVariables(variables...))
    return !gqlError(err)
}

// Send a GraphQL query to GitHub.
func gqlQuery(query any, variables ...map[string]any) bool {
    ctx := context.Background()
    err := gqlClient().Query(ctx, query, gqlVariables(variables...))
    return !gqlError(err)
}

// The current GraphQL client.
func gqlClient() *githubv4.Client {
    if gql_client == nil {
        gql_client = gqlConnect()
    }
    return gql_client
}

// Create an authenticated GraphQL client.
func gqlConnect() *githubv4.Client {
    tok := oauth2.Token{AccessToken: authToken()}
	src := oauth2.StaticTokenSource(&tok)
	cli := oauth2.NewClient(context.Background(), src)
    return githubv4.NewClient(cli)
}

// Merge 0 or more maps.
func gqlVariables(variables ...map[string]any) map[string]any {
    vars := map[string]any{}
    for _, arg := range variables {
        vars = util.MapMerge(vars, arg)
    }
    return vars
}

// Report an error.
//  NOTE: returns false if `err` is nil.
func gqlError(err error) bool {
    if err == nil { return false }
    if gqlABORT { log.Fatalln(err) }
    fmt.Printf("GraphQL error: %v\n", err)
    return true
}

// ============================================================================
// Exported functions
// ============================================================================

// Execute a number of GitHub GraphGL API queries.
func TrialGraphQL(_ ...string) {
    gqlExample1()
    gqlExample2()
    gqlExample3()
    gqlExample4()
    gqlExample5(SAMPLE_ORG, SAMPLE_REPO, SAMPLE_ISSUE)
}

// ============================================================================
// Internal functions
// ============================================================================

func gqlExample1() {
	var Query struct {
		Viewer struct {
			Login     githubv4.String
			CreatedAt time.Time
			AvatarURL string `graphql:"avatarUrl(size: 72)"`
		}
	}
	gqlQuery(&Query)

    tag := "EXAMPLE 1:"
    fmt.Println()
	fmt.Println(tag, Query.Viewer.Login)
	fmt.Println(tag, Query.Viewer.CreatedAt)
	fmt.Println(tag, Query.Viewer.AvatarURL)
}

func gqlExample2() {
	var Query struct {
		Organization struct {
			Login   githubv4.String
			Name    githubv4.String
		}`graphql:"organization(login: \"uvalib\")"`
	}
	gqlQuery(&Query)

    tag := "EXAMPLE 2:"
    fmt.Println()
	fmt.Println(tag, Query.Organization.Login)
	fmt.Println(tag, Query.Organization.Name)
}

func gqlExample3() {
	var Query struct {
		Organization struct {
			Login   githubv4.String
			Name    githubv4.String
		}`graphql:"organization(login: $login)"`
	}
    vars := map[string]any { "login": githubv4.String("uvalib") }
	gqlQuery(&Query, vars)

    tag := "EXAMPLE 3:"
    fmt.Println()
	fmt.Println(tag, Query.Organization.Login)
	fmt.Println(tag, Query.Organization.Name)
}

func gqlExample4() {
    // organization(login: "uvalib") {
    //     repository(name: "emma") {
    //         issues(first: 10) {
    //             nodes {
    //                 id
    //                 title
    //                 bodyHTML
    //             }
    //         }
    //     }
    // }
	var Query struct {
		Organization struct {
			Login   githubv4.String
			Name    githubv4.String
            Repository struct {
                Name githubv4.String
                Issues struct {
                    Nodes []struct {
                        Id    githubv4.String
                        Title githubv4.String
                        Body  githubv4.String
                    }
                } `graphql:"issues(first: 10)"`
            } `graphql:"repository(name: $repo_name)"`
		} `graphql:"organization(login: $org_login)"`
	}
    vars := map[string]any{
        "org_login": githubv4.String(SAMPLE_ORG),
        "repo_name": githubv4.String(SAMPLE_REPO),
    }
	gqlQuery(&Query, vars)

    tag := "EXAMPLE 4:"
    fmt.Println()
	fmt.Println(tag, Query.Organization.Login)
	fmt.Println(tag, Query.Organization.Name)
	fmt.Println(tag, Query.Organization.Repository.Name)
    for idx, iss := range Query.Organization.Repository.Issues.Nodes {
        t := fmt.Sprintf("%s Issue %d:", tag, idx)
        fmt.Println()
        fmt.Println(t, "id    =", iss.Id)
        fmt.Println(t, "title =", iss.Title)
        fmt.Println(t, "body  =", iss.Body)
    }
}

func gqlExample5(org, name string, issueNumber int) {
    // organization(login: org) {
    //     repository(name: name) {
    //         issue(number: issueNumber) {
    //             id
    //             title
    //             bodyHTML
    //         }
    //     }
    // }
	var Query struct {
        Repository struct {
            Name githubv4.String
            Issue struct {
                Id    githubv4.String
                Title githubv4.String
                Body  githubv4.String
            } `graphql:"issue(number: $issue_num)"`
        } `graphql:"repository(owner: $repo_owner, name: $repo_name)"`
	}
    vars := map[string]any{
        "repo_owner": githubv4.String(org),
        "repo_name":  githubv4.String(name),
        "issue_num":  githubv4.Int(issueNumber),
    }
	gqlQuery(&Query, vars)

    tag := "EXAMPLE 5:"
    fmt.Println()
	fmt.Println(tag, Query.Repository.Name)
	fmt.Println(tag, Query.Repository.Issue.Id)
	fmt.Println(tag, Query.Repository.Issue.Title)
	fmt.Println(tag, Query.Repository.Issue.Body)
}
