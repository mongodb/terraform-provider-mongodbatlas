package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const clustersPath = "groups/%s/clusters"

//ClusterService is an interface for interfacing with the Project IP Whitelist
// endpoints of the MongoDB Atlas API.
//See more: https://docs.atlas.mongodb.com/reference/api/whitelist/
type ClusterService interface {
	List(context.Context, string, *ListOptions) ([]Cluster, *Response, error)
	Get(context.Context, string, string) (*Cluster, *Response, error)
	Create(context.Context, string, *Cluster) (*Cluster, *Response, error)
	Update(context.Context, string, string, *Cluster) (*Cluster, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

//ClusterServiceOp handles communication with the Cluster related methods
// of the MongoDB Atlas API
type ClusterServiceOp struct {
	client *Client
}

var _ ClusterService = &ClusterServiceOp{}

type AutoScaling struct {
	DiskGBEnabled *bool `json:"diskGBEnabled,omitempty"`
}

type BiConnector struct {
	Enabled        *bool  `json:"enabled,omitempty"`
	ReadPreference string `json:"readPreference,omitempty"`
}

type ProviderSettings struct {
	BackingProviderName string `json:"backingProviderName,omitempty"`
	DiskIOPS            *int64 `json:"diskIOPS,omitempty"`
	DiskTypeName        string `json:"diskTypeName,omitempty"`
	EncryptEBSVolume    *bool  `json:"encryptEBSVolume,omitempty"`
	InstanceSizeName    string `json:"instanceSizeName,omitempty"`
	ProviderName        string `json:"providerName,omitempty"`
	RegionName          string `json:"regionName,omitempty"`
	VolumeType          string `json:"volumeType,omitempty"`
}

type RegionsConfig struct {
	AnalyticsNodes *int64 `json:"analyticsNodes,omitempty"`
	ElectableNodes *int64 `json:"electableNodes,omitempty"`
	Priority       *int64 `json:"priority,omitempty"`
	ReadOnlyNodes  *int64 `json:"readOnlyNodes,omitempty"`
}

type ReplicationSpec struct {
	ID            string                   `json:"id,omitempty"`
	NumShards     *int64                   `json:"numShards,omitempty"`
	ZoneName      string                   `json:"zoneName,omitempty"`
	RegionsConfig map[string]RegionsConfig `json:"regionsConfig,omitempty"`
}

// Cluster represents MongoDB cluster.
type Cluster struct {
	AutoScaling              AutoScaling              `json:"autoScaling,omitempty"`
	BackupEnabled            *bool                    `json:"backupEnabled,omitempty"`
	BiConnector              BiConnector              `json:"biConnector,omitempty"`
	ClusterType              string                   `json:"clusterType,omitempty"`
	DiskSizeGB               *float64                 `json:"diskSizeGB,omitempty"`
	EncryptionAtRestProvider string                   `json:"encryptionAtRestProvider,omitempty"`
	ID                       string                   `json:"id,omitempty"`
	GroupID                  string                   `json:"groupId,omitempty"`
	MongoDBVersion           string                   `json:"mongoDBVersion,omitempty"`
	MongoDBMajorVersion      string                   `json:"mongoDBMajorVersion,omitempty"`
	MongoURI                 string                   `json:"mongoURI,omitempty"`
	MongoURIUpdated          string                   `json:"mongoURIUpdated,omitempty"`
	MongoURIWithOptions      string                   `json:"mongoURIWithOptions,omitempty"`
	Name                     string                   `json:"name,omitempty"`
	NumShards                *int64                   `json:"numShards"`
	Paused                   *bool                    `json:"paused,omitempty"`
	ProviderBackupEnabled    *bool                    `json:"providerBackupEnabled,omitempty"`
	ProviderSettings         ProviderSettings         `json:"providerSettings,omitempty"`
	ReplicationFactor        *int64                   `json:"replicationFactor,omitempty"`
	ReplicationSpec          map[string]RegionsConfig `json:"replicationSpec,omitempty"`
	ReplicationSpecs         []ReplicationSpec        `json:"replicationSpecs,omitempty"`
	SrvAddress               string                   `json:"srvAddress,omitempty"`
	StateName                string                   `json:"stateName,omitempty"`
}

// clustersResponse is the response from the ClusterService.List.
type clustersResponse struct {
	Links      []*Link   `json:"links,omitempty"`
	Results    []Cluster `json:"results,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}

//List all whitelist entries in the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/whitelist-get-all/
func (s *ClusterServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]Cluster, *Response, error) {
	path := fmt.Sprintf(clustersPath, groupID)

	//Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(clustersResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

//Get gets the whitelist entry specified to {WHITELIST-ENTRY} from the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/whitelist-get-one-entry/
func (s *ClusterServiceOp) Get(ctx context.Context, groupID string, whiteListEntry string) (*Cluster, *Response, error) {
	if whiteListEntry == "" {
		return nil, nil, NewArgError("whiteListEntry", "must be set")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	escapedEntry := url.PathEscape(whiteListEntry)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Add one or more whitelist entries to the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/database-users-create-a-user/
func (s *ClusterServiceOp) Create(ctx context.Context, groupID string, createRequest *Cluster) (*Cluster, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(clustersPath, groupID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Update one or more whitelist entries in the project associated to {GROUP-ID}
//See more: https://docs.atlas.mongodb.com/reference/api/whitelist-update-one/
func (s *ClusterServiceOp) Update(ctx context.Context, groupID string, whitelistEntry string, updateRequest *Cluster) (*Cluster, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, whitelistEntry)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Delete the whitelist entry specified to {WHITELIST-ENTRY} from the project associated to {GROUP-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/whitelist-delete-one/
func (s *ClusterServiceOp) Delete(ctx context.Context, groupID string, whitelistEntry string) (*Response, error) {
	if whitelistEntry == "" {
		return nil, NewArgError("whitelistEntry", "must be set")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	escapedEntry := url.PathEscape(whitelistEntry)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}
