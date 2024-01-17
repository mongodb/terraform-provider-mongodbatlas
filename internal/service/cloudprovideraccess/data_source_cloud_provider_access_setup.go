package cloudprovideraccess

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSourceSetup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasCloudProviderAccessSetupRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS"}, false),
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"aws_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"atlas_aws_account_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"atlas_assumed_role_external_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"azure_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"atlas_azure_app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasCloudProviderAccessSetupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	roleID := d.Get("role_id").(string)

	role, _, err := conn.CloudProviderAccess.GetRole(ctx, projectID, roleID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, err))
	}

	roleSchema := roleToSchemaSetup(role)
	for key, val := range roleSchema {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(fmt.Errorf(ErrorCloudProviderGetRead, err))
		}
	}

	d.SetId(id.UniqueId())
	return nil
}
