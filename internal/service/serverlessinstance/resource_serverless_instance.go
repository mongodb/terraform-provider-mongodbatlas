package serverlessinstance

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/atlas-sdk/v20231115005/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorServerlessInstanceSetting = "error setting `%s` for serverless instance (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: resourceSchema(),
	}
}

func resourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"provider_settings_backing_provider_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"provider_settings_provider_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"provider_settings_region_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"connection_strings_standard_srv": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"connection_strings_private_endpoint_srv": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"create_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mongo_db_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": {
			Type:       schema.TypeSet,
			Optional:   true,
			Computed:   true,
			ConfigMode: schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"href": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"rel": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				}},
		},
		"state_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"termination_protection_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"continuous_backup_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"tags": &advancedcluster.RSTagsSchema,
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	name := d.Get("name").(string)

	serverlessProviderSettings := &matlas.ServerlessProviderSettings{
		BackingProviderName: d.Get("provider_settings_backing_provider_name").(string),
		ProviderName:        d.Get("provider_settings_provider_name").(string),
		RegionName:          d.Get("provider_settings_region_name").(string),
	}

	serverlessBackupOptions := &matlas.ServerlessBackupOptions{
		ServerlessContinuousBackupEnabled: pointy.Bool(d.Get("continuous_backup_enabled").(bool)),
	}

	serverlessInstanceRequest := &matlas.ServerlessCreateRequestParams{
		Name:                         name,
		ProviderSettings:             serverlessProviderSettings,
		ServerlessBackupOptions:      serverlessBackupOptions,
		TerminationProtectionEnabled: pointy.Bool(d.Get("termination_protection_enabled").(bool)),
	}

	if _, ok := d.GetOk("tags"); ok {
		tagsSlice := advancedcluster.ExpandTagSliceFromSetSchema(d)
		serverlessInstanceRequest.Tag = &tagsSlice
	}

	_, _, err := conn.ServerlessInstances.Create(ctx, projectID, serverlessInstanceRequest)
	if err != nil {
		return diag.Errorf("error creating serverless instance: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING", "PENDING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceRefreshFunc(ctx, d.Get("name").(string), projectID, conn),
		Timeout:    3 * time.Hour,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating MongoDB Serverless Instance: %s", err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"name":       name,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	instanceName := ids["name"]

	instance, _, err := connV2.ServerlessInstancesApi.GetServerlessInstance(ctx, projectID, instanceName).Execute()
	if err != nil {
		// case 404: deleted in the backend case
		if strings.Contains(err.Error(), "404") && !d.IsNewResource() {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting serverless instance information: %s", err)
	}

	if err := d.Set("id", instance.GetId()); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "id", d.Id(), err)
	}

	if err := d.Set("provider_settings_backing_provider_name", instance.ProviderSettings.GetBackingProviderName()); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "provider_settings_backing_provider_name", d.Id(), err)
	}

	if err := d.Set("provider_settings_provider_name", instance.ProviderSettings.GetProviderName()); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "provider_settings_provider_name", d.Id(), err)
	}

	if err := d.Set("provider_settings_region_name", instance.ProviderSettings.GetRegionName()); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "provider_settings_region_name", d.Id(), err)
	}

	if err := d.Set("connection_strings_standard_srv", instance.ConnectionStrings.GetStandardSrv()); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "connection_strings_standard_srv", d.Id(), err)
	}

	if err := d.Set("connection_strings_private_endpoint_srv", flattenSRVConnectionString(instance.ConnectionStrings.GetPrivateEndpoint())); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "connection_strings_private_endpoint_srv", d.Id(), err)
	}

	if err := d.Set("create_date", conversion.TimePtrToStringPtr(instance.CreateDate)); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "create_date", d.Id(), err)
	}

	if err := d.Set("mongo_db_version", instance.GetMongoDBVersion()); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "mongo_db_version", d.Id(), err)
	}

	if err := d.Set("links", conversion.FlattenLinks(instance.GetLinks())); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "links", d.Id(), err)
	}

	if err := d.Set("state_name", instance.GetStateName()); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "state_name", d.Id(), err)
	}

	if err := d.Set("termination_protection_enabled", instance.GetTerminationProtectionEnabled()); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "termination_protection_enabled", d.Id(), err)
	}

	if err := d.Set("continuous_backup_enabled", instance.ServerlessBackupOptions.GetServerlessContinuousBackupEnabled()); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "continuous_backup_enabled", d.Id(), err)
	}

	if err := d.Set("tags", conversion.FlattenTags(instance.GetTags())); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "tags", d.Id(), err)
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	instanceName := ids["name"]

	if d.HasChange("termination_protection_enabled") || d.HasChange("continuous_backup_enabled") || d.HasChange("tags") {
		serverlessBackupOptions := &matlas.ServerlessBackupOptions{
			ServerlessContinuousBackupEnabled: pointy.Bool(d.Get("continuous_backup_enabled").(bool)),
		}

		ServerlessUpdateRequestParams := &matlas.ServerlessUpdateRequestParams{
			ServerlessBackupOptions:      serverlessBackupOptions,
			TerminationProtectionEnabled: pointy.Bool(d.Get("termination_protection_enabled").(bool)),
		}

		if d.HasChange("tags") {
			tags := advancedcluster.ExpandTagSliceFromSetSchema(d)
			ServerlessUpdateRequestParams.Tag = &tags
		}

		_, _, err := conn.ServerlessInstances.Update(ctx, projectID, instanceName, ServerlessUpdateRequestParams)
		if err != nil {
			return diag.Errorf("error updating serverless instance: %s", err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING", "PENDING"},
			Target:     []string{"IDLE"},
			Refresh:    resourceRefreshFunc(ctx, d.Get("name").(string), projectID, conn),
			Timeout:    3 * time.Hour,
			MinTimeout: 1 * time.Minute,
			Delay:      3 * time.Minute,
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error updating MongoDB Serverless Instance: %s", err)
		}
	}
	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	serverlessName := ids["name"]

	_, err := conn.ServerlessInstances.Delete(ctx, projectID, serverlessName)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting MongoDB Serverless Instance (%s): %s", serverlessName, err))
	}

	log.Println("[INFO] Waiting for MongoDB Serverless Instance to be destroyed")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceRefreshFunc(ctx, serverlessName, projectID, conn),
		Timeout:    3 * time.Hour,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting MongoDB Serverless Instance (%s): %s", serverlessName, err))
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID, name, err := splitImportID(d.Id())
	if err != nil {
		return nil, err
	}

	u, _, err := connV2.ServerlessInstancesApi.GetServerlessInstance(ctx, *projectID, *name).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import cluster %s in project %s, error: %s", *name, *projectID, err)
	}

	if err := d.Set("project_id", u.GetGroupId()); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "project_id", u.GetId(), err)
	}

	if err := d.Set("name", u.Name); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "name", u.GetId(), err)
	}

	if err := d.Set("continuous_backup_enabled", u.ServerlessBackupOptions.ServerlessContinuousBackupEnabled); err != nil {
		log.Printf(advancedcluster.ErrorClusterSetting, "continuous_backup_enabled", u.GetId(), err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": *projectID,
		"name":       u.GetName(),
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceRefreshFunc(ctx context.Context, name, projectID string, client *matlas.Client) retry.StateRefreshFunc {
	return func() (any, string, error) {
		c, resp, err := client.ServerlessInstances.Get(ctx, projectID, name)

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && c == nil && resp == nil {
			return nil, "", err
		} else if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			if resp.StatusCode == 503 {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		if c.StateName != "" {
			log.Printf("[DEBUG] status for MongoDB Serverless Instance: %s: %s", name, c.StateName)
		}

		return c, c.StateName, nil
	}
}

func flattenSRVConnectionString(list []admin.ServerlessConnectionStringsPrivateEndpointList) []any {
	ret := make([]any, len(list))
	for i, elm := range list {
		ret[i] = elm.GetSrvConnectionString()
	}
	return ret
}

func splitImportID(id string) (projectID, instanceName *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a serverless instance, use the format {project_id}-{name}")
		return
	}

	projectID = &parts[1]
	instanceName = &parts[2]

	return
}
