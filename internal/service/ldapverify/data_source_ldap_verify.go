package ldapverify

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bind_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
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
					},
				},
			},
			"validations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	requestID := d.Get("request_id").(string)
	ldapResp, _, err := connV2.LDAPConfigurationApi.GetUserSecurityVerify(ctx, projectID, requestID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorRead, projectID, err))
	}
	if err := d.Set("hostname", ldapResp.Request.GetHostname()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "hostname", d.Id(), err))
	}
	if err := d.Set("port", ldapResp.Request.GetPort()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "port", d.Id(), err))
	}
	if err := d.Set("bind_username", ldapResp.Request.GetBindUsername()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "bind_username", d.Id(), err))
	}
	if err := d.Set("links", conversion.FlattenLinks(ldapResp.GetLinks())); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "links", d.Id(), err))
	}
	if err := d.Set("validations", flattenValidations(ldapResp.GetValidations())); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "validations", d.Id(), err))
	}
	if err := d.Set("request_id", ldapResp.GetRequestId()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "request_id", d.Id(), err))
	}
	if err := d.Set("status", ldapResp.GetStatus()); err != nil {
		return diag.FromErr(fmt.Errorf(errorSettings, "status", d.Id(), err))
	}
	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"request_id": ldapResp.GetRequestId(),
	}))
	return nil
}
