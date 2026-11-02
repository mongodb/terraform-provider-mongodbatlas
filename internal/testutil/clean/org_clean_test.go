package clean_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312027/admin"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/clean"
)

const (
	itemsPerPage                   = 100
	keepProjectsCreatedWithinHours = 5
	// Resource cleanup for a project can be slow, especially when there are active clusters, that can take more than 10 minutes to delete
	// Once 5 minutes are passed, we give up deleting and hope for the project to be deleted within the next run
	retryInterval = 60 * time.Second
	retryAttempts = 5
)

var (
	defaultProjectPrefixes = []string{
		"cfn-", // general CFN tests
		"ct-",  // CFN contract tests
		"test-acc-tf-p-",
		"atlas-examples-e2e-", // terraform-mongodbatlas-modules/atlas-examples e2e tests
	}
	// keptPrefixes has the prefix of the projects that we want to delete their resources but keep the projects themselves.
	// Useful when a feature flag or cloud provider is configured outside of the test
	keptPrefixes = []string{
		"test-acc-tf-p-keep",
	}
)

func TestSingleProjectRemoval(t *testing.T) {
	projectToClean := os.Getenv("MONGODB_ATLAS_CLEAN_PROJECT_ID")
	if projectToClean == "" {
		t.Skip("skipping test; set MONGODB_ATLAS_CLEAN_PROJECT_ID=project-id to run")
	}
	client := acc.ConnV2()
	dryRun, _ := strconv.ParseBool(os.Getenv("DRY_RUN"))
	changes, err := clean.RemoveProjectResources(t.Context(), t.Logf, dryRun, client, projectToClean)
	require.NoError(t, err)
	if changes != "" {
		t.Logf("project %s %s", projectToClean, changes)
	}
	if !dryRun {
		require.NoError(t, clean.DeleteProject(t.Context(), client, projectToClean))
	}
}

// Using a test to simplify logging and parallelization
func TestCleanProjectAndClusters(t *testing.T) {
	cleanOrg, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_CLEAN_ORG"))
	if !cleanOrg {
		t.Skip("skipping test; set MONGODB_ATLAS_CLEAN_ORG=true to run")
	}
	client := acc.ConnV2()
	dryRun, _ := strconv.ParseBool(os.Getenv("DRY_RUN"))
	onlyZeroClusters, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_CLEAN_ONLY_WHEN_NO_CLUSTERS"))
	skipGracePeriod, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_CLEAN_SKIP_GRACE_PERIOD"))
	skipProjectsAfter := time.Now().Add(-keepProjectsCreatedWithinHours * time.Hour)
	if skipGracePeriod {
		skipProjectsAfter = time.Now()
	}
	retryAttemptsStr := os.Getenv("MONGODB_ATLAS_CLEAN_RETRY_ATTEMPTS")
	runRetries := retryAttempts
	if retryAttemptsStr != "" {
		attempts, err := strconv.Atoi(retryAttemptsStr)
		require.NoError(t, err)
		runRetries = attempts
	}
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	require.NotEmpty(t, orgID, "MONGODB_ATLAS_ORG_ID must be set")
	projects := readAllProjects(t.Context(), t, client, orgID)
	projectsBefore := len(projects)
	t.Logf("found %d projects (DRY_RUN=%t)", projectsBefore, dryRun)
	projectsToClean := map[string]string{}
	projectInfos := []string{}
	prefixes := projectPrefixes()
	for _, p := range projects {
		skipReason := projectSkipReason(&p, skipProjectsAfter, onlyZeroClusters, prefixes)
		projectName := p.GetName()
		projectID := p.GetId()
		if skipReason != "" {
			t.Logf("skip project %s (%s), reason: %s", projectName, projectID, skipReason)
			continue
		}
		projectInfos = append(projectInfos, fmt.Sprintf("Project created at %s name %s (%s)", p.GetCreated().Format(time.RFC3339), projectName, p.GetId()))
		projectsToClean[projectName] = projectID
	}
	t.Logf("deleting project resources and optionally delete projects for %d projects", len(projectsToClean))
	slices.Sort(projectInfos)
	t.Log(strings.Join(projectInfos, "\n"))
	var deleteErrors int
	var emptyProjectCount int
	for name, projectID := range projectsToClean {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			changes, err := clean.RemoveProjectResources(t.Context(), t.Logf, dryRun, client, projectID)
			if errors.Is(err, clean.ErrUnauthorized) {
				t.Logf("skip project %s (%s): unauthorized (transient, will retry next run)", name, projectID)
				return
			}
			require.NoError(t, err)
			if changes != "" {
				t.Logf("project %s %s", name, changes)
			}
			if skipProjectDelete(name) {
				t.Logf("keep project empty, but no delete %s (%s)", name, projectID)
				emptyProjectCount++
				return
			}
			for i := range runRetries {
				attempt := i + 1
				if attempt > 1 {
					time.Sleep(retryInterval)
				}
				t.Logf("attempt %d to delete project %s", attempt, name)
				if dryRun {
					return
				}
				err = clean.DeleteProject(t.Context(), client, projectID)
				if err == nil {
					return
				}
				retryCode := clean.ProjectDeleteRetryCode(err)
				if retryCode != "" {
					t.Logf("attempt %d, project %s has active resources, waiting, error: %s", attempt, projectID, retryCode)
					continue
				}
			}
			t.Logf("failed to delete project %s: %s", name, err)
			deleteErrors++
		})
	}
	t.Cleanup(func() {
		projectsAfter := readAllProjects(context.Background(), t, client, orgID) // reason: using context.Background() here intentionally because t.Context() is canceled at cleanup
		t.Logf("SUMMARY\nProjects changed from %d to %d\ndelete_errors=%d\nempty_project_count=%d\nDRY_RUN=%t", projectsBefore, len(projectsAfter), deleteErrors, emptyProjectCount, dryRun)
	})
}

// Using a test to simplify logging; kept separate from TestCleanProjectAndClusters so org users are cleaned even when project cleanup fails.
func TestCleanOrgUsers(t *testing.T) {
	cleanOrg, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_CLEAN_ORG"))
	if !cleanOrg {
		t.Skip("skipping test; set MONGODB_ATLAS_CLEAN_ORG=true to run")
	}
	dryRun, _ := strconv.ParseBool(os.Getenv("DRY_RUN"))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	require.NotEmpty(t, orgID, "MONGODB_ATLAS_ORG_ID must be set")
	prefixes := userEmailPrefixes()
	t.Logf("cleaning users with email prefixes %v (DRY_RUN=%t)", prefixes, dryRun)
	removedUsers, err := removeTestUsers(t.Context(), t, dryRun, acc.ConnV2(), orgID, prefixes)
	if errors.Is(err, clean.ErrUnauthorized) {
		t.Skipf("skipping test; unauthorized accessing org users (transient, will retry next run): %s", err)
	}
	require.NoError(t, err)
	t.Logf("SUMMARY\nremoved_users=%d\nDRY_RUN=%t", removedUsers, dryRun)
}

// defaultUserEmailPrefixes mirrors the emails built by acc.RandomEmail (`acc.prefixName` + "-").
var defaultUserEmailPrefixes = []string{"test-acc-tf-"}

// removeTestUsers removes every org member whose email starts with one of the given prefixes.
// In dry-run mode it only counts the users it would remove.
func removeTestUsers(ctx context.Context, t *testing.T, dryRun bool, client *admin.APIClient, orgID string, prefixes []string) (int, error) {
	t.Helper()
	users, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.OrgUserResponse], *http.Response, error) {
		orgUsers, resp, err := client.MongoDBCloudUsersAPI.ListOrgUsers(ctx, orgID).ItemsPerPage(itemsPerPage).PageNum(pageNum).Execute()
		if err != nil {
			return nil, resp, clean.SkipUnauthorizedErr(resp, err)
		}
		return orgUsers, resp, nil
	})
	if err != nil {
		return 0, err
	}
	removed := 0
	for i := range users {
		user := users[i]
		if email := user.GetUsername(); !hasAnyPrefix(email, prefixes) {
			continue
		}
		if dryRun {
			t.Logf("would remove test user %s", user.GetUsername())
			removed++
			continue
		}
		resp, err := client.MongoDBCloudUsersAPI.RemoveOrgUser(ctx, orgID, user.GetId()).Execute()
		if err != nil {
			if admin.IsErrorCode(err, "RESOURCE_NOT_FOUND") {
				continue // already removed by a concurrent run
			}
			return removed, clean.SkipUnauthorizedErr(resp, err)
		}
		t.Logf("removed test user %s", user.GetUsername())
		removed++
	}
	return removed, nil
}

func readAllProjects(ctx context.Context, t *testing.T, client *admin.APIClient, orgID string) []admin.Group {
	t.Helper()
	projects, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.Group], *http.Response, error) {
		return client.OrganizationsAPI.GetOrgGroups(ctx, orgID).ItemsPerPage(itemsPerPage).PageNum(pageNum).Execute()
	})
	require.NoError(t, err)
	return projects
}

func projectPrefixes() []string {
	if prefixes := envPrefixes("MONGODB_ATLAS_CLEAN_PROJECT_PREFIXES"); len(prefixes) > 0 {
		return prefixes
	}
	return defaultProjectPrefixes
}

func userEmailPrefixes() []string {
	if prefixes := envPrefixes("MONGODB_ATLAS_CLEAN_USER_EMAIL_PREFIXES"); len(prefixes) > 0 {
		return prefixes
	}
	return defaultUserEmailPrefixes
}

// envPrefixes reads a comma-separated prefix list, dropping blanks. Returns nil when unset or empty.
func envPrefixes(envVar string) []string {
	prefixes := []string{}
	for prefix := range strings.SplitSeq(os.Getenv(envVar), ",") {
		if trimmed := strings.TrimSpace(prefix); trimmed != "" {
			prefixes = append(prefixes, trimmed)
		}
	}
	return prefixes
}

func hasAnyPrefix(value string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func projectSkipReason(p *admin.Group, skipProjectsAfter time.Time, onlyEmpty bool, prefixes []string) string {
	if !hasAnyPrefix(p.GetName(), prefixes) {
		return "name does not match cleanup prefixes"
	}
	if p.GetCreated().After(skipProjectsAfter) {
		return "created after " + skipProjectsAfter.Format("2006-01-02T15:04")
	}
	if onlyEmpty && p.GetClusterCount() > 0 {
		return "has clusters"
	}
	return ""
}

func skipProjectDelete(name string) bool {
	for _, keepPrefix := range keptPrefixes {
		if strings.HasPrefix(name, keepPrefix) {
			return true
		}
	}
	return false
}
