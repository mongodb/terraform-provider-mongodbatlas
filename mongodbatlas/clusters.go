package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const clustersPath = "groups/%s/clusters"

//ClustersService is an interface for interfacing with the Clusters
// endpoints of the MongoDB Atlas API.
//See more: https://docs.atlas.mongodb.com/reference/api/clusters/
type ClustersService interface {
	List(context.Context, string, *ListOptions) ([]Cluster, *Response, error)
	Get(context.Context, string, string) (*Cluster, *Response, error)
	Create(context.Context, string, *Cluster) (*Cluster, *Response, error)
	Update(context.Context, string, string, *Cluster) (*Cluster, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	UpdateProcessArgs(context.Context, string, string, *ProcessArgs) (*ProcessArgs, *Response, error)
	GetProcessArgs(context.Context, string, string) (*ProcessArgs, *Response, error)
}

//ClustersServiceOp handles communication with the Cluster related methods
// of the MongoDB Atlas API
type ClustersServiceOp struct {
	client *Client
}

var _ ClustersService = &ClustersServiceOp{}

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
	NumShards                *int64                   `json:"numShards,omitempty"`
	Paused                   *bool                    `json:"paused,omitempty"`
	ProviderBackupEnabled    *bool                    `json:"providerBackupEnabled,omitempty"`
	ProviderSettings         *ProviderSettings        `json:"providerSettings,omitempty"`
	ReplicationFactor        *int64                   `json:"replicationFactor,omitempty"`
	ReplicationSpec          map[string]RegionsConfig `json:"replicationSpec,omitempty"`
	ReplicationSpecs         []ReplicationSpec        `json:"replicationSpecs,omitempty"`
	SrvAddress               string                   `json:"srvAddress,omitempty"`
	StateName                string                   `json:"stateName,omitempty"`
}

type ProcessArgs struct {
	FailIndexKeyTooLong              *bool  `json:"failIndexKeyTooLong"`
	JavascriptEnabled                *bool  `json:"javascriptEnabled"`
	MinimumEnabledTLSProtocol        string `json:"minimumEnabledTlsProtocol"`
	NoTableScan                      *bool  `json:"noTableScan"`
	OplogSizeMB                      *int64 `json:"oplogSizeMB"`
	SampleSizeBIConnector            *int64 `json:"sampleSizeBIConnector"`
	SampleRefreshIntervalBIConnector *int64 `json:"sampleRefreshIntervalBIConnector"`
}

// clustersResponse is the response from the ClustersService.List.
type clustersResponse struct {
	Links      []*Link   `json:"links,omitempty"`
	Results    []Cluster `json:"results,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}

//List all clusters in the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/clusters-get-all/
func (s *ClustersServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]Cluster, *Response, error) {
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

//Get gets the cluster specified to {ClUSTER-NAME} from the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/clusters-get-one/
func (s *ClustersServiceOp) Get(ctx context.Context, groupID string, clusterName string) (*Cluster, *Response, error) {
	if clusterName == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
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

//Add a cluster to the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/clusters-create-one/
func (s *ClustersServiceOp) Create(ctx context.Context, groupID string, createRequest *Cluster) (*Cluster, *Response, error) {
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

//Update a cluster in the project associated to {GROUP-ID}
//See more: https://docs.atlas.mongodb.com/reference/api/clusters-modify-one/
func (s *ClustersServiceOp) Update(ctx context.Context, groupID string, clusterName string, updateRequest *Cluster) (*Cluster, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, clusterName)

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

//Delete the cluster specified to {CLUSTER-NAME} from the project associated to {GROUP-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/clusters-delete-one/
func (s *ClustersServiceOp) Delete(ctx context.Context, groupID string, clusterName string) (*Response, error) {
	if clusterName == "" {
		return nil, NewArgError("clusterName", "must be set")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

//UpdateProcessArgs Modifies Advanced Configuration Options for One Cluster
//See more: https://docs.atlas.mongodb.com/reference/api/clusters-modify-advanced-configuration-options/
func (s *ClustersServiceOp) UpdateProcessArgs(ctx context.Context, groupID string, clusterName string, updateRequest *ProcessArgs) (*ProcessArgs, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	path := fmt.Sprintf("%s/%s/processArgs", basePath, clusterName)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessArgs)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//GetProcessArgs gets the Advanced Configuration Options for One Cluster
//See more: https://docs.atlas.mongodb.com/reference/api/clusters-get-advanced-configuration-options/#get-advanced-configuration-options-for-one-cluster
func (s *ClustersServiceOp) GetProcessArgs(ctx context.Context, groupID string, clusterName string) (*ProcessArgs, *Response, error) {
	if clusterName == "" {
		return nil, nil, NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(clustersPath, groupID)
	escapedEntry := url.PathEscape(clusterName)
	path := fmt.Sprintf("%s/%s/processArgs", basePath, escapedEntry)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessArgs)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
