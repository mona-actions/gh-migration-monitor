package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v53/github"
	"github.com/mona-actions/gh-migration-monitor/internal/models"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// githubClient implements the GitHubClient interface
type githubClient struct {
	restClient    *github.Client
	graphqlClient *githubv4.Client
	rateLimiter   *http.Client
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(token string, isLegacy bool) (GitHubClient, error) {
	if token == "" {
		return nil, fmt.Errorf("github token is required")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	// Keep a reference to the original transport to prevent infinite recursion
	originalTransport := tc.Transport

	// Add custom headers for legacy migrations
	if isLegacy {
		tc.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set("Graphql-Features", "gh_migrator_import_to_dotcom")
			return originalTransport.RoundTrip(req)
		})
	}

	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)
	if err != nil {
		return nil, fmt.Errorf("failed to create rate limiter: %w", err)
	}

	return &githubClient{
		restClient:    github.NewClient(rateLimiter),
		graphqlClient: githubv4.NewClient(rateLimiter),
		rateLimiter:   rateLimiter,
	}, nil
}

// roundTripperFunc allows us to implement RoundTripper as a function
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// ListMigrations implements GitHubClient.ListMigrations
func (c *githubClient) ListMigrations(ctx context.Context, org string, isLegacy bool) ([]models.Migration, error) {
	if isLegacy {
		return c.listLegacyMigrations(ctx, org)
	}
	return c.listGEIMigrations(ctx, org)
}

// listGEIMigrations retrieves migrations using the new GEI API
func (c *githubClient) listGEIMigrations(ctx context.Context, org string) ([]models.Migration, error) {
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

	variables := map[string]interface{}{
		"orgName": githubv4.String(org),
		"first":   githubv4.Int(100),
		"after":   (*githubv4.String)(nil),
	}

	var migrations []models.Migration

	for {
		if err := c.graphqlClient.Query(ctx, &query, variables); err != nil {
			return nil, &APIError{
				StatusCode: 0,
				Message:    fmt.Sprintf("failed to query GEI migrations for org %s", org),
				Err:        err,
			}
		}

		for _, edge := range query.Organization.RepositoryMigrations.Edges {
			createdAt, err := time.Parse(time.RFC3339, edge.Node.CreatedAt)
			if err != nil {
				log.Printf("Failed to parse created_at time %s: %v", edge.Node.CreatedAt, err)
				createdAt = time.Time{}
			}

			migration := models.Migration{
				ID:              edge.Node.Id,
				RepositoryName:  edge.Node.RepositoryName,
				State:           models.State(edge.Node.State),
				CreatedAt:       createdAt,
				FailureReason:   edge.Node.FailureReason,
				MigrationLogURL: edge.Node.MigrationLogUrl,
			}

			migrations = append(migrations, migration)
		}

		if !query.Organization.RepositoryMigrations.PageInfo.HasNextPage {
			break
		}
		variables["after"] = githubv4.NewString(query.Organization.RepositoryMigrations.PageInfo.EndCursor)
	}

	return migrations, nil
}

// listLegacyMigrations retrieves migrations using the legacy migration API
func (c *githubClient) listLegacyMigrations(ctx context.Context, org string) ([]models.Migration, error) {
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

	opt := &github.ListOptions{PerPage: 100}
	variables := map[string]interface{}{
		"orgName": githubv4.String(org),
		"first":   githubv4.Int(100),
		"after":   (*githubv4.String)(nil),
		"guid":    githubv4.String(""),
	}

	var migrations []models.Migration

	for {
		legacyMigrations, resp, err := c.restClient.Migrations.ListMigrations(ctx, org, opt)
		if err != nil {
			return nil, &APIError{
				StatusCode: 0,
				Message:    fmt.Sprintf("failed to list legacy migrations for org %s", org),
				Err:        err,
			}
		}

		for _, migration := range legacyMigrations {
			if migration.GUID == nil {
				continue
			}

			variables["guid"] = githubv4.String(*migration.GUID)
			if err := c.graphqlClient.Query(ctx, &migrationStatusQuery, variables); err != nil {
				log.Printf("Error executing Migration Status GraphQL query: %v", err)
				continue
			}

			for _, resource := range migrationStatusQuery.Organization.Migration.MigratableResources.Nodes {
				if resource.ModelName == "repository" {
					createdAt := time.Time{}
					if migration.CreatedAt != nil {
						if parsed, err := time.Parse(time.RFC3339, *migration.CreatedAt); err == nil {
							createdAt = parsed
						}
					}

					migrationURL := ""
					if migration.URL != nil {
						migrationURL = *migration.URL
					}

					state := ""
					if migration.State != nil {
						state = strings.ToUpper(*migration.State)
					}

					m := models.Migration{
						ID:              *migration.GUID,
						RepositoryName:  string(resource.TargetUrl),
						State:           models.State(state),
						CreatedAt:       createdAt,
						FailureReason:   "Unavailable for legacy migrations",
						MigrationLogURL: migrationURL,
					}

					migrations = append(migrations, m)
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
		opt.Page = resp.NextPage
	}

	return migrations, nil
}
