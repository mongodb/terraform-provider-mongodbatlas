package advancedcluster

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/utility"
)

func (r *advancedClusterRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	conn := r.Client.Atlas
	var plan tfAdvancedClusterRSModel

	if diags := processCreateConfig(ctx, req, &plan, resp); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	projectID := plan.ProjectID.ValueString()

	request := clusterCreateRequest(ctx, &plan)
	cluster, _, err := conn.AdvancedClusters.Create(ctx, projectID, request)
	if err != nil {
		resp.Diagnostics.AddError("Unable to CREATE cluster. Error during create in Atlas", fmt.Sprintf(errorClusterAdvancedCreate, err))
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := waitClusterCreate(ctx, conn, timeout, projectID, cluster.Name); err != nil {
		resp.Diagnostics.AddError("Unable to CREATE cluster. Error during create in Atlas", err.Error())
		return
	}

	// update the advanced configuration now that the cluster has created successfully
	if v := plan.AdvancedConfiguration; !v.IsUnknown() {
		advancedConfig := newAdvancedConfiguration(ctx, v)
		if _, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, cluster.Name, advancedConfig); err != nil {
			resp.Diagnostics.AddError("Error during cluster CREATE", fmt.Sprintf(errorAdvancedClusterAdvancedConfUpdate, cluster.Name, err))
			return
		}
	}

	if err := pauseClusterIfRequired(ctx, conn, &plan, timeout); err != nil {
		resp.Diagnostics.AddError("Unable to UPDATE cluster. An error occurred when attempting to pause cluster in Atlas.", err.Error())
		return
	}

	cluster, response, err := conn.AdvancedClusters.Get(ctx, projectID, cluster.Name)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error during cluster READ from Atlas", fmt.Sprintf(errorClusterAdvancedRead, cluster.Name, err.Error()))
		return
	}

	newClusterModel, diags := newTfAdvancedClusterRSModel(ctx, conn, cluster, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterModel)...)
}

func processCreateConfig(ctx context.Context, req resource.CreateRequest,
	plan *tfAdvancedClusterRSModel,
	resp *resource.CreateResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(req.Plan.Get(ctx, plan)...)
	if diags.HasError() {
		return diags
	}

	if err := validateCreateConfig(ctx, plan); err != nil {
		diags.AddError("Invalid Create Configuration", err.Error())
		return diags
	}
	return diags
}

func clusterCreateRequest(ctx context.Context, plan *tfAdvancedClusterRSModel) *matlas.AdvancedCluster {
	request := &matlas.AdvancedCluster{
		Name:             plan.Name.ValueString(),
		ClusterType:      plan.ClusterType.ValueString(),
		ReplicationSpecs: newReplicationSpecs(ctx, plan.ReplicationSpecs),
	}

	if v := plan.BackupEnabled; !v.IsUnknown() {
		request.BackupEnabled = v.ValueBoolPointer()
	}
	if v := plan.DiskSizeGb; !v.IsUnknown() {
		request.DiskSizeGB = v.ValueFloat64Pointer()
	}
	if v := plan.EncryptionAtRestProvider; !v.IsUnknown() {
		request.EncryptionAtRestProvider = v.ValueString()
	}
	if v := plan.MongoDBMajorVersion; !v.IsUnknown() {
		request.MongoDBMajorVersion = utility.FormatMongoDBMajorVersion(v.ValueString())
	}
	if v := plan.PitEnabled; !v.IsUnknown() {
		request.PitEnabled = v.ValueBoolPointer()
	}
	if v := plan.RootCertType; !v.IsUnknown() {
		request.RootCertType = v.ValueString()
	}
	if v := plan.TerminationProtectionEnabled; !v.IsUnknown() {
		request.TerminationProtectionEnabled = v.ValueBoolPointer()
	}
	if v := plan.VersionReleaseSystem; !v.IsUnknown() {
		request.VersionReleaseSystem = v.ValueString()
	}
	if v := plan.BiConnectorConfig; !v.IsUnknown() {
		request.BiConnector = newBiConnectorConfig(ctx, plan.BiConnectorConfig)
	}
	request.Labels = append(newLabels(ctx, plan.Labels), DefaultLabel)
	request.Tags = newTags(ctx, plan.Tags)

	return request
}

func validateCreateConfig(ctx context.Context, plan *tfAdvancedClusterRSModel) error {
	if plan.AcceptDataRisksAndForceReplicaSetReconfig.ValueString() != defaultString {
		return fmt.Errorf("accept_data_risks_and_force_replica_set_reconfig can not be set in creation, only in update")
	}

	if err := validateTfConfig(ctx, plan); err != nil {
		return fmt.Errorf("invalid Create Configuration: %s", err.Error())
	}
	return nil
}

func validateTfConfig(ctx context.Context, plan *tfAdvancedClusterRSModel) error {
	var advancedConfig *matlas.ProcessArgs
	if v := plan.AdvancedConfiguration; !v.IsUnknown() {
		advancedConfig = newAdvancedConfiguration(ctx, v)
		if advancedConfig != nil && advancedConfig.OplogSizeMB != nil && *advancedConfig.OplogSizeMB <= 0 {
			return fmt.Errorf("`advanced_configuration.oplog_size_mb` cannot be <= 0")
		}
	}
	if v := plan.Labels; !v.IsNull() && ContainsLabelOrKey(newLabels(ctx, v), DefaultLabel) {
		return fmt.Errorf("you should not set `Infrastructure Tool` label, it is used for internal purposes")
	}
	return nil
}

func waitClusterCreate(ctx context.Context, conn *matlas.Client, timeout time.Duration, projectID, clusterName string) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING", "PENDING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterAdvancedRefreshFunc(ctx, clusterName, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceClusterAdvancedRefreshFunc(ctx context.Context, name, projectID string, client *matlas.Client) retry.StateRefreshFunc {
	return func() (any, string, error) {
		c, resp, err := client.AdvancedClusters.Get(ctx, projectID, name)

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && c == nil && resp == nil {
			return nil, "", err
		}

		if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			if resp.StatusCode == 503 {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		if c.StateName != "" {
			log.Printf("[DEBUG] status for MongoDB cluster: %s: %s", name, c.StateName)
		}

		return c, c.StateName, nil
	}
}
