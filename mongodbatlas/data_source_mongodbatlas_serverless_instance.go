package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasServerlessInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasServerlessInstanceRead,
		Schema:      returnServerlessInstanceDSSchema(),
	}
}

func dataSourceMongoDBAtlasServerlessInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	projectID, projectIDOk := d.GetOk("project_id")
	instanceName, instanceNameOk := d.GetOk("name")

	if !(projectIDOk && instanceNameOk) {
		return diag.Errorf("project_id and name must be configured")
	}

	serverlessInstance, _, err := conn.ServerlessInstances.Get(ctx, projectID.(string), instanceName.(string))
	if err != nil {
		return diag.Errorf("error getting serverless instance information: %s", err)
	}

	if err := d.Set("id", serverlessInstance.ID); err != nil {
		return diag.Errorf("error setting `is` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_settings_backing_provider_name", serverlessInstance.ProviderSettings.BackingProviderName); err != nil {
		return diag.Errorf("error setting `provider_settings_backing_provider_name` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_settings_provider_name", serverlessInstance.ProviderSettings.ProviderName); err != nil {
		return diag.Errorf("error setting `provider_settings_provider_name` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("provider_settings_region_name", serverlessInstance.ProviderSettings.RegionName); err != nil {
		return diag.Errorf("error setting `provider_settings_region_name` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("connection_strings_standard_srv", serverlessInstance.ConnectionStrings.StandardSrv); err != nil {
		return diag.Errorf("error setting `connection_strings_standard_srv` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("create_date", serverlessInstance.CreateDate); err != nil {
		return diag.Errorf("error setting `create_date` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("mongo_db_version", serverlessInstance.MongoDBVersion); err != nil {
		return diag.Errorf("error setting `mongo_db_version` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("links", flattenServerlessInstanceLinks(serverlessInstance.Links)); err != nil {
		return diag.Errorf("error setting `links` for serverless instance (%s): %s", d.Id(), err)
	}

	if err := d.Set("state_name", serverlessInstance.StateName); err != nil {
		return diag.Errorf("error setting `state_name` for serverless instance (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID.(string),
		"name":       instanceName.(string),
	}))

	return nil
}

func returnServerlessInstanceDSSchema() map[string]*schema.Schema {
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
	}
}
