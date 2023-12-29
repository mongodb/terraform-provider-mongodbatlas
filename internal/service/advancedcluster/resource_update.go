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
)

func (r *advancedClusterRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	conn := r.Client.Atlas
	var state, plan tfAdvancedClusterRSModel

	if diags := processUpdateConfig(ctx, req, resp, &plan, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	projectID, clusterName := decodeClusterID(state.ID.ValueString())

	if diags := updateOrUpgradeCluster(ctx, conn, &state, &plan, projectID, clusterName); diags.HasError() {
		resp.Diagnostics.Append(diags...)
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

func processUpdateConfig(ctx context.Context, req resource.UpdateRequest,
	resp *resource.UpdateResponse,
	plan *tfAdvancedClusterRSModel,
	state *tfAdvancedClusterRSModel) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(req.Plan.Get(ctx, plan)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(req.State.Get(ctx, state)...)
	if diags.HasError() {
		return diags
	}

	if err := validateTfConfig(ctx, plan); err != nil {
		diags.AddError("Invalid Create Configuration", err.Error())
		return diags
	}
	return diags
}

func updateOrUpgradeCluster(ctx context.Context, conn *matlas.Client, state, plan *tfAdvancedClusterRSModel,
	projectID, clusterName string) diag.Diagnostics {
	var diags diag.Diagnostics

	timeout, _ := plan.Timeouts.Update(ctx, defaultTimeout)

	if upgradeRequest := getUpgradeRequest(ctx, state, plan); upgradeRequest != nil {
		if _, _, err := handleClusterUpgrade(ctx, conn, upgradeRequest, projectID, clusterName, timeout); err != nil {
			diags.AddError("Unable to UPDATE cluster. An error occurred while upgrading cluster.", err.Error())
		}
	} else {
		diags.Append(handleClusterUpdate(ctx, conn, state, plan, timeout)...)
	}
	return diags
}

func decodeClusterID(id string) (projectID, clusterName string) {
	ids := conversion.DecodeStateID(id)
	return ids["project_id"], ids["cluster_name"]
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

func handleClusterUpgrade(ctx context.Context, conn *matlas.Client, request *matlas.Cluster, projectID, name string, timeout time.Duration) (*matlas.Cluster, *matlas.Response, error) {
	request.Name = name

	cluster, resp, err := conn.Clusters.Upgrade(ctx, projectID, request)
	if err != nil {
		return nil, nil, err
	}

	if err := waitClusterUpdate(ctx, conn, timeout, projectID, cluster.Name); err != nil {
		return nil, nil, err
	}

	return cluster, resp, nil
}

func handleClusterUpdate(ctx context.Context, conn *matlas.Client, state, plan *tfAdvancedClusterRSModel, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	projectID, clusterName := decodeClusterID(state.ID.ValueString())

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
	if hasStringUpdated(plan.MongoDBMajorVersion, state.MongoDBMajorVersion) {
		cluster.MongoDBMajorVersion = plan.MongoDBMajorVersion.ValueString()
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

	// Labels & Tags are optional so state/plan will either be null or known
	if !reflect.DeepEqual(plan.Labels, state.Labels) {
		if ContainsLabelOrKey(newLabels(ctx, plan.Labels), DefaultLabel) {
			diags.AddError("Unable to UPDATE cluster. An error occurred when updating labels.", "you should not set `Infrastructure Tool` label, it is used for internal purposes")
			return diags
		}
		cluster.Labels = newLabels(ctx, plan.Labels)
	}
	if !reflect.DeepEqual(plan.Tags, state.Tags) {
		cluster.Tags = newTags(ctx, plan.Tags)
	}

	if updated, newPlan, d := biConnectorConfigIfUpdated(ctx, plan.BiConnectorConfig, state.BiConnectorConfig); !d.HasError() && updated {
		cluster.BiConnector = newPlan
	}
	if updated, newPlan, d := replicationSpecsIfUpdated(ctx, plan.ReplicationSpecs, state.ReplicationSpecs); !d.HasError() && updated {
		cluster.ReplicationSpecs = newPlan
	}
	if updated, newPlan, d := advancedConfigIfUpdated(ctx, plan.AdvancedConfiguration, state.AdvancedConfiguration); !d.HasError() && updated {
		if !reflect.DeepEqual(newPlan, matlas.ProcessArgs{}) {
			_, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, clusterName, newPlan)
			if err != nil {
				diags.AddError("Unable to UPDATE cluster. An error occurred when updating advanced_configuration.", err.Error())
				return diags
			}
		}
	}

	// Has changes
	if !reflect.DeepEqual(cluster, clusterChangeDetect) {
		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, resp, err := updateCluster(ctx, conn, cluster, projectID, clusterName, timeout)
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

	if err := pauseClusterIfRequired(ctx, conn, plan, timeout); err != nil {
		diags.AddError("Unable to UPDATE cluster. An error occurred when attempting to pause cluster in Atlas.", err.Error())
		return diags
	}

	return diags
}

func pauseClusterIfRequired(ctx context.Context, conn *matlas.Client, plan *tfAdvancedClusterRSModel, timeout time.Duration) error {
	if plan.Paused.ValueBool() {
		clusterRequest := &matlas.AdvancedCluster{
			Paused: pointy.Bool(true),
		}

		if _, _, err := updateCluster(ctx, conn, clusterRequest, plan.ProjectID.ValueString(), plan.Name.ValueString(), timeout); err != nil {
			return err
		}
	}
	return nil
}

func updateCluster(ctx context.Context, conn *matlas.Client, request *matlas.AdvancedCluster, projectID, name string, timeout time.Duration,
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

	if !planVal.IsUnknown() && !stateVal.IsUnknown() {
		updated, d := hasReplicationSpecsUpdated(ctx, planVal, stateVal)

		if d.HasError() {
			return false, nil, d
		}

		if updated {
			updatedSpecs, d := getUpdatedReplicationSpecs(ctx, planVal, stateVal)
			return true, updatedSpecs, d
		}

		return false, nil, d
	}

	return false, nil, d
}

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
		if i > len(stateRepSpecs) {
			continue
		}

		ss := stateRepSpecs[i]
		ps := planRepSpecs[i]

		updatedRepSpec := newReplicationSpec(ctx, &ps)
		updatedRepSpec.ID = ss.ID.ValueString()

		res = append(res, updatedRepSpec)
	}

	if i <= len(planRepSpecs) {
		for i < len(planRepSpecs) {
			tmp := newReplicationSpec(ctx, &planRepSpecs[i])
			res = append(res, tmp)
			i++
		}
	}

	return res, diags
}

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
	var diags diag.Diagnostics

	if !planRepSpec.NumShards.IsUnknown() && !planRepSpec.NumShards.Equal(stateRepSpec.NumShards) {
		return true, diags
	}

	// compare ZoneName if it is known, or check against default if it is unknown
	if !planRepSpec.ZoneName.IsUnknown() {
		if !planRepSpec.ZoneName.Equal(stateRepSpec.ZoneName) {
			return true, diags
		}
	} else if stateRepSpec.ZoneName.ValueString() != DefaultZoneName {
		return true, diags
	}

	updated, d := hasRegionConfigsUpdated(ctx, planRepSpec.RegionsConfigs, stateRepSpec.RegionsConfigs)
	if d.HasError() || updated {
		return updated, append(diags, d...)
	}

	return false, diags
}

func hasRegionConfigsUpdated(ctx context.Context, planVal, stateVal types.List) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	var planRegionConfigs, stateRegionConfigs []tfRegionsConfigModel

	if d := planVal.ElementsAs(ctx, &planRegionConfigs, false); d.HasError() {
		return false, append(diags, d...)
	}
	if d := stateVal.ElementsAs(ctx, &stateRegionConfigs, false); d.HasError() {
		return false, append(diags, d...)
	}

	if len(planRegionConfigs) != len(stateRegionConfigs) {
		return true, diags
	}

	for i := range planRegionConfigs {
		planConfig := &planRegionConfigs[i]
		stateConfig := &stateRegionConfigs[i]
		if updated, d := hasRegionConfigUpdated(ctx, planConfig, stateConfig); d.HasError() || updated {
			return updated, append(diags, d...)
		}
	}

	return false, diags
}

func hasRegionConfigUpdated(ctx context.Context, planConfig, stateConfig *tfRegionsConfigModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !planConfig.BackingProviderName.IsUnknown() && !planConfig.BackingProviderName.Equal(stateConfig.BackingProviderName) {
		return true, diags
	}
	if !planConfig.Priority.IsUnknown() && !planConfig.Priority.Equal(stateConfig.Priority) {
		return true, diags
	}
	if !planConfig.RegionName.IsUnknown() && !planConfig.RegionName.Equal(stateConfig.RegionName) {
		return true, diags
	}
	if !planConfig.ProviderName.IsUnknown() && !planConfig.ProviderName.Equal(stateConfig.ProviderName) {
		return true, diags
	}

	if updated, d := hasRegionConfigSpecUpdated(ctx, planConfig.AnalyticsSpecs, stateConfig.AnalyticsSpecs); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigSpecUpdated(ctx, planConfig.ElectableSpecs, stateConfig.ElectableSpecs); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigSpecUpdated(ctx, planConfig.ReadOnlySpecs, stateConfig.ReadOnlySpecs); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigAutoScalingSpecUpdated(ctx, planConfig.AutoScaling, stateConfig.AutoScaling); d.HasError() || updated {
		return updated, append(diags, d...)
	}
	if updated, d := hasRegionConfigAutoScalingSpecUpdated(ctx, planConfig.AnalyticsAutoScaling, stateConfig.AnalyticsAutoScaling); d.HasError() || updated {
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
	if d := planVal.ElementsAs(ctx, &planSpecs, false); d.HasError() {
		return false, append(diags, d...)
	}
	if d := stateVal.ElementsAs(ctx, &stateSpecs, false); d.HasError() {
		return false, append(diags, d...)
	}

	if len(planSpecs) != len(stateSpecs) {
		return true, diags
	}

	ps := planSpecs[0]
	ss := stateSpecs[0]

	hasUpdated := (!ps.DiskIOPS.IsUnknown() && !ps.DiskIOPS.Equal(ss.DiskIOPS)) ||
		(!ps.EBSVolumeType.IsUnknown() && !ps.EBSVolumeType.Equal(ss.EBSVolumeType)) ||
		(!ps.InstanceSize.IsUnknown() && !ps.InstanceSize.Equal(ss.InstanceSize)) ||
		(!ps.NodeCount.IsUnknown() && !ps.NodeCount.Equal(ss.NodeCount))

	return hasUpdated, diags
}

func hasOptionalStringUpdated(plan, state types.String) bool {
	if (plan.IsUnknown() || plan.IsNull()) && (state.IsUnknown() || state.IsNull()) {
		return false
	}

	// If one of the values is unknown or null but not both, or if both are present but different, it's an update.
	return plan.ValueString() != state.ValueString()
}

func hasStringUpdated(plan, state types.String) bool {
	return !plan.IsUnknown() && !plan.IsNull() &&
		!state.IsUnknown() && !state.IsNull() &&
		plan.ValueString() != state.ValueString()
}

func hasBoolUpdated(plan, state types.Bool) bool {
	return !plan.IsUnknown() && !plan.IsNull() &&
		!state.IsUnknown() && !state.IsNull() &&
		plan.ValueBool() != state.ValueBool()
}

func hasFloatUpdated(plan, state types.Float64) bool {
	return !plan.IsUnknown() && !plan.IsNull() &&
		!state.IsUnknown() &&
		plan.ValueFloat64() != state.ValueFloat64()
}

func doesAdvancedReplicationSpecMatchAPI(tfObject *tfReplicationSpecRSModel, apiObject *matlas.AdvancedReplicationSpec) bool {
	return tfObject.ID.ValueString() == apiObject.ID || (tfObject.ID.IsNull() && tfObject.ZoneName.ValueString() == apiObject.ZoneName)
}

func removeDefaultLabel(labels []TfLabelModel) []TfLabelModel {
	result := make([]TfLabelModel, 0)

	for _, item := range labels {
		if item.Key.ValueString() == DefaultLabel.Key && item.Value.ValueString() == DefaultLabel.Value {
			continue
		}
		result = append(result, item)
	}

	return result
}
