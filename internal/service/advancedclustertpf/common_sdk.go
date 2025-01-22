package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/update"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	"go.mongodb.org/atlas-sdk/v20241113004/admin"
)

func createCluster20240805(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin20240805.ClusterDescription20240805, ids *ClusterReader) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV220240805.ClustersApi.CreateCluster(ctx, ids.ProjectID, req).Execute()
	if err != nil {
		diags.AddError(errorCreateLegacy20240805, defaultAPIErrorDetails(ids.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, ids, operationCreateLegacy, diags)
}

func createCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, ids *ClusterReader) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.ClustersApi.CreateCluster(ctx, ids.ProjectID, req).Execute()
	if err != nil {
		diags.AddError(errorCreate, defaultAPIErrorDetails(ids.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, ids, operationCreate, diags)
}

func updateCluster(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, req *admin.ClusterDescription20240805, ids *ClusterReader, operationName string) *admin.ClusterDescription20240805 {
	_, _, err := client.AtlasV2.ClustersApi.UpdateCluster(ctx, ids.ProjectID, ids.ClusterName, req).Execute()
	if err != nil {
		diags.AddError(errorUpdate, defaultAPIErrorDetails(ids.ClusterName, err))
		return nil
	}
	return AwaitChanges(ctx, client, ids, operationName, diags)
}

func updateAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, reqLegacy *admin20240530.ClusterDescriptionProcessArgs, reqNew *admin.ClusterDescriptionProcessArgs20240805, ids *ClusterReader) (legacy *admin20240530.ClusterDescriptionProcessArgs, latest *admin.ClusterDescriptionProcessArgs20240805, changed bool) {
	var (
		err             error
		advConfig       *admin.ClusterDescriptionProcessArgs20240805
		legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs
		projectID       = ids.ProjectID
		clusterName     = ids.ClusterName
	)
	if !update.IsZeroValues(reqLegacy) {
		changed = true
		legacyAdvConfig, _, err = client.AtlasV220240530.ClustersApi.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, reqLegacy).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfUpdateLegacy, defaultAPIErrorDetails(clusterName, err))
			return nil, nil, false
		}
		_ = AwaitChanges(ctx, client, ids, operationAdvancedConfigurationUpdate20240530, diags)
		if diags.HasError() {
			return nil, nil, false
		}
	}
	if !update.IsZeroValues(reqNew) {
		changed = true
		advConfig, _, err = client.AtlasV2.ClustersApi.UpdateClusterAdvancedConfiguration(ctx, projectID, clusterName, reqNew).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfUpdate, defaultAPIErrorDetails(clusterName, err))
			return nil, nil, false
		}
		_ = AwaitChanges(ctx, client, ids, operationAdvancedConfigurationUpdate, diags)
		if diags.HasError() {
			return nil, nil, false
		}
	}
	return legacyAdvConfig, advConfig, changed
}

func readIfUnsetAdvancedConfiguration(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, ids *ClusterReader, configLegacy *admin20240530.ClusterDescriptionProcessArgs, configNew *admin.ClusterDescriptionProcessArgs20240805) (legacy *admin20240530.ClusterDescriptionProcessArgs, latest *admin.ClusterDescriptionProcessArgs20240805) {
	var (
		err         error
		projectID   = ids.ProjectID
		clusterName = ids.ClusterName
	)

	if configLegacy == nil {
		configLegacy, _, err = client.AtlasV220240530.ClustersApi.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfReadLegacy, defaultAPIErrorDetails(clusterName, err))
			return
		}
	}
	if configNew == nil {
		configNew, _, err = client.AtlasV2.ClustersApi.GetClusterAdvancedConfiguration(ctx, projectID, clusterName).Execute()
		if err != nil {
			diags.AddError(errorAdvancedConfRead, defaultAPIErrorDetails(clusterName, err))
			return
		}
	}
	return configLegacy, configNew
}
