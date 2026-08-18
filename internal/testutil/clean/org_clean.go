package clean

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312023/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

// ErrUnauthorized signals a transient HTTP 401 while accessing a project's resources (e.g. a
// project a concurrent run is creating or tearing down). Callers skip the project and retry next run.
var ErrUnauthorized = errors.New("unauthorized accessing project resources, skipping for this run")

// SkipUnauthorizedErr maps a transient 401 to ErrUnauthorized so the caller can skip the project;
// any other error is returned unchanged.
func SkipUnauthorizedErr(resp *http.Response, err error) error {
	if validate.StatusUnauthorized(resp) {
		return ErrUnauthorized
	}
	return err
}

type PendingOrgUser struct {
	InvitedAt time.Time
	ID        string
	Username  string
}

func ListStalePendingOrgUsers(ctx context.Context, client *admin.APIClient, orgID string, prefixes []string, invitedBefore time.Time) ([]PendingOrgUser, error) {
	users, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.OrgUserResponse], *http.Response, error) {
		return client.MongoDBCloudUsersAPI.ListOrgUsers(ctx, orgID).
			OrgMembershipStatus("PENDING").ItemsPerPage(500).PageNum(pageNum).IncludeCount(true).Execute()
	})
	if err != nil {
		return nil, err
	}
	stale := []PendingOrgUser{}
	for i := range users {
		user := &users[i]
		if !hasAnyPrefix(user.GetUsername(), prefixes) {
			continue
		}
		if invitedAt := user.GetInvitationCreatedAt(); !invitedAt.IsZero() && invitedAt.After(invitedBefore) {
			continue
		}
		stale = append(stale, PendingOrgUser{ID: user.GetId(), Username: user.GetUsername(), InvitedAt: user.GetInvitationCreatedAt()})
	}
	return stale, nil
}

func RemovePendingOrgUsers(ctx context.Context, dryRun bool, client *admin.APIClient, orgID string, users []PendingOrgUser) (int, error) {
	if dryRun {
		return len(users), nil
	}
	removed := 0
	for i := range users {
		if _, err := client.MongoDBCloudUsersAPI.RemoveOrgUser(ctx, orgID, users[i].ID).Execute(); err != nil {
			return removed, err
		}
		removed++
	}
	return removed, nil
}

func hasAnyPrefix(name string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// RemoveStreamInstances deletes all stream instances in the project.
// It will also remove all stream processors associated with the stream instance.
func RemoveStreamInstances(ctx context.Context, dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	streamInstances, resp, err := client.StreamsAPI.ListStreamWorkspaces(ctx, projectID).Execute()
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}

	for _, instance := range streamInstances.GetResults() {
		instanceName := *instance.Name

		if !dryRun {
			_, err = client.StreamsAPI.DeleteStreamWorkspace(ctx, projectID, instanceName).Execute()
			if err != nil && admin.IsErrorCode(err, "STREAM_TENANT_HAS_STREAM_PROCESSORS") {
				streamProcessors, spResp, spErr := client.StreamsAPI.GetStreamProcessors(ctx, projectID, instanceName).Execute()
				if spErr != nil {
					return 0, SkipUnauthorizedErr(spResp, spErr)
				}

				processors := streamProcessors.GetResults()
				for i := range processors {
					_, err = client.StreamsAPI.DeleteStreamProcessor(ctx, projectID, instanceName, processors[i].Name).Execute()
					if err != nil {
						return 0, err
					}
				}

				_, err = client.StreamsAPI.DeleteStreamWorkspace(ctx, projectID, instanceName).Execute()
				if err != nil {
					return 0, err
				}
			} else if err != nil {
				return 0, err
			}
		}
	}
	return len(streamInstances.GetResults()), nil
}

// RemovePrivateLinkConnections deletes all Stream Processing Private Link connections in the project;
// left behind, they block project deletion with CANNOT_CLOSE_GROUP_ACTIVE_STREAMS_RESOURCE.
func RemovePrivateLinkConnections(ctx context.Context, dryRun bool, client *admin.APIClient, projectID string) (int, error) {
	connections, resp, err := client.StreamsAPI.ListPrivateLinkConnections(ctx, projectID).Execute()
	if err != nil {
		return 0, SkipUnauthorizedErr(resp, err)
	}
	results := connections.GetResults()
	for i := range results {
		if !dryRun {
			_, err = client.StreamsAPI.DeletePrivateLinkConnection(ctx, projectID, results[i].GetId()).Execute()
			if admin.IsErrorCode(err, "STREAM_PRIVATE_LINK_IN_USE") {
				continue // still referenced by a stream connection, leave it for the next run
			}
			if err != nil {
				return 0, err
			}
		}
	}
	return len(results), nil
}
