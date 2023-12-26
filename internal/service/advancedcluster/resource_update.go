package advancedcluster

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mwielbut/pointy"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/utility"
)

func (r *advancedClusterRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	conn := r.Client.Atlas
	var state, plan, tfconfig tfAdvancedClusterRSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &tfconfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := conversion.DecodeStateID(state.ID.ValueString())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	timeout, _ := plan.Timeouts.Update(ctx, defaultTimeout)

	if upgradeRequest := getUpgradeRequest(ctx, &state, &plan); upgradeRequest != nil {
		_, _, err := upgradeCluster(ctx, conn, upgradeRequest, projectID, clusterName, timeout)

		if err != nil {
			resp.Diagnostics.AddError("Unable to UPDATE cluster. An error occurred while upgrading cluster.", err.Error())
			return
		}
	} else {
		resp.Diagnostics.Append(updateCluster(ctx, conn, &state, &plan, timeout)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	cluster, response, err := conn.AdvancedClusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error during cluster READ from Atlas", fmt.Sprintf(errorClusterAdvancedRead, clusterName, err.Error()))
		return
	}

	log.Printf("[DEBUG] GET ClusterAdvanced %+v", cluster)
	newClusterModel, diags := newTfAdvancedClusterRSModel(ctx, conn, cluster, &plan)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterModel)...)
}

func getUpgradeRequest(ctx context.Context, state, plan *tfAdvancedClusterRSModel) *matlas.Cluster {
	if reflect.DeepEqual(plan.ReplicationSpecs, state.ReplicationSpecs) {
		return nil
	}

	currentSpecs := newReplicationSpecs(ctx, state.ReplicationSpecs)
	updatedSpecs := newReplicationSpecs(ctx, plan.ReplicationSpecs)

	if len(currentSpecs) != 1 || len(updatedSpecs) != 1 || len(currentSpecs[0].RegionConfigs) != 1 || len(updatedSpecs[0].RegionConfigs) != 1 {
		return nil
	}

	currentRegion := currentSpecs[0].RegionConfigs[0]
	updatedRegion := updatedSpecs[0].RegionConfigs[0]
	currentSize := currentRegion.ElectableSpecs.InstanceSize

	if currentRegion.ElectableSpecs.InstanceSize == updatedRegion.ElectableSpecs.InstanceSize || !IsSharedTier(currentSize) {
		return nil
	}

	return &matlas.Cluster{
		ProviderSettings: &matlas.ProviderSettings{
			ProviderName:     updatedRegion.ProviderName,
			InstanceSizeName: updatedRegion.ElectableSpecs.InstanceSize,
			RegionName:       updatedRegion.RegionName,
		},
	}
}

func upgradeCluster(ctx context.Context, conn *matlas.Client, request *matlas.Cluster, projectID, name string, timeout time.Duration) (*matlas.Cluster, *matlas.Response, error) {
	request.Name = name

	cluster, resp, err := conn.Clusters.Upgrade(ctx, projectID, request)
	if err != nil {
		return nil, nil, err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
		Target:     []string{"IDLE"},
		Refresh:    ResourceClusterRefreshFunc(ctx, name, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	return cluster, resp, nil
}

func updateCluster(ctx context.Context, conn *matlas.Client, state, plan *tfAdvancedClusterRSModel, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	ids := conversion.DecodeStateID(state.ID.ValueString())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	cluster := new(matlas.AdvancedCluster)
	clusterChangeDetect := new(matlas.AdvancedCluster)

	if hasBoolUpdated(plan.BackupEnabled, state.BackupEnabled) {
		cluster.BackupEnabled = plan.BackupEnabled.ValueBoolPointer()
	}
	if hasStringUpdated(plan.ClusterType, state.ClusterType) {
		cluster.ClusterType = plan.ClusterType.ValueString()
	}
	if hasFloatUpdated(plan.DiskSizeGb, state.DiskSizeGb) {
		cluster.DiskSizeGB = plan.DiskSizeGb.ValueFloat64Pointer()
	}
	if hasStringUpdated(plan.EncryptionAtRestProvider, state.EncryptionAtRestProvider) {
		cluster.EncryptionAtRestProvider = plan.EncryptionAtRestProvider.ValueString()
	}

	// TODO: remove db version logic
	if !plan.MongoDBMajorVersion.Equal(state.MongoDBMajorVersion) {
		cluster.MongoDBMajorVersion = utility.FormatMongoDBMajorVersion(plan.MongoDBMajorVersion.ValueString())
	}

	if hasBoolUpdated(plan.PitEnabled, state.PitEnabled) {
		cluster.PitEnabled = plan.PitEnabled.ValueBoolPointer()
	}

	if hasStringUpdated(plan.RootCertType, state.RootCertType) {
		cluster.RootCertType = plan.RootCertType.ValueString()
	}
	if hasBoolUpdated(plan.TerminationProtectionEnabled, state.TerminationProtectionEnabled) {
		cluster.TerminationProtectionEnabled = plan.TerminationProtectionEnabled.ValueBoolPointer()
	}
	if hasOptionalStringUpdated(plan.AcceptDataRisksAndForceReplicaSetReconfig, state.AcceptDataRisksAndForceReplicaSetReconfig) {
		cluster.AcceptDataRisksAndForceReplicaSetReconfig = plan.AcceptDataRisksAndForceReplicaSetReconfig.ValueString()
	}
	if hasBoolUpdated(plan.Paused, state.Paused) && !plan.Paused.ValueBool() {
		cluster.Paused = plan.Paused.ValueBoolPointer()
	}
	// if isUpdatedBool(plan.RetainBackupsEnabled, state.RetainBackupsEnabled) {
	// 	cluster.RetainBackupsEnabled = plan.RetainBackupsEnabled.ValueBoolPointer()
	// }

	if updated, newPlan, d := biConnectorConfigIfUpdated(ctx, plan.BiConnectorConfig, state.BiConnectorConfig); !d.HasError() && updated {
		cluster.BiConnector = newPlan
	}

	// Labels is optional so state/plan will either be null or known
	if !reflect.DeepEqual(plan.Labels, state.Labels) {
		if ContainsLabelOrKey(newLabels(ctx, plan.Labels), DefaultLabel) {
			diags.AddError("Unable to UPDATE cluster. An error occurred when updating labels.", "you should not set `Infrastructure Tool` label, it is used for internal purposes")
			return diags
		}
		cluster.Labels = newLabels(ctx, plan.Labels)
	}

	// Tags is optional so state/plan will either be null or known
	if !reflect.DeepEqual(plan.Tags, state.Tags) {
		cluster.Tags = newTags(ctx, plan.Tags)
	}

	if updated, newPlan, d := replicationSpecsIfUpdated(ctx, plan.ReplicationSpecs, state.ReplicationSpecs); !d.HasError() && updated {
		cluster.ReplicationSpecs = newPlan
	}

	// var tfRepSpecsPlan, tfRepSpecsState []tfReplicationSpecRSModel
	// if isReplicationSpecUpdated(plan.ReplicationSpecs, state.ReplicationSpecs) {
	// 	// TODO remove:
	// 	plan.ReplicationSpecs.ElementsAs(ctx, &tfRepSpecsPlan, true)
	// 	state.ReplicationSpecs.ElementsAs(ctx, &tfRepSpecsState, true)

	// 	if !reflect.DeepEqual(tfRepSpecsPlan, tfRepSpecsState) {
	// 		if !reflect.DeepEqual(tfRepSpecsPlan[0].RegionsConfigs, tfRepSpecsState[0].RegionsConfigs) {
	// 			var tfRegionConfigsPlan, tfRegionConfigsState []tfRegionsConfigModel
	// 			tfRepSpecsPlan[0].RegionsConfigs.ElementsAs(ctx, &tfRegionConfigsPlan, true)
	// 			tfRepSpecsState[0].RegionsConfigs.ElementsAs(ctx, &tfRegionConfigsState, true)
	// 		}
	// 	}
	// 	cluster.ReplicationSpecs = newReplicationSpecs(ctx, plan.ReplicationSpecs)
	// }

	// TODO add comment in PR:
	// This logic has been updated from SDKv2 implementation where if the user removes advanced_confgiuration block
	// we would not call Update API, instead, now we send empty request object with all null values so API can reset defaults wherever applicable
	// and return those.
	if updated, newPlan, d := advancedConfigIfUpdated(ctx, plan.AdvancedConfiguration, state.AdvancedConfiguration); !d.HasError() && updated {
		// advancedConfReq := newAdvancedConfiguration(ctx, ac)
		// if !reflect.DeepEqual(newPlan, matlas.ProcessArgs{}) { // TODO check if required, see comment inside update func
		_, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, clusterName, newPlan)
		if err != nil {
			diags.AddError("Unable to UPDATE cluster. An error occurred when updating advanced_configuration.", err.Error())
			return diags
		}
		// }
	}

	// Has changes
	if !reflect.DeepEqual(cluster, clusterChangeDetect) {
		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, resp, err := updateAdvancedCluster(ctx, conn, cluster, projectID, clusterName, timeout)
			if err != nil {
				if resp == nil || resp.StatusCode == 400 {
					return retry.NonRetryableError(fmt.Errorf(errorClusterAdvancedUpdate, clusterName, err))
				}
				return retry.RetryableError(fmt.Errorf(errorClusterAdvancedUpdate, clusterName, err))
			}
			return nil
		})
		if err != nil {
			diags.AddError("Unable to UPDATE cluster. An error occurred when updating cluster in Atlas.", err.Error())
			return diags
		}
	}

	if plan.Paused.ValueBool() {
		clusterRequest := &matlas.AdvancedCluster{
			Paused: pointy.Bool(true),
		}

		_, _, err := updateAdvancedCluster(ctx, conn, clusterRequest, projectID, clusterName, timeout)
		if err != nil {
			diags.AddError("Unable to UPDATE cluster. An error occurred when attempting to pause cluster in Atlas.", err.Error())
			return diags
		}
	}

	return diags
}

func updateAdvancedCluster(ctx context.Context, conn *matlas.Client, request *matlas.AdvancedCluster, projectID, name string, timeout time.Duration,
) (*matlas.AdvancedCluster, *matlas.Response, error) {
	cluster, resp, err := conn.AdvancedClusters.Update(ctx, projectID, name, request)
	if err != nil {
		return nil, nil, err
	}

	if err := waitClusterUpdate(ctx, conn, timeout, projectID, cluster.Name); err != nil {
		return nil, nil, err
	}

	return cluster, resp, nil
}

func waitClusterUpdate(ctx context.Context, conn *matlas.Client, timeout time.Duration, projectID, clusterName string) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterAdvancedRefreshFunc(ctx, clusterName, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func biConnectorConfigIfUpdated(ctx context.Context, planVal, stateVal types.List) (bool, *matlas.BiConnector, diag.Diagnostics) {
	var d diag.Diagnostics

	var planConfig, stateConfig []TfBiConnectorConfigModel
	planVal.ElementsAs(ctx, &planConfig, true)
	stateVal.ElementsAs(ctx, &stateConfig, true)

	if !planVal.IsUnknown() && !stateVal.IsUnknown() {
		if updated, d := hasBiConnectorConfigUpdated(planConfig, stateConfig); d.HasError() || !updated {
			return false, nil, d
		}

		// if certain attributes in the plan are unknown, we fetch those from the state to create complete object to send to the request
		// while only using plan values for the updated/user configured attributes
		if planConfig[0].Enabled.IsUnknown() {
			planConfig[0].Enabled = stateConfig[0].Enabled
		}
		if planConfig[0].ReadPreference.IsUnknown() {
			planConfig[0].ReadPreference = stateConfig[0].ReadPreference
		}

		updatedPlan, d := types.ListValueFrom(ctx, TfBiConnectorConfigType, planConfig)
		if d.HasError() {
			return false, nil, d
		}

		return true, newBiConnectorConfig(ctx, updatedPlan), d
	}

	return false, nil, d
}

func hasBiConnectorConfigUpdated(planConfig, stateConfig []TfBiConnectorConfigModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(planConfig) < 1 || len(stateConfig) < 1 {
		return false, diags
	}

	pc := planConfig[0]
	sc := stateConfig[0]

	if !pc.Enabled.IsUnknown() && !pc.Enabled.Equal(sc.Enabled) {
		return true, diags
	}
	if !pc.ReadPreference.IsUnknown() && !pc.ReadPreference.Equal(sc.ReadPreference) {
		return true, diags
	}

	return false, diags
}

func advancedConfigIfUpdated(ctx context.Context, planVal, stateVal types.List) (bool, *matlas.ProcessArgs, diag.Diagnostics) {
	var d diag.Diagnostics

	var planConfig, stateConfig []TfAdvancedConfigurationModel
	planVal.ElementsAs(ctx, &planConfig, true)
	stateVal.ElementsAs(ctx, &stateConfig, true)

	// updated, d := hasAdvancedConfigUpdated(planConfig, stateConfig)
	// if d.HasError() || !updated {
	// 	return false, nil, d
	// }

	// if removed
	// if len(planConfig) == 0 && len(stateConfig) == 1 {
	// 	return true, &matlas.ProcessArgs{}, d
	// }

	if !planVal.IsUnknown() && !stateVal.IsUnknown() {
		if updated, d := hasAdvancedConfigUpdated(planConfig, stateConfig); d.HasError() || !updated {
			return false, nil, d
		}
		// if certain attributes in the plan are unknown, we fetch those from the state to create complete object to send to the request
		// while only using plan values for the updated/user configured attributes
		if planConfig[0].DefaultReadConcern.IsUnknown() {
			planConfig[0].DefaultReadConcern = stateConfig[0].DefaultReadConcern
		}
		if planConfig[0].DefaultWriteConcern.IsUnknown() {
			planConfig[0].DefaultWriteConcern = stateConfig[0].DefaultWriteConcern
		}
		if planConfig[0].FailIndexKeyTooLong.IsUnknown() {
			planConfig[0].FailIndexKeyTooLong = stateConfig[0].FailIndexKeyTooLong
		}
		if planConfig[0].JavascriptEnabled.IsUnknown() {
			planConfig[0].JavascriptEnabled = stateConfig[0].JavascriptEnabled
		}
		if planConfig[0].MinimumEnabledTLSProtocol.IsUnknown() {
			planConfig[0].MinimumEnabledTLSProtocol = stateConfig[0].MinimumEnabledTLSProtocol
		}
		if planConfig[0].NoTableScan.IsUnknown() {
			planConfig[0].NoTableScan = stateConfig[0].NoTableScan
		}
		if planConfig[0].SampleSizeBiConnector.IsUnknown() {
			planConfig[0].SampleSizeBiConnector = stateConfig[0].SampleSizeBiConnector
		}
		if planConfig[0].SampleRefreshIntervalBiConnector.IsUnknown() {
			planConfig[0].SampleRefreshIntervalBiConnector = stateConfig[0].SampleRefreshIntervalBiConnector
		}
		if planConfig[0].OplogSizeMB.IsUnknown() {
			planConfig[0].OplogSizeMB = stateConfig[0].OplogSizeMB
		}
		if planConfig[0].OplogMinRetentionHours.IsNull() {
			planConfig[0].OplogMinRetentionHours = stateConfig[0].OplogMinRetentionHours
		}
		if planConfig[0].TransactionLifetimeLimitSeconds.IsUnknown() {
			planConfig[0].TransactionLifetimeLimitSeconds = stateConfig[0].TransactionLifetimeLimitSeconds
		}

		updatedPlan, d := types.ListValueFrom(ctx, tfAdvancedConfigurationType, planConfig)
		if d.HasError() {
			return true, nil, d
		}

		return true, newAdvancedConfiguration(ctx, updatedPlan), d
	}
	return false, nil, d
}

func hasAdvancedConfigUpdated(planConfig, stateConfig []TfAdvancedConfigurationModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(planConfig) < 1 || len(stateConfig) < 1 {
		return false, diags
	}

	pc := planConfig[0]
	sc := stateConfig[0]

	if !pc.DefaultReadConcern.IsUnknown() && !pc.DefaultReadConcern.Equal(sc.DefaultReadConcern) {
		return true, diags
	}
	if !pc.DefaultWriteConcern.IsUnknown() && !pc.DefaultWriteConcern.Equal(sc.DefaultWriteConcern) {
		return true, diags
	}
	if !pc.FailIndexKeyTooLong.IsUnknown() && !pc.FailIndexKeyTooLong.Equal(sc.FailIndexKeyTooLong) {
		return true, diags
	}
	if !pc.JavascriptEnabled.IsUnknown() && !pc.JavascriptEnabled.Equal(sc.JavascriptEnabled) {
		return true, diags
	}
	if !pc.MinimumEnabledTLSProtocol.IsUnknown() && !pc.MinimumEnabledTLSProtocol.Equal(sc.MinimumEnabledTLSProtocol) {
		return true, diags
	}
	if !pc.NoTableScan.IsUnknown() && !pc.NoTableScan.Equal(sc.NoTableScan) {
		return true, diags
	}
	if !pc.SampleSizeBiConnector.IsUnknown() && !pc.SampleSizeBiConnector.Equal(sc.SampleSizeBiConnector) {
		return true, diags
	}
	if !pc.SampleRefreshIntervalBiConnector.IsUnknown() && !pc.SampleRefreshIntervalBiConnector.Equal(sc.SampleRefreshIntervalBiConnector) {
		return true, diags
	}
	if !pc.OplogSizeMB.IsUnknown() && !pc.OplogSizeMB.Equal(sc.OplogSizeMB) {
		return true, diags
	}
	if !pc.OplogMinRetentionHours.IsNull() && !pc.OplogMinRetentionHours.Equal(sc.OplogMinRetentionHours) {
		return true, diags
	}
	if !pc.TransactionLifetimeLimitSeconds.IsUnknown() && !pc.TransactionLifetimeLimitSeconds.Equal(sc.TransactionLifetimeLimitSeconds) {
		return true, diags
	}
	return false, diags
}

func replicationSpecsIfUpdated(ctx context.Context, planVal, stateVal types.List) (bool, []*matlas.AdvancedReplicationSpec, diag.Diagnostics) {
	var d diag.Diagnostics

	if !planVal.IsUnknown() && !stateVal.IsUnknown() { // check if updated/added
		updated, d := hasReplicationSpecsUpdated(ctx, planVal, stateVal)

		if d.HasError() {
			return false, nil, d
		}

		if updated {
			updatedSpecs, d := getUpdatedReplicationSpecs(ctx, planVal, stateVal)
			return true, updatedSpecs, d
			// return true, getUpdatedReplicationSpecs(), d
		}

		return false, nil, d
	}

	return false, nil, d
}

// getUpdatedReplicationSpecs This method creates API request objects to update replication specs.
// The API request objects are cerated by iterating over state replication_specs and replaces attribute
// values that are known in the plan. This is because the state replication_specs can have other Computed values
// while the plan only has values configured by the user or values from any plan modifiers or defaults set.
func getUpdatedReplicationSpecs(ctx context.Context, planVal, stateVal types.List) ([]*matlas.AdvancedReplicationSpec, diag.Diagnostics) {
	var diags diag.Diagnostics
	var res []*matlas.AdvancedReplicationSpec

	var planRepSpecs, stateRepSpecs []tfReplicationSpecRSModel
	if diags = planVal.ElementsAs(ctx, &planRepSpecs, false); diags.HasError() {
		return nil, diags
	}
	if diags = stateVal.ElementsAs(ctx, &stateRepSpecs, false); diags.HasError() {
		return nil, diags
	}

	var i int
	for i = 0; i < len(planRepSpecs); i++ {
		if i < len(stateRepSpecs) { // if rep_specs removed from config, we don't take them in API object list
			ss := stateRepSpecs[i]
			ps := planRepSpecs[i]

			// updatedRepSpec, diags := getUpdatedReplicationSpec(ctx, &ps, &ss)
			// if diags.HasError() {
			// 	return nil, diags
			// }
			updatedRepSpec := newReplicationSpec(ctx, &ps)
			updatedRepSpec.ID = ss.ID.ValueString()

			res = append(res, updatedRepSpec)
		}
	}

	if i <= len(planRepSpecs) { // new replication_specs added
		for i < len(planRepSpecs) {
			tmp := newReplicationSpec(ctx, &planRepSpecs[i])
			res = append(res, tmp)
			i++
		}
	}

	return res, diags
}

// func getUpdatedReplicationSpec(ctx context.Context, ps, ss *tfReplicationSpecRSModel) (*matlas.AdvancedReplicationSpec, diag.Diagnostics) {
// 	var diags diag.Diagnostics // TODO remove

// 	newSpec := *newReplicationSpec(ctx, ps)

// 	// if v := ps.NumShards; !v.IsUnknown() {
// 	// 	newSpec.NumShards = int(v.ValueInt64())
// 	// }
// 	// if v := ps.ZoneName; !v.IsUnknown() {
// 	// 	newSpec.ZoneName = v.ValueString()
// 	// }

// 	// newSpec.RegionConfigs = newRegionConfigs(ctx, ps.RegionsConfigs)
// 	return &newSpec, diags
// }

// func getUpdatedRegionConfigs(ctx context.Context, planVal, stateVal types.List) ([]*matlas.AdvancedRegionConfig, diag.Diagnostics) {
// 	var diags diag.Diagnostics
// 	var regionConfigs []*matlas.AdvancedRegionConfig

// 	regionConfigs = newRegionConfigs(ctx, planVal)
// 	return

// }

// hasReplicationSpecsUpdated This method checks for any attribute if known in the replication_specs plan
// should be same as it's state value. This needs to be checked as plan attributes unless configured by the user
// or any plan modifiers or defaults will always be unknown (this happens because most values in these objects
// are Optional & Computed) and hence, cannot be compared with their corresponding state value.
func hasReplicationSpecsUpdated(ctx context.Context, planVal, stateVal types.List) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	var planRepSpecs, stateRepSpecs []tfReplicationSpecRSModel
	if diags = planVal.ElementsAs(ctx, &planRepSpecs, true); diags.HasError() {
		return false, diags
	}
	if diags = stateVal.ElementsAs(ctx, &stateRepSpecs, true); diags.HasError() {
		return false, diags
	}

	if len(planRepSpecs) != len(stateRepSpecs) {
		return true, diags
	}

	for i := range planRepSpecs {
		if updated, d := hasReplicationSpecUpdated(ctx, &planRepSpecs[i], &stateRepSpecs[i]); d.HasError() || updated {
			return updated, append(diags, d...)
		}
	}

	return false, diags
}

func hasReplicationSpecUpdated(ctx context.Context, planRepSpec, stateRepSpec *tfReplicationSpecRSModel) (bool, diag.Diagnostics) {
	var hasUpdated bool
	var diags diag.Diagnostics

	if !planRepSpec.NumShards.IsUnknown() {
		hasUpdated = !planRepSpec.NumShards.Equal(stateRepSpec.NumShards)
	}

	if !planRepSpec.ZoneName.IsUnknown() { // if user has defined zone_name in config
		hasUpdated = hasUpdated || !planRepSpec.ZoneName.Equal(stateRepSpec.ZoneName)

		// if user has NOT defined zone_name in config, we set it to defaultZoneName during create so we check against that:
	} else if planRepSpec.ZoneName.IsUnknown() && stateRepSpec.ZoneName.ValueString() != DefaultZoneName {
		return true, diags
	}

	// TODO refactor
	if updated, d := hasRegionConfigsUpdated(ctx, planRepSpec.RegionsConfigs, stateRepSpec.RegionsConfigs); d.HasError() || updated {
		return updated, d
	}

	return hasUpdated, diags
}

func hasRegionConfigsUpdated(ctx context.Context, planVal, stateVal types.List) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	var planRegionConfigs, stateRegionConfigs []tfRegionsConfigModel
	planVal.ElementsAs(ctx, &planRegionConfigs, false)
	stateVal.ElementsAs(ctx, &stateRegionConfigs, false)

	if len(planRegionConfigs) != len(stateRegionConfigs) {
		return true, diags
	}

	for i := range planRegionConfigs {
		hasUpdated, d := hasRegionConfigUpdated(ctx, &planRegionConfigs[i], &stateRegionConfigs[i])

		if hasUpdated || d.HasError() {
			return hasUpdated, d
		}
	}

	return false, diags
}

func hasRegionConfigUpdated(ctx context.Context, planRegionConfig, stateRegionConfig *tfRegionsConfigModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// TODO refactor
	if !planRegionConfig.BackingProviderName.IsUnknown() && !planRegionConfig.BackingProviderName.Equal(stateRegionConfig.BackingProviderName) {
		return true, diags
	}
	if !planRegionConfig.Priority.IsUnknown() && !planRegionConfig.Priority.Equal(stateRegionConfig.Priority) {
		return true, diags
	}
	if !planRegionConfig.RegionName.IsUnknown() && !planRegionConfig.RegionName.Equal(stateRegionConfig.RegionName) {
		return true, diags
	}
	if !planRegionConfig.ProviderName.IsUnknown() && !planRegionConfig.ProviderName.Equal(stateRegionConfig.ProviderName) {
		return true, diags
	}

	if updated, d := hasRegionConfigSpecUpdated(ctx, planRegionConfig.AnalyticsSpecs, stateRegionConfig.AnalyticsSpecs); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigSpecUpdated(ctx, planRegionConfig.ElectableSpecs, stateRegionConfig.ElectableSpecs); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	// failing here,, readonly spec can be nul in state, so how to check if updated
	if updated, d := hasRegionConfigSpecUpdated(ctx, planRegionConfig.ReadOnlySpecs, stateRegionConfig.ReadOnlySpecs); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigAutoScalingSpecUpdated(ctx, planRegionConfig.AutoScaling, stateRegionConfig.AutoScaling); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigAutoScalingSpecUpdated(ctx, planRegionConfig.AnalyticsAutoScaling, stateRegionConfig.AnalyticsAutoScaling); d.HasError() || updated {
		return updated, append(diags, d...)
	}

	return false, diags
}

func hasRegionConfigAutoScalingSpecUpdated(ctx context.Context, planVal, stateVal types.List) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if planVal.IsUnknown() || stateVal.IsNull() {
		return false, diags
	}

	var planSpecs, stateSpecs []tfRegionsConfigAutoScalingSpecsModel
	if d := planVal.ElementsAs(ctx, &planSpecs, false); diags.HasError() {
		return true, append(diags, d...)
	}
	if d := stateVal.ElementsAs(ctx, &stateSpecs, false); diags.HasError() {
		return true, append(diags, d...)
	}

	if len(planSpecs) != len(stateSpecs) {
		return true, diags
	}

	ps := planSpecs[0]
	ss := stateSpecs[0]

	if !ps.ComputeEnabled.IsUnknown() && !ps.ComputeEnabled.Equal(ss.ComputeEnabled) ||
		!ps.ComputeMaxInstanceSize.IsUnknown() && !ps.ComputeMaxInstanceSize.Equal(ss.ComputeMaxInstanceSize) ||
		!ps.ComputeMinInstanceSize.IsUnknown() && !ps.ComputeMinInstanceSize.Equal(ss.ComputeMinInstanceSize) ||
		!ps.ComputeScaleDownEnabled.IsUnknown() && !ps.ComputeScaleDownEnabled.Equal(ss.ComputeScaleDownEnabled) ||
		!ps.DiskGBEnabled.IsUnknown() && !ps.DiskGBEnabled.Equal(ss.DiskGBEnabled) {
		return true, diags
	}

	return false, diags
}

func hasRegionConfigSpecUpdated(ctx context.Context, planVal, stateVal types.List) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if planVal.IsUnknown() || stateVal.IsNull() {
		return false, diags
	}

	var planSpecs, stateSpecs []tfRegionsConfigSpecsModel
	if d := planVal.ElementsAs(ctx, &planSpecs, false); diags.HasError() {
		return true, append(diags, d...)
	}
	if d := stateVal.ElementsAs(ctx, &stateSpecs, false); diags.HasError() {
		return true, append(diags, d...)
	}

	if len(planSpecs) != len(stateSpecs) {
		return true, diags
	}

	ps := planSpecs[0]
	ss := stateSpecs[0]

	// TODO refactor
	var hasUpdated bool
	if v := ps.DiskIOPS; !v.IsUnknown() {
		hasUpdated = v.Equal(ss.DiskIOPS)
	}
	if v := ps.EBSVolumeType; !v.IsUnknown() {
		hasUpdated = hasUpdated || !v.Equal(ss.EBSVolumeType)
	}
	if v := ps.InstanceSize; !v.IsUnknown() {
		hasUpdated = hasUpdated || !v.Equal(ss.InstanceSize)
	}
	if v := ps.NodeCount; !v.IsUnknown() {
		hasUpdated = hasUpdated || !v.Equal(ss.NodeCount)
	}

	return hasUpdated, diags
}

// 1. for O/C attributes plan values will be unknown unless
// they are in the config or have been updated
// 2. Optional attributes if not configured will be unknown in the plan
// 3. Required or attributes with default values will be known in the plan
// 4. Computed attributes will be unknown in the plan, and can't be updated by us anyway, so we don't need to check for them
func hasStringUpdated(p, s types.String) bool {
	if !p.IsUnknown() && !p.IsNull() && !s.IsUnknown() && !s.IsNull() {
		return p.ValueString() != s.ValueString()
	}
	return false
}

func hasOptionalStringUpdated(p, s types.String) bool {
	isPlanValPresent := !p.IsUnknown() && !p.IsNull()
	isStateValPresent := !s.IsUnknown() && !s.IsNull()

	if (isPlanValPresent && !isStateValPresent) || (!isPlanValPresent && isStateValPresent) {
		if isPlanValPresent && isStateValPresent {
			return p.ValueString() != s.ValueString()
		}
		return true
	}
	return false
}

func hasBoolUpdated(p, s types.Bool) bool {
	if !p.IsUnknown() && !p.IsNull() && !s.IsUnknown() && !s.IsNull() {
		return p.ValueBool() != s.ValueBool()
	}
	return false
}

func hasFloatUpdated(p, s types.Float64) bool {
	if !p.IsUnknown() && !p.IsNull() && !s.IsUnknown() {
		return p.ValueFloat64() != s.ValueFloat64()
	}
	return false
}
