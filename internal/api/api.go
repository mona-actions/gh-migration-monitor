package api

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v53/github"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

var query struct {
	Organization struct {
		RepositoryMigrations struct {
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage githubv4.Boolean
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

var migrationStatusQuery struct {
	Organization struct {
		Migration struct {
			Guid                githubv4.String
			Id                  githubv4.ID
			State               githubv4.String
			UploadUrl           githubv4.String
			MigratableResources struct {
				Nodes []struct {
					TargetUrl githubv4.String
					ModelName githubv4.String
				}
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
			} `graphql:"migratableResources(first: $first, after: $after)"`
		} `graphql:"migration(guid: $guid)"`
	} `graphql:"organization(login: $orgName)"`
}

func newGHClient() *githubv4.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: viper.Get("GITHUB_TOKEN").(string)})
	tc := oauth2.NewClient(context.Background(), ts)

	// Keep a reference to the original transport to prevent infinite recursion
	originalTransport := tc.Transport

	if viper.GetBool("ISLEGACY") {
		//log.Println("Using custom header")
		tc.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			// Set the headers
			req.Header.Set("Graphql-Features", "gh_migrator_import_to_dotcom")

			// Use the original RoundTripper to perform the request
			return originalTransport.RoundTrip(req)
		})
	}

	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)

	if err != nil {
		panic(err)
	}

	return githubv4.NewClient(rateLimiter)
}

func newGHRestClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: viper.Get("GITHUB_TOKEN").(string)})
	tc := oauth2.NewClient(ctx, ts)
	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)

	if err != nil {
		panic(err)
	}

	return github.NewClient(rateLimiter)
}

func GetOrgMigrations() []map[string]string {

	if viper.GetBool("ISLEGACY") {
		return GetLegacyMigrations()
	} else {
		return GetGEIMigrations()
	}
}

func GetGEIMigrations() []map[string]string {

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

func GetLegacyMigrations() []map[string]string {

	client := newGHRestClient()
	gqlClient := newGHClient()

	opt := &github.ListOptions{PerPage: 100}

	variables := map[string]interface{}{
		"orgName": githubv4.String(viper.Get("GITHUB_ORGANIZATION").(string)),
		"first":   githubv4.Int(100),
		"after":   (*githubv4.String)(nil),
		"guid":    githubv4.String(""),
	}

	//log.Println("Getting legacy migrations")

	var rm = []map[string]string{}

	for {
		migrations, resp, err := client.Migrations.ListMigrations(context.Background(), viper.Get("GITHUB_ORGANIZATION").(string), opt)

		if err != nil {
			log.Printf("Error getting List of Migrations: %v", err)
			return nil
		}

		for _, migration := range migrations {
			variables["guid"] = githubv4.String(*migration.GUID)
			err := gqlClient.Query(context.Background(), &migrationStatusQuery, variables)
			if err != nil {
				log.Printf("Error executing Migration Status GraphQL query: %v", err)
				return nil
			}
			for _, repoMigration := range migrationStatusQuery.Organization.Migration.MigratableResources.Nodes {
				if repoMigration.ModelName == "repository" {
					rm = append(rm, map[string]string{
						"Id":              *migration.GUID,
						"CreatedAt":       *migration.CreatedAt,
						"FailureReason":   "Unavailable for legacy migrations",
						"RepositoryName":  string(repoMigration.TargetUrl),
						"State":           strings.ToUpper(*migration.State),
						"MigrationLogUrl": *migration.URL,
					})
				}
			}

			if !migrationStatusQuery.Organization.Migration.MigratableResources.PageInfo.HasNextPage {
				continue
			}
			variables["after"] = githubv4.NewString(migrationStatusQuery.Organization.Migration.MigratableResources.PageInfo.EndCursor)
		}
		if resp.NextPage == 0 {
			break
		}
	}
	return rm
}
