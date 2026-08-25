package acc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312023/admin"
)

const (
	CollectionRestoreAPIVersion          = "application/vnd.atlas.2025-03-12+json"
	CollectionRestoreSeedDatabase        = "sample_mflix"
	CollectionRestoreSeedDatabaseRenamed = "sample_mflix_renamed"
	CollectionRestoreSeedRestaurantsNS   = "sample_restaurants.restaurants"
	CollectionRestoreSeedCollectionName  = "movies"
	CollectionRestoreMissingCollectionNS = "sample_mflix.never_there"
	collectionRestoreSnapshotRetention   = 1
)

var collectionRestoreDestRegions = []string{"US_EAST_1", "US_EAST_2", "EU_WEST_1"}

var (
	collectionRestoreFixtureOnce sync.Once
	collectionRestoreFixture     *CollectionRestoreFixture
	collectionRestoreFixtureErr  error
	collectionRestoreDestIdx     atomic.Uint32
)

// CollectionRestoreFixture is the shared source cluster and snapshot for collection-restore acceptance tests.
type CollectionRestoreFixture struct {
	ProjectID               string
	ClusterName             string
	SnapshotID              string
	AfterSnapshotUTCSeconds int64
	External                bool
}

type collectionRestoreEnv struct {
	ProjectID   string
	ClusterName string
	SnapshotID  string
}

func collectionRestoreEnvFromOS() (collectionRestoreEnv, bool) {
	env := collectionRestoreEnv{
		ProjectID:   os.Getenv("MONGODB_ATLAS_PROJECT_ID"),
		ClusterName: os.Getenv("MONGODB_ATLAS_CLUSTER_NAME"),
		SnapshotID:  os.Getenv("MONGODB_ATLAS_SNAPSHOT_ID"),
	}
	ok := env.ProjectID != "" && env.ClusterName != "" && env.SnapshotID != ""
	return env, ok
}

func afterSnapshotUTCSeconds(created time.Time) int64 {
	return created.UTC().Unix()
}

// CollectionRestoreEnvFromOSForTest exposes env-var parsing for unit tests.
func CollectionRestoreEnvFromOSForTest() (projectID, clusterName, snapshotID string, ok bool) {
	env, ok := collectionRestoreEnvFromOS()
	return env.ProjectID, env.ClusterName, env.SnapshotID, ok
}

// AfterSnapshotUTCSecondsForTest exposes snapshot timestamp conversion for unit tests.
func AfterSnapshotUTCSecondsForTest(created time.Time) int64 {
	return afterSnapshotUTCSeconds(created)
}

// CloudBackupCollectionRestoreFixture returns the process-wide source cluster and snapshot.
// When MONGODB_ATLAS_PROJECT_ID, MONGODB_ATLAS_CLUSTER_NAME, and MONGODB_ATLAS_SNAPSHOT_ID are set, those values are reused.
func CloudBackupCollectionRestoreFixture(tb testing.TB) *CollectionRestoreFixture {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "sharedInfo not initialized, use acc.Run() to run tests that require shared resources")
	collectionRestoreFixtureOnce.Do(func() {
		collectionRestoreFixture, collectionRestoreFixtureErr = initCollectionRestoreFixture(tb)
	})
	require.NoError(tb, collectionRestoreFixtureErr)
	return collectionRestoreFixture
}

func initCollectionRestoreFixture(tb testing.TB) (*CollectionRestoreFixture, error) {
	tb.Helper()
	ctx := tb.Context()
	if env, ok := collectionRestoreEnvFromOS(); ok {
		if err := requireCloudBackupAndPIT(ctx, env.ProjectID, env.ClusterName); err != nil {
			return nil, err
		}
		snapshot, err := waitForCompletedSnapshot(ctx, env.ProjectID, env.ClusterName, env.SnapshotID)
		if err != nil {
			return nil, err
		}
		return &CollectionRestoreFixture{
			ProjectID:               env.ProjectID,
			ClusterName:             env.ClusterName,
			SnapshotID:              env.SnapshotID,
			AfterSnapshotUTCSeconds: afterSnapshotUTCSeconds(snapshot.GetCreatedAt()),
			External:                true,
		}, nil
	}

	projectID, clusterName := clusterNameExecution(tb, true, true, true)
	if err := requireCloudBackupAndPIT(ctx, projectID, clusterName); err != nil {
		return nil, err
	}
	snapshotID, afterSeconds, err := takeCompletedOnDemandSnapshot(ctx, projectID, clusterName)
	if err != nil {
		return nil, err
	}
	return &CollectionRestoreFixture{
		ProjectID:               projectID,
		ClusterName:             clusterName,
		SnapshotID:              snapshotID,
		AfterSnapshotUTCSeconds: afterSeconds,
	}, nil
}

func requireCloudBackupAndPIT(ctx context.Context, projectID, clusterName string) error {
	src, _, err := ConnV2().ClustersAPI.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		return fmt.Errorf("get cluster %s: %w", clusterName, err)
	}
	return errIfCloudBackupOrPITDisabled(clusterName, src.GetBackupEnabled(), src.GetPitEnabled())
}

func errIfCloudBackupOrPITDisabled(clusterName string, backupEnabled, pitEnabled bool) error {
	if !backupEnabled {
		return fmt.Errorf("cluster %s does not have cloud backup enabled", clusterName)
	}
	if !pitEnabled {
		return fmt.Errorf("cluster %s does not have point-in-time restore enabled", clusterName)
	}
	return nil
}

// ErrIfCloudBackupOrPITDisabledForTest exposes backup/PIT flag checking for unit tests.
func ErrIfCloudBackupOrPITDisabledForTest(clusterName string, backupEnabled, pitEnabled bool) error {
	return errIfCloudBackupOrPITDisabled(clusterName, backupEnabled, pitEnabled)
}

// NextCollectionRestoreDestRegion returns a distinct AWS region for each dest cluster in this process.
func NextCollectionRestoreDestRegion() string {
	i := collectionRestoreDestIdx.Add(1) - 1
	return collectionRestoreDestRegions[int(i)%len(collectionRestoreDestRegions)]
}

func takeCompletedOnDemandSnapshot(ctx context.Context, projectID, clusterName string) (snapshotID string, afterSeconds int64, err error) {
	retention := collectionRestoreSnapshotRetention
	params := &admin.DiskBackupOnDemandSnapshotRequest{
		Description:     new("tf acc collection restore fixture"),
		RetentionInDays: &retention,
	}
	snapshot, _, err := ConnV2().CloudBackupsAPI.TakeSnapshots(ctx, projectID, clusterName, params).Execute()
	if err != nil {
		return "", 0, fmt.Errorf("take snapshot: %w", err)
	}
	snapshotID = snapshot.GetId()
	if _, err := waitForCompletedSnapshot(ctx, projectID, clusterName, snapshotID); err != nil {
		return "", 0, err
	}
	return snapshotID, time.Now().UTC().Unix(), nil
}

func waitForCompletedSnapshot(ctx context.Context, projectID, clusterName, snapshotID string) (*admin.DiskBackupReplicaSet, error) {
	requestParams := &admin.GetClusterBackupSnapshotApiParams{
		GroupId:     projectID,
		ClusterName: clusterName,
		SnapshotId:  snapshotID,
	}
	stateConf := retry.StateChangeConf{
		Pending:    []string{"queued", "inProgress"},
		Target:     []string{"completed"},
		Timeout:    2 * time.Hour,
		MinTimeout: 1 * time.Minute,
		Delay:      1 * time.Minute,
		Refresh: func() (any, string, error) {
			cur, _, getErr := ConnV2().CloudBackupsAPI.GetClusterBackupSnapshotWithParams(ctx, requestParams).Execute()
			if getErr != nil {
				return nil, "", getErr
			}
			status := cur.GetStatus()
			if status == "failed" {
				return nil, status, fmt.Errorf("snapshot %s failed", snapshotID)
			}
			return cur, status, nil
		},
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return nil, fmt.Errorf("wait snapshot %s: %w", snapshotID, err)
	}
	snapshot, _, err := ConnV2().CloudBackupsAPI.GetClusterBackupSnapshot(ctx, projectID, clusterName, snapshotID).Execute()
	if err != nil {
		return nil, fmt.Errorf("get snapshot %s: %w", snapshotID, err)
	}
	return snapshot, nil
}

type collectionRestoreJobView struct {
	ID                string `json:"id"`
	State             string `json:"state"`
	TargetClusterName string `json:"targetClusterName"`
}

type collectionRestoreJobList struct {
	Results []collectionRestoreJobView `json:"results"`
}

func collectionRestoreJobCall(projectID, clusterName, jobID string) config.APICallParams {
	pathParams := map[string]string{
		"projectId":   projectID,
		"clusterName": clusterName,
	}
	relativePath := "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/collectionRestoreJobs"
	if jobID != "" {
		pathParams["jobId"] = jobID
		relativePath += "/{jobId}"
	}
	return config.APICallParams{
		VersionHeader: CollectionRestoreAPIVersion,
		RelativePath:  relativePath,
		PathParams:    pathParams,
		Method:        http.MethodGet,
	}
}

func closeOnAPIErr(resp *http.Response, err error) error {
	if err != nil && resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	return err
}

func readJSONBody(resp *http.Response, dest any) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dest)
}

// GetCollectionRestoreJob GETs one job via the untyped client and the resource version header.
func GetCollectionRestoreJob(ctx context.Context, projectID, clusterName, jobID string) (map[string]any, error) {
	resp, err := MongoDBClient.UntypedAPICall(ctx, collectionRestoreJobCall(projectID, clusterName, jobID), nil)
	if err := closeOnAPIErr(resp, err); err != nil {
		return nil, err
	}
	var body map[string]any
	if err := readJSONBody(resp, &body); err != nil {
		return nil, err
	}
	return body, nil
}

// LatestCollectionRestoreJobID returns the newest job on the source cluster whose target matches destClusterName.
func LatestCollectionRestoreJobID(ctx context.Context, projectID, clusterName, destClusterName string) (string, error) {
	resp, err := MongoDBClient.UntypedAPICall(ctx, collectionRestoreJobCall(projectID, clusterName, ""), nil)
	if err := closeOnAPIErr(resp, err); err != nil {
		return "", err
	}
	var list collectionRestoreJobList
	if err := readJSONBody(resp, &list); err != nil {
		return "", err
	}
	fallback := ""
	for _, job := range list.Results {
		if job.TargetClusterName != destClusterName || job.ID == "" {
			continue
		}
		switch job.State {
		case "INITIALIZING", "IN_PROGRESS", "FINALIZING":
			return job.ID, nil
		default:
			fallback = job.ID
		}
	}
	if fallback != "" {
		return fallback, nil
	}
	return "", fmt.Errorf("no collection restore job targeting %s", destClusterName)
}

// WaitCollectionRestoreJobState polls until job state matches wantState.
func WaitCollectionRestoreJobState(ctx context.Context, projectID, clusterName, jobID, wantState string) error {
	stateConf := retry.StateChangeConf{
		Pending:    []string{"INITIALIZING", "IN_PROGRESS", "FINALIZING"},
		Target:     []string{wantState},
		Timeout:    constant.DefaultTimeout,
		MinTimeout: 1 * time.Minute,
		Delay:      30 * time.Second,
		Refresh: func() (any, string, error) {
			body, err := GetCollectionRestoreJob(ctx, projectID, clusterName, jobID)
			if err != nil {
				return nil, "", err
			}
			state, _ := body["state"].(string)
			if (state == "FAILED" || state == "CANCELED") && state != wantState {
				return body, state, fmt.Errorf("collection restore job %s ended with state %s", jobID, state)
			}
			return body, state, nil
		},
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
