package api

import (
	"context"
	"fmt"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
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
			Edges []struct {
				Node struct {
					Id              string
					CreatedAt       string
					FailureReason   string
					RepositoryName  string
					State           string
					MigrationLogUrl string
				}
			}
		} `graphql:"repositoryMigrations(first: $first, after: $after)"`
	} `graphql:"organization(login: $orgName)"`
}

func newGHClient() *githubv4.Client {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: fmt.Sprint(viper.Get("GITHUB_TOKEN"))})
	httpClient := oauth2.NewClient(context.Background(), src)
	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(httpClient.Transport)

	if err != nil {
		panic(err)
	}

	return githubv4.NewClient(rateLimiter)
}

func GetOrgMigrations() []map[string]string {
	client := newGHClient()

	variables := map[string]interface{}{
		"orgName": githubv4.String(viper.Get("GITHUB_ORGANIZATION").(string)),
		"first":   githubv4.Int(100),
		"after":   (*githubv4.String)(nil),
	}

	var rm = []map[string]string{}
	for {
		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			panic(err)
		}

		for _, repoMigration := range query.Organization.RepositoryMigrations.Edges {
			rm = append(rm, map[string]string{
				"Id":              repoMigration.Node.Id,
				"CreatedAt":       repoMigration.Node.CreatedAt,
				"FailureReason":   repoMigration.Node.FailureReason,
				"RepositoryName":  repoMigration.Node.RepositoryName,
				"State":           repoMigration.Node.State,
				"MigrationLogUrl": repoMigration.Node.MigrationLogUrl})
		}

		if !query.Organization.RepositoryMigrations.PageInfo.HasNextPage {
			break
		}
		variables["after"] = githubv4.NewString(query.Organization.RepositoryMigrations.PageInfo.EndCursor)
	}

	return rm
}
