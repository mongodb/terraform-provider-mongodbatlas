package serverlessinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: fmt.Sprintf(constant.DeprecationDataSourceByDateWithExternalLink, constant.ServerlessSharedEOLDate, "https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/serverless-shared-migration-guide"),
		ReadContext:        dataSourceRead,
		Schema:             dataSourceSchema(),
	}
}

func dataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
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
			Computed: true,
		},
		"provider_settings_provider_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provider_settings_region_name": {
			Type:     schema.TypeString,
			Computed: true,
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
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"href": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"rel": {
						Type:     schema.TypeString,
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
			Computed: true,
		},
		"continuous_backup_enabled": {
			Deprecated: fmt.Sprintf(constant.DeprecationParamByDateWithExternalLink, "March 2025", "https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/serverless-shared-migration-guide"),
			Type:       schema.TypeBool,
			Optional:   true,
			Computed:   true,
		},
		"auto_indexing": {
			Deprecated: fmt.Sprintf(constant.DeprecationParamByDateWithExternalLink, "March 2025", "https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/serverless-shared-migration-guide"),
			Type:       schema.TypeBool,
			Optional:   true,
			Computed:   true,
		},
		// TODO: TEMPORARY CHANGE, DON'T MERGE
		// TODO: TEMPORARY CHANGE, DON'T MERGE
		"tags": &advancedcluster.DSTagsSchema,
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID, projectIDOk := d.GetOk("project_id")
	instanceName, instanceNameOk := d.GetOk("name")

	if !projectIDOk || !instanceNameOk {
		return diag.Errorf("project_id and name must be configured")
	}

	instance, _, err := connV2.ServerlessInstancesApi.GetServerlessInstance(ctx, projectID.(string), instanceName.(string)).Execute()
	if err != nil {
		return diag.Errorf("error getting serverless instance information: %s", err)
	}

	if err := d.Set("id", instance.GetId()); err != nil {
		return diag.Errorf("error setting `is` for serverless instance (%s): %s", d.Id(), err)
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

	if len(instance.ConnectionStrings.GetPrivateEndpoint()) > 0 {
		if err := d.Set("connection_strings_private_endpoint_srv", flattenSRVConnectionString(instance.ConnectionStrings.GetPrivateEndpoint())); err != nil {
			return diag.Errorf(errorServerlessInstanceSetting, "connection_strings_private_endpoint_srv", d.Id(), err)
		}
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

	autoIndexing, _, err := connV2.PerformanceAdvisorApi.GetServerlessAutoIndexing(ctx, projectID.(string), instanceName.(string)).Execute()
	if err != nil {
		return diag.Errorf("error getting serverless instance information for auto_indexing: %s", err)
	}
	if err := d.Set("auto_indexing", autoIndexing); err != nil {
		return diag.Errorf(errorServerlessInstanceSetting, "auto_indexing", d.Id(), err)
	}

	if err := d.Set("tags", conversion.FlattenTags(instance.GetTags())); err != nil {
		return diag.Errorf(advancedcluster.ErrorClusterAdvancedSetting, "tags", d.Id(), err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID.(string),
		"name":       instanceName.(string),
	}))

	return nil
}
