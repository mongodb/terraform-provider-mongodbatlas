package streamworkspace

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streaminstance"
)

// AsInstanceModel returns a TFStreamInstanceModel with workspace_name mapped to instance_name
// This eliminates the need for conversion functions by reusing the same underlying data
func (m *TFStreamsWorkspaceModel) AsInstanceModel() *streaminstance.TFStreamInstanceModel {
	return &streaminstance.TFStreamInstanceModel{
		ID:                m.ID,
		InstanceName:      m.WorkspaceName, // Map workspace_name to instance_name
		ProjectID:         m.ProjectID,
		DataProcessRegion: m.DataProcessRegion,
		StreamConfig:      m.StreamConfig,
		Hostnames:         m.Hostnames,
	}
}

// FromInstanceModel populates this workspace model from a TFStreamInstanceModel
// This eliminates the need for conversion functions by directly updating fields
func (m *TFStreamsWorkspaceModel) FromInstanceModel(instanceModel *streaminstance.TFStreamInstanceModel) {
	m.ID = instanceModel.ID
	m.WorkspaceName = instanceModel.InstanceName // Map instance_name to workspace_name
	m.ProjectID = instanceModel.ProjectID
	m.DataProcessRegion = instanceModel.DataProcessRegion
	m.StreamConfig = instanceModel.StreamConfig
	m.Hostnames = instanceModel.Hostnames
}
