package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var query struct {
	Organization struct {
		RepositoryMigrations struct {
			PageInfo struct {
				StartCursor     githubv4.String
				EndCursor       githubv4.String
				HasNextPage     githubv4.Boolean
				HasPreviousPage githubv4.Boolean
			}
			Nodes []struct {
				Id              string
				CreatedAt       string
				FailureReason   string
				RepositoryName  string
				State           string
				MigrationLogUrl string
			}
		} `graphql:"repositoryMigrations(first: $first, after: $after)"`
	} `graphql:"organization(login: $orgName)"`
}

func newGHClient() *githubv4.Client {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: fmt.Sprint(viper.Get("GITHUB_TOKEN"))})
	httpClient := oauth2.NewClient(context.Background(), src)

	return githubv4.NewClient(httpClient)
}

func GetOrgMigrations() []interface{} {
	client := newGHClient()

	variables := map[string]interface{}{
		"orgName": githubv4.String(viper.Get("GITHUB_ORGANIZATION").(string)),
		"first":   githubv4.Int(100),
		"after":   (*githubv4.String)(nil),
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		panic(err)
	}

	// Collect the query results
	var results []interface{}
	for _, node := range query.Organization.RepositoryMigrations.Nodes {
		out, _ := json.Marshal(node)
		results = append(results, out)
	}

	return results
}
