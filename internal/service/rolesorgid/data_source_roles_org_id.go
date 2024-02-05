package rolesorgid

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	apiKeyOrgList, _, err := connV2.RootApi.GetSystemStatus(ctx).Execute()
	if err != nil {
		return diag.Errorf("error getting API Key's org assigned (%s): ", err)
	}
	for _, role := range apiKeyOrgList.ApiKey.GetRoles() {
		if strings.HasPrefix(role.GetRoleName(), "ORG_") {
			if err := d.Set("org_id", role.GetOrgId()); err != nil {
				return diag.Errorf(constant.ErrorSettingAttribute, "org_id", err)
			}
			d.SetId(role.GetOrgId())
			return nil
		}
	}
	d.SetId(id.UniqueId())
	return nil
}
